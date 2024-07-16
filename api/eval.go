package api

import (
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
)

func EvalGUI(code string) {
	defer func() {
		if err := recover(); err != nil {
			if userErr, ok := err.(engine.UserErr); ok {
				util.LogErr(userErr)
			} else {
				panic(err)
			}
		}
	}()
	engine.Eval(code)
}
