package main

// File que for distributing multiple input files over GPUs.

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/MathieuMoalic/amumax/cuda/cu"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
)

var (
	exitStatus       atom = 0
	numOK, numFailed atom = 0, 0
)

func RunQueue(files []string) {
	s := NewStateTab(files)
	s.PrintTo(os.Stdout)
	addr := fmt.Sprint(*engine.Flag_webui_queue_host, ":", *engine.Flag_webui_queue_port)
	go s.ListenAndServe(addr)
	s.Run()
	fmt.Println(numOK.get(), "OK, ", numFailed.get(), "failed")
	os.Exit(int(exitStatus))
}

// StateTab holds the queue state (list of jobs + statuses).
// All operations are atomic.
type stateTab struct {
	lock sync.Mutex
	jobs []job
	next int
}

// Job info.
type job struct {
	inFile  string // input file to run
	webAddr string // http address for gui of running process
	uid     int
}

// NewStateTab constructs a queue for the given input files.
// After construction, it is accessed atomically.
func NewStateTab(inFiles []string) *stateTab {
	s := new(stateTab)
	s.jobs = make([]job, len(inFiles))
	for i, f := range inFiles {
		s.jobs[i] = job{inFile: f, uid: i}
	}
	return s
}

// StartNext advances the next job and marks it running, setting its webAddr to indicate the GUI url.
// A copy of the job info is returned, the original remains unmodified.
// ok is false if there is no next job.
func (s *stateTab) StartNext(webAddr string) (next job, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.next >= len(s.jobs) {
		return job{}, false
	}
	s.jobs[s.next].webAddr = webAddr
	jobCopy := s.jobs[s.next]
	s.next++
	return jobCopy, true
}

// Finish marks the job with j's uid as finished.
func (s *stateTab) Finish(j job) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.jobs[j.uid].webAddr = ""
}

// Runs all the jobs in stateTab.
func (s *stateTab) Run() {
	nGPU := cu.DeviceGetCount()
	idle := initGPUs(nGPU)
	for {
		gpu := <-idle
		addr := fmt.Sprint(":", 35367+gpu)
		j, ok := s.StartNext(addr)
		if !ok {
			break
		}
		go func() {
			run(j.inFile, gpu)
			s.Finish(j)
			idle <- gpu
		}()
	}
	// drain remaining tasks (one already done)
	for i := 1; i < nGPU; i++ {
		<-idle
	}
}

type atom int32

func (a *atom) get() int { return int(atomic.LoadInt32((*int32)(a))) }
func (a *atom) inc()     { atomic.AddInt32((*int32)(a), 1) }

func run(inFile string, gpu int) {
	// overridden flags
	gpuFlag := fmt.Sprint(`-gpu=`, gpu)
	// httpFlag := fmt.Sprint(`-http=`, webAddr)

	// pass through flags
	flags := []string{gpuFlag}
	flag.Visit(func(f *flag.Flag) {
		if f.Name != "gpu" && f.Name != "http" && f.Name != "failfast" {
			flags = append(flags, fmt.Sprintf("-%v=%v", f.Name, f.Value))
		}
	})
	flags = append(flags, inFile)

	_ = exec.Command(os.Args[0], flags...)
	numOK.inc()
}

func initGPUs(nGpu int) chan int {
	if nGpu == 0 {
		util.Log.ErrAndExit("no GPUs available")
	}
	idle := make(chan int, nGpu)
	for i := 0; i < nGpu; i++ {
		idle <- i
	}
	return idle
}

func (s *stateTab) PrintTo(w io.Writer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, j := range s.jobs {
		fmt.Fprintf(w, "%3d %v %v\n", i, j.inFile, j.webAddr)
	}
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (s *stateTab) RenderHTML(w io.Writer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Fprintln(w, ` 
<!DOCTYPE html> <html> <head> 
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="refresh" content="1">
	<style media="all" type="text/css">

		* {
			box-sizing: border-box;
		}

		html body {
			background: rgb(40, 42, 54);
			color: rgb(248, 248, 242);
			margin-left: 5%;
			margin-right: 5%;
			font-family: sans-serif;
			font-size: 14px;
		}
		table { border-collapse: collapse; }
		td        { padding: 1px 5px; }
		hr        { border-style: none; border-top: 1px solid #CCCCCC; }
		a         { color: #50fa7b; text-decoration: none; }
		div       { margin-left: 20px; margin-top: 5px; margin-bottom: 20px; }
		div#footer{ color:gray; font-size:14px; border:none; }
		.ErrorBox { color: red; font-weight: bold; font-size: 1em; } 
		.TextBox  { border:solid; border-color:#BBBBBB; border-width:1px; padding-left:4px; }
		textarea  { border:solid; border-color:#BBBBBB; border-width:1px; padding-left:4px; color:gray; font-size: 1em; }
	</style>
	</head><body>
	<span style="color:#ffb86c; font-weight:bold; font-size:1.5em"> amumax queue status </span><br/>
	<hr/>
	<pre>`)

	for _, j := range s.jobs {
		if j.webAddr != "" {
			fmt.Fprint(w, `<b>`, j.uid, ` <a href="`, "http://", GetLocalIP()+j.webAddr, `">`, j.inFile, " ", j.webAddr, "</a></b>\n")
		} else {
			fmt.Fprint(w, j.uid, " ", j.inFile, "\n")
		}
	}
	fmt.Fprintln(w, `</pre><hr/></body></html>`)
}

func (s *stateTab) ListenAndServe(addr string) {
	http.Handle("/", s)
	go func() {
		util.Log.PanicIfError(http.ListenAndServe(addr, nil))
	}()
}

func (s *stateTab) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.RenderHTML(w)
}
