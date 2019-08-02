package controllers

import (
	"github.com/valyala/fasthttp"
	"fmt"

)

type PanicMonger struct {
	onPanic []func()
}

func (this *PanicMonger) Catch(ctx *fasthttp.RequestCtx) {
	if err := recover(); err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, fmt.Sprint(err))
		for _, f := range this.onPanic {
			f()
		}
	}
}

func (this *PanicMonger) Append(f func()) {
	this.onPanic = append(this.onPanic, f)
}

func (this *PanicMonger) Prepend(f func()) {
	this.onPanic = append([]func() { f }, this.onPanic...)
}
