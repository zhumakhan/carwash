package controllers

import (
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"carwashes/utils"
	"github.com/valyala/fasthttp"
	"strconv"
)

func AddCarType(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch userT := user.(type){
	case m.Owner:
		washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		washId := uint(washId64)
		if userT.Role == 1{
			if err := database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Find(&m.CarWash{}); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
		}else if userT.Role == 2{
			if err := database.GetDB().Where(&m.CarWash{ID : userT.WashID}).Find(&m.CarWash{}); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
		}
		var carType m.CarType
		(&carType).UnmarshalJSON(ctx.PostBody())

		if err := database.GetDB().Create(&m.CarType{CarWashID: washId, Name: carType.Name}); err.Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
		respondWithMessage(ctx, fasthttp.StatusOK, "OK")
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
}

//easyjson:json
type CarTypes struct{
	CarTypes []m.CarType
}
func GetCarTypes(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch userT := user.(type){
	case m.Owner:
		washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		washId := uint(washId64)
		var carTypes CarTypes
		if userT.Role == 1{
			if err := database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Find(&m.CarWash{}); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
			if err := database.GetDB().Where(&m.CarType{CarWashID : washId}).Find(&(carTypes.CarTypes)); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}	
		}else if userT.Role == 2{
			if err := database.GetDB().Where(&m.CarWash{ID : userT.WashID}).Find(&m.CarWash{}); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
			if err := database.GetDB().Where(&m.CarType{CarWashID : userT.WashID}).Find(&(carTypes.CarTypes)); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
		}
		
		respondWithJSON(ctx, fasthttp.StatusOK, carTypes)
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
}

func DeleteCarType(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch userT := user.(type){
	case m.Owner:
		washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		washId := uint(washId64)
		carTypeId64,err := strconv.ParseUint(ctx.UserValue(c.CarTypeIdPathVar).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		carTypeId := uint(carTypeId64)
		if userT.Role == 1{
			if err := database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Find(&m.CarWash{}); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}	
		}else if userT.Role == 2{
			if err := database.GetDB().Where(&m.CarWash{ID : userT.WashID}).Find(&m.CarWash{}); err.Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
		}
		
		carType := m.CarType{ID : carTypeId, CarWashID: washId}
		if err := database.GetDB().Where(&carType).Delete(&carType); err.Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
		respondWithMessage(ctx, fasthttp.StatusOK, "OK")
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
}

func ChangeCarType(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch userT := user.(type){
	case m.Owner:
		washId64,err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		washId := uint(washId64)
		carTypeId64,err := strconv.ParseUint(ctx.UserValue(c.CarTypeIdPathVar).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		carTypeId := uint(carTypeId64)

		if err := database.GetDB().Where(&m.CarWash{ID : washId, Owner : userT.UUID}).Find(&m.CarWash{}); err.Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}

		var carType m.CarType
		(&carType).UnmarshalJSON(ctx.PostBody())
		carType.ID = carTypeId
		carType.CarWashID = washId
		if err := database.GetDB().Model(&carType).Update(&carType); err.Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, carType)
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
}
