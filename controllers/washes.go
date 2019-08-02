package controllers

import (
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"carwashes/utils"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)
//TODO:save WWWwith photo in server
func InsertWash(ctx *fasthttp.RequestCtx){
	wash := m.CarWash{}

	if err := (&wash).UnmarshalJSON(ctx.PostBody()); err != nil {
		respondWithError(ctx,fasthttp.StatusBadRequest, err.Error())
		return
	}
	err := getWashRegistrationErrors(wash)
	if  err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	user := utils.Authorize(ctx)
	switch userT := user.(type) {
	case m.Owner:
		wash.Owner = userT.UUID
	case m.Admin:
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
	wash.CreatedAt = time.Now().UTC().Truncate(time.Second)
	wash.PaidUntil = wash.CreatedAt	//tarif
	db := database.GetDB().Create(&wash)
	if db.Error != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, "Unknown Error")
		return
	}
	//create archive booking table for wash
	m.WashId = wash.ID
	if db = database.GetPastDB().CreateTable(&m.PastBooking{}); db.Error != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, db.Error.Error())
		return
	}
	respondWithJSON(ctx, fasthttp.StatusOK, wash)
}
/*
washes
Authorization=accestoken,
if Authorization argument is empty, then return all washes
for turaQ admins:
washes?approved=true to return status 1 washes
washes?approved=false or washes/:wash-id to return not approved washes
*/
func GetWashes(ctx *fasthttp.RequestCtx){
	var washes []m.CarWash
	var washescpy m.CarWashes
	status := uint(0)
	if string(ctx.QueryArgs().Peek("approved")) == "true"{
		status = 1
	}
	user := utils.Authorize(ctx)
	if user == nil{
		database.GetDB().Preload("Services").Find(&washes)
		for i := 0; i < len(washes); i++ {
			washes[i].Owner = ""
			var books []m.Booking
			database.GetDB().Where(m.Booking{CarWash : washes[i].ID}).Take(&books)
			max := uint(0)
			for j := 0; j < len(books); j++{
				if max < books[j].Order{
					max = books[j].Order
				}
			}	
			washes[i].QueueSize = max
		}
	}else{
		switch userT := user.(type){
		case m.Owner:
			if userT.Role == 2{
				database.GetDB().Where(&m.CarWash{ID : userT.WashID}).Preload("Services").Preload("Services.Costs").Preload("CarTypes").Find(&washes)
			}else if userT.Role == 1{
				database.GetDB().Where(&m.CarWash{Owner : userT.UUID}).Preload("Services").Preload("Services.Costs").Preload("CarTypes").Find(&washes)
			}
		case m.Admin:
			database.GetDB().Where(&m.CarWash{Status : status}).Preload("Services").Preload("Services.Costs").Preload("CarTypes").Find(&washes)
		default:
			respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
			return
		}
	}
	washescpy.Washes = washes
	respondWithJSON(ctx, fasthttp.StatusOK, washescpy)
}
/*
*washes/:wash-id
*/
func GetWash(ctx *fasthttp.RequestCtx){
	var wash m.CarWash


	user := utils.Authorize(ctx)
	washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
	washId := uint(washId64)
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	if washId == 0{
		respondWithError(ctx, fasthttp.StatusBadRequest, "RecordNotFound")
		return
	}
	switch userT := user.(type){
	case m.Owner:
		if database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Preload("Services").Preload("Services.Costs").Preload("CarTypes").Find(&wash).RecordNotFound(){
			respondWithError(ctx, fasthttp.StatusBadRequest, "RecordNotFound")
			return
		}
	case m.Admin:
		if database.GetDB().Where(&m.CarWash{ID : washId}).Preload("Services").Preload("CarTypes").Find(&wash).RecordNotFound(){
			respondWithError(ctx, fasthttp.StatusBadRequest, "RecordNotFound")
			return
		}
	case m.User:
		if database.GetDB().Where(&m.CarWash{ID : washId}).Preload("Services").Preload("CarTypes").Find(&wash).RecordNotFound(){
			respondWithError(ctx, fasthttp.StatusBadRequest, "RecordNotFound")
			return
		}
		wash.ID = 0
		wash.Status = 0
		wash.Longitude = 0.0
		wash.Latitude = 0.0
		wash.Owner = ""
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
	}
//may require to return with bookings, services, history for today
	respondWithJSON(ctx, fasthttp.StatusOK, wash)
}
func UpdateWash(ctx *fasthttp.RequestCtx){
	var wash, newwash m.CarWash
	washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	washId := uint(washId64)
	if  err := (&newwash).UnmarshalJSON(ctx.PostBody()); err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	user := utils.Authorize(ctx)
	switch userT := user.(type){
	case m.Owner:
		if err := database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Find(&wash); err.Error!= nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
	case m.Admin:
		if err := database.GetDB().Where(&m.CarWash{ID : washId}).Find(&wash); err.Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}

	if len(newwash.Name) > 0 {
		wash.Name = newwash.Name
	}
	if len(newwash.Address) > 0{
		wash.Address= newwash.Address
	}
	if newwash.Longitude != 0{
		wash.Longitude = newwash.Longitude
	}
	if newwash.Latitude != 0{
		wash.Latitude = newwash.Latitude
	}
	if len(newwash.Owner) > 0{
		wash.Owner = newwash.Owner
	}
	if len(newwash.Photo) > 0{
		wash.Photo = newwash.Photo
	}
	//Tarifs
	if ctx.QueryArgs().Peek("period") != nil {
		months64, err := strconv.ParseInt(string(ctx.QueryArgs().Peek("period")), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		months := int(months64)
		present := time.Now().UTC().Truncate(time.Second)
		if present.Before(wash.PaidUntil) {
			wash.PaidUntil = wash.PaidUntil.AddDate(0, months, 0)
		} else {
			wash.PaidUntil = time.Now().AddDate(0, months, 0).UTC().Truncate(time.Second)
		}
	}
	//
	if err := database.GetDB().Save(&wash); err.Error != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
		return
	}
	respondWithMessage(ctx, fasthttp.StatusOK, "OK")
}
func DeleteWash(ctx *fasthttp.RequestCtx){
	var wash m.CarWash
	washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	washId := uint(washId64)
	user := utils.Authorize(ctx)
	switch userT := user.(type){
	case m.Owner:
		if err := database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Find(&wash); err.Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
	case m.Admin:
		if err := database.GetDB().Where(&m.CarWash{ID : washId}).Find(&wash); err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
	if err := database.GetDB().Delete(&wash); err.Error != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
		return
	}
	respondWithMessage(ctx, fasthttp.StatusOK, "OK")
}
