package api

import "github.com/MathieuMoalic/amumax/engine"

type Header struct {
	Path     string  `json:"path"`
	Progress float64 `json:"progress"`
	Status   string  `json:"status"`
}

func newHeader() *Header {
	status := ""
	if engine.Pause {
		status = "paused"
	} else {
		status = "running"

	}
	return &Header{
		Path:   engine.OD(),
		Status: status,
	}
}
