package queue

// File que for distributing multiple input files over GPUs.

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/MathieuMoalic/amumax/src/api"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/url"
)

var (
	exitStatus       atom = 0
	numOK, numFailed atom = 0, 0
)

func RunQueue(files []string, flags *flags.FlagsType) {
	s := NewStateTab(files)
	host, port, path, err := url.ParseAddrPath(flags.WebUIQueueAddress)
	log_old.Log.PanicIfError(err)
	if path != "" {
		log_old.Log.ErrAndExit("Path not supported for queue web UI")
	}
	addr, _, err := api.FindAvailablePort(host, port)
	log_old.Log.PanicIfError(err)
	log_old.Log.Info("Queue web UI at %v", addr)
	s.printJobList()
	go s.ListenAndServe(addr)
	s.Run(flags)
	log_old.Log.Command(fmt.Sprintf("%d OK; %d Failed", numOK.get(), numFailed.get()))
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
func (s *stateTab) Run(flags *flags.FlagsType) {
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
			run(j.inFile, gpu, flags)
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

func run(inFile string, gpu int, flags *flags.FlagsType) {
	// invalid flags: Version, Update, Gpu, Interactive, OutputDir, SelfTest
	// add all of the other flags to the command line
	cmd := []string{os.Args[0]}

	// Add valid flags to the command line
	if flags.Debug {
		cmd = append(cmd, "--debug")
	}
	if flags.Vet {
		cmd = append(cmd, "--vet")
	}
	if flags.CacheDir != fmt.Sprintf("%v/amumax_kernels", os.TempDir()) {
		cmd = append(cmd, "--cache", flags.CacheDir)
	}
	if flags.Silent {
		cmd = append(cmd, "--silent")
	}
	if flags.Sync {
		cmd = append(cmd, "--sync")
	}
	if flags.ForceClean {
		cmd = append(cmd, "--force-clean")
	}
	if flags.SkipExists {
		cmd = append(cmd, "--skip-exist")
	}
	if flags.HideProgressBar {
		cmd = append(cmd, "--hide-progress-bar")
	}
	if flags.Tunnel != "" {
		cmd = append(cmd, "--tunnel", flags.Tunnel)
	}
	if flags.Insecure {
		cmd = append(cmd, "--insecure")
	}
	if flags.NewEngine {
		cmd = append(cmd, "--new-parser")
	}
	if flags.WebUIDisabled {
		cmd = append(cmd, "--webui-disable")
	}
	if flags.WebUIAddress != "localhost:35367" {
		cmd = append(cmd, "--webui-addr", flags.WebUIAddress)
	}
	// GPU and Input File
	cmd = append(cmd, "--gpu", fmt.Sprintf("%d", gpu), inFile)

	// cmd := []string{os.Args[0], "-g", fmt.Sprint(gpu), inFile}
	// log.Log.Command(fmt.Sprintf("Running %s on GPU %d", inFile, gpu))
	// concat all the flags and the input file
	cmdString := ""
	for _, c := range cmd {
		cmdString += c + " "
	}
	log_old.Log.Command(fmt.Sprintf("Running %s", cmdString))
	err := exec.Command(cmd[0], cmd[1:]...).Run()
	if err != nil {
		log_old.Log.Command(fmt.Sprintf("FAILED %s on GPU %d: %v", inFile, gpu, err))
		exitStatus = 1
		numFailed.inc()
		return
	}
	log_old.Log.Command(fmt.Sprintf("DONE %s on GPU %d", inFile, gpu))
	numOK.inc()
}

func initGPUs(nGpu int) chan int {
	if nGpu == 0 {
		log_old.Log.ErrAndExit("no GPUs available")
	}
	idle := make(chan int, nGpu)
	for i := 0; i < nGpu; i++ {
		idle <- i
	}
	return idle
}

func (s *stateTab) printJobList() {
	s.lock.Lock()
	defer s.lock.Unlock()
	log_old.Log.Command("Job list:")
	for i, j := range s.jobs {
		log_old.Log.Command(fmt.Sprintf("%3d %v %v", i, j.inFile, j.webAddr))
	}
	log_old.Log.Command("Starting ...")
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
		log_old.Log.PanicIfError(http.ListenAndServe(addr, nil))
	}()
}

func (s *stateTab) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.RenderHTML(w)
}
