package api

import "github.com/MathieuMoalic/amumax/src/engine"

type HeaderState struct {
	Path   string `msgpack:"path"`
	Status string `msgpack:"status"`
}

func initHeaderAPI() *HeaderState {
	status := ""
	if engine.Pause {
		status = "paused"
	} else {
		status = "running"
	}
	return &HeaderState{
		Path:   engine.OD(),
		Status: status,
	}
}

func (h *HeaderState) Update() {
	status := ""
	if engine.Pause {
		status = "paused"
	} else {
		status = "running"
	}
	h.Path = engine.OD()
	h.Status = status
}
