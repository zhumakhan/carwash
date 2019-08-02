package controllers

import (
	"carwashes/database"
	m "carwashes/models"
	"carwashes/utils"
	"github.com/valyala/fasthttp"
	"strconv"
)

//easyjson:json
type Clients struct {
	Clients []m.Client
}

//TODO: FULLTEXT KEYS AND MATCH SEARCH

func GetClients(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	washStr := string(args.Peek("wash"))
	user := utils.Authorize(ctx)
	var washId uint64
	var err error
	if washStr != "" {
		washId, err = strconv.ParseUint(washStr, 10, 64)
	} else {
		respondWithError(ctx, fasthttp.StatusBadRequest, "specify wash")
		return
	}
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	switch userT := user.(type){
	case m.Owner:
		db := database.GetDB()
		if err := db.Where(&m.CarWash{ID : uint(washId), Owner : userT.UUID}).Find(&m.CarWash{}).Error; err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		var clients Clients
		if err := db.Model(&m.Booking{}).Joins("LEFT JOIN bookings ON bookings.client_uuid = clients.uuid").
			Where(&m.Booking{ CarWash : uint(washId), CarNumber : string(args.Peek("car-number")) }).
			Find(&clients.Clients).Error;
			err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, clients)
		return
	}
}
