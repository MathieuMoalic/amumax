package api

import "github.com/MathieuMoalic/amumax/engine"

type Header struct {
	Path   string `msgpack:"path"`
	Status string `msgpack:"status"`
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
