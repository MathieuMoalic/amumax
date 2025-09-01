package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type MetricsState struct {
	ws *WebSocketManager

	PID             int     `msgpack:"pid"`
	Error           string  `msgpack:"error"`
	CpuPercent      float64 `msgpack:"cpuPercent"`
	CpuPercentTotal float64 `msgpack:"cpuPercentTotal"`
	RamPercent      float32 `msgpack:"ramPercent"`
	RamPercentTotal float64 `msgpack:"ramPercentTotal"`

	GpuName               string  `msgpack:"gpuName"`
	GpuUUID               string  `msgpack:"gpuUUID"`
	GpuUtilizationPercent int     `msgpack:"gpuUtilizationPercent"`
	GpuTemperature        int     `msgpack:"gpuTemperature"`
	GpuPowerDraw          float64 `msgpack:"gpuPowerDraw"`
	GpuPowerLimit         float64 `msgpack:"gpuPowerLimit"`
	GpuVramUsed           int     `msgpack:"gpuVramUsed"`
	GpuVramTotal          int     `msgpack:"gpuVramTotal"`
}

func initMetricsAPI(e *echo.Group, ws *WebSocketManager) *MetricsState {
	pid := os.Getpid()
	metricState := &MetricsState{PID: pid, ws: ws}
	metricState.Update()
	e.POST("/api/metrics/reset-error", metricState.postMetricsReset)
	return metricState
}

func (m *MetricsState) Update() {
	if m.Error != "" {
		return
	}
	m.getProgramCPUUsage()
	m.getTotalCPUUsage()
	m.getTotalRAMUsage()
	m.getGPUStats1()
	m.getGPUStats2()
}

// Function to get GPU stats via `nvidia-smi`
func (m *MetricsState) getGPUStats1() {
	// Get all compute apps, filter by pid, and get the gpu uuid, name, and vram used
	cmd := exec.Command("nvidia-smi",
		"--query-compute-apps=pid,used_gpu_memory,gpu_name,gpu_uuid",
		"--format=csv,nounits,noheader",
	)
	expectedFields := 4
	output, err := cmd.Output()
	if err != nil {
		m.Error = "Error getting gpu stats"
		log.Log.Warn("%s", m.Error)
		return
	}
	var pid_fields []string
	// find the line with the pid
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ", ")
		// the last line is empty
		if line == "" {
			continue
		}
		if len(fields) == expectedFields {
			pid, err1 := strconv.Atoi(fields[0])
			if err1 != nil {
				m.Error = "Error parsing pid"
				log.Log.Warn("%s", m.Error)
				return
			}
			if pid == m.PID {
				pid_fields = fields
			}
		} else {
			m.Error = fmt.Sprintf("Expected %d fields in nvidia-smi output, got %v: %s", expectedFields, len(fields), output)
			log.Log.Warn("%s", m.Error)
			return
		}
	}
	if pid_fields == nil {
		m.Error = "Couldn't find the process ID in nvidia-smi output, this is normal if running inside a container"
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuVramUsed, err = strconv.Atoi(pid_fields[1])
	if err != nil {
		m.Error = "Error parsing vram used"
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuName = pid_fields[2]
	m.GpuUUID = pid_fields[3]
}

func (m *MetricsState) getGPUStats2() {
	if m.GpuUUID == "" {
		return
	}
	// filter using the gpu uuid
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=uuid,temperature.gpu,power.draw,memory.total,utilization.gpu,power.limit",
		"--format=csv,nounits,noheader",
	)
	expectedFields := 6
	output, err := cmd.Output()
	if err != nil {
		m.Error = "Error getting gpu stats"
		log.Log.Warn("%s", m.Error)
		return
	}
	var uuid_fields []string
	// find the line with the uuid
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// the last line is empty
		if line == "" {
			continue
		}
		fields := strings.Split(line, ", ")
		if len(fields) == expectedFields {
			if fields[0] == m.GpuUUID {
				uuid_fields = fields
			}
		} else {
			m.Error = fmt.Sprintf("Expected %d fields in nvidia-smi output, got %v: %s", expectedFields, len(fields), output)
			log.Log.Warn("%s", m.Error)
			return
		}
	}
	if uuid_fields == nil {
		m.Error = "No correct UUID found in nvidia-smi output"
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuTemperature, err = strconv.Atoi(uuid_fields[1]) // temperature in C
	if err != nil {
		m.Error = fmt.Sprintf("Error parsing GPU temperature: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuPowerDraw, err = strconv.ParseFloat(uuid_fields[2], 32) // power draw in W
	if err != nil {
		m.Error = fmt.Sprintf("Error parsing power draw: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuVramTotal, err = strconv.Atoi(uuid_fields[3]) // vram total in MiB
	if err != nil {
		m.Error = fmt.Sprintf("Error parsing vram total: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuUtilizationPercent, err = strconv.Atoi(uuid_fields[4]) // gpu utilization in %
	if err != nil {
		m.Error = fmt.Sprintf("Error parsing gpu utilization: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.GpuPowerLimit, err = strconv.ParseFloat(uuid_fields[5], 32) // power limit in W
	if err != nil {
		m.Error = fmt.Sprintf("Error parsing power limit: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
}

// Get total system CPU usage
func (m *MetricsState) getTotalCPUUsage() {
	totalCPUArray, err := cpu.Percent(0, false)
	if err != nil {
		m.Error = fmt.Sprintf("error getting total CPU usage: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.CpuPercentTotal = totalCPUArray[0]
}

// Get program-specific CPU usage
func (m *MetricsState) getProgramCPUUsage() {
	proc, err := process.NewProcess(int32(m.PID))
	if err != nil {
		m.Error = fmt.Sprintf("error getting process information: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	numCPU, err := cpu.Counts(true)
	if err != nil {
		m.Error = fmt.Sprintf("error getting CPU count: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		m.Error = fmt.Sprintf("error getting program CPU usage: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.CpuPercent = cpuPercent / float64(numCPU)

	m.RamPercent, err = proc.MemoryPercent()
	if err != nil {
		m.Error = fmt.Sprintf("error getting program RAM usage: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
}

// Get system RAM usage
func (m *MetricsState) getTotalRAMUsage() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		m.Error = fmt.Sprintf("error getting RAM usage: %v", err)
		log.Log.Warn("%s", m.Error)
		return
	}
	m.RamPercentTotal = vmStat.UsedPercent
}

func (m *MetricsState) postMetricsReset(c echo.Context) error {
	m.Error = ""
	m.ws.broadcastEngineState() // Use the instance to call the method
	return c.JSON(http.StatusOK, "")
}
