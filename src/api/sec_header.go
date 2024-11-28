package api

import "github.com/MathieuMoalic/amumax/src/engine_old"

type HeaderState struct {
	Path    string  `msgpack:"path"`
	Status  string  `msgpack:"status"`
	Version *string `msgpack:"version"`
}

func initHeaderAPI() *HeaderState {
	status := ""
	if engine_old.Pause {
		status = "paused"
	} else {
		status = "running"
	}
	return &HeaderState{
		Path:    engine_old.OD(),
		Status:  status,
		Version: &engine_old.VERSION,
	}
}

func (h *HeaderState) Update() {
	status := ""
	if engine_old.Pause {
		status = "paused"
	} else {
		status = "running"
	}
	h.Path = engine_old.OD()
	h.Status = status
}
