package controllers

import (
	"github.com/valyala/fasthttp"

	"carwashes/utils"
	m "carwashes/models"
)

type easyjsonInterface interface {
	MarshalJSON() ([]byte, error)
}

//easyjson:json
type responseMessage struct {
	Message string
}
//easyjson:json
type errorMessage struct {
	Error string
}

func respondWithMessage(ctx *fasthttp.RequestCtx, statusCode int, resp string) {
	respondWithJSON(ctx, statusCode, responseMessage{resp})
}

func respondWithError(ctx *fasthttp.RequestCtx, statusCode int, err string) {
	respondWithJSON(ctx, statusCode, errorMessage{err})
}

func respondWithJSON(ctx *fasthttp.RequestCtx, statusCode int, payload easyjsonInterface) {
	ctx.Response.Header.SetStatusCode(statusCode)
	wb, err := payload.MarshalJSON()
	if err != nil {
		ctx.SetContentType("text/plain")
		ctx.WriteString(err.Error())
	} else {
		ctx.SetContentType("application/json")
		ctx.Write(wb)
	}
}

func Hello(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch v := user.(type) {
		case m.Owner:
			respondWithMessage(ctx, 200, "Hello, Owner! Your email is " +  v.Email)
		case m.Client:
			respondWithMessage(ctx, 200, "Hello, Client!")
		default:
			respondWithMessage(ctx, 200, "Hello, Anonymous!")
	}
}
