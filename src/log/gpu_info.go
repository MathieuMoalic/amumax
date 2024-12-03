package log

// Used in logs and metadata
type GpuInfo struct {
	CudaVersion   string
	CUDACC        string
	DevName       string
	TotalMem      string
	DriverVersion string
	GPUCC         string
}
