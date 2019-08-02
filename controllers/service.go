package controllers

import (
	"github.com/valyala/fasthttp"
)

func GetStatus(ctx *fasthttp.RequestCtx)  {
	respondWithMessage(ctx, fasthttp.StatusOK, "OK")
}