package controllers

import (
	"github.com/valyala/fasthttp"
	"carwashes/database"
	c "carwashes/constants"
	m "carwashes/models"
	"strconv"
	"carwashes/utils"
//	
)

//easyjson:json
type ServicesResponse struct {
	Services []m.Service
}

func GetWashServices(ctx *fasthttp.RequestCtx) {
	washIdStr := ctx.UserValue(c.WashIdPathVar)
	washId, err := strconv.ParseUint(washIdStr.(string), 10, 64)
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	service := m.Service{ CarWash: uint(washId) }
	services := ServicesResponse{}
	if err := db.Where(&service).Find(&(services.Services)).Error; err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	for i, s := range services.Services {
		if err := db.Where(&m.ServiceCost{ Service: s.ID }).Find(&(services.Services[i].Costs)).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(ctx, fasthttp.StatusOK, services)
}

func GetWashService(ctx *fasthttp.RequestCtx) {
	washIdStr := ctx.UserValue(c.WashIdPathVar)
	serviceIdStr := ctx.UserValue(c.ServiceIdPathVar)
	washId, err := strconv.ParseUint(washIdStr.(string), 10, 64)
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	serviceId, err := strconv.ParseUint(serviceIdStr.(string), 10, 64)
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	service := m.Service{ CarWash: uint(washId), ID: uint(serviceId) }
	if err := db.Where(&service).Take(&service).Error; err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	if err := db.Where(&m.ServiceCost{ Service: service.ID }).Find(&(service.Costs)).Error; err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(ctx, fasthttp.StatusOK, service)
}

func InsertWashService(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch owner := user.(type) {
	case m.Owner:
		service := m.Service{}
		body := ctx.PostBody()
		if err := service.UnmarshalJSON(body); err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}

		service.ID = 0

		db := database.GetDB()
		cw := m.CarWash{ Owner : owner.UUID, ID : service.CarWash }
		query := db.Where(&cw).Find(&cw)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You do not own this wash")
			return
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		serviceFind := service
		serviceFind.Costs = nil
		if !db.Model(&service).Where(&serviceFind).Find(&serviceFind).RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Already Exist")
			return
		}
		if err := db.Create(&service).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, service)
	default:
		respondWithError(ctx, fasthttp.StatusUnauthorized, "You are not authorized")
	}
}

func ChangeWashService(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch owner := user.(type) {
	case m.Owner:
		washIdStr := ctx.UserValue(c.WashIdPathVar)
		washId, err := strconv.ParseUint(washIdStr.(string), 10, 64)
		if err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		serviceIdStr := ctx.UserValue(c.ServiceIdPathVar)
		serviceId, err := strconv.ParseUint(serviceIdStr.(string), 10, 64)
		if err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		service := &m.Service{
			ID : uint(serviceId),
			CarWash : uint(washId),
		}
		db := database.GetDB()
		if err := db.Where(&service).Take(&service).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}

		cw := m.CarWash{ Owner : owner.UUID, ID : uint(washId) }
		query := db.Where(&cw).Find(&cw)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You do not own this wash")
			return
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}

		body := ctx.PostBody()
		if err := service.UnmarshalJSON(body); err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		service.ID = uint(serviceId)
		service.CarWash = uint(washId)
		query = db.Where(&m.ServiceCost{ Service : uint(serviceId) }).Delete(&m.ServiceCost{})
		if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}

		for i, sc := range service.Costs {
			if sc.CarModel == "" {
				respondWithError(ctx, fasthttp.StatusBadRequest, "Wrong Model")
				return
			}
			query = db.Where(&m.ServiceCost{ Service : uint(serviceId), CarModel : sc.CarModel }).FirstOrCreate(&sc)
			if query.Error != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
				return
			}
			sc.Cost, sc.Duration = service.Costs[i].Cost, service.Costs[i].Duration
			db.Save(&sc)
		}
		service.Costs = nil
		if err := db.Save(&service).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, service)
	default:
		respondWithError(ctx, fasthttp.StatusUnauthorized, "You are not authorized")
	}
}

func DeleteWashService(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch owner := user.(type) {
	case m.Owner:
		washIdStr := ctx.UserValue(c.WashIdPathVar)
		washId, err := strconv.ParseUint(washIdStr.(string), 10, 64)
		if err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		service := &m.Service{}

		service.CarWash = uint(washId)
		var servId uint64
		servId, err = strconv.ParseUint(ctx.UserValue(c.ServiceIdPathVar).(string), 10, 64)
		service.ID = uint(servId)

		db := database.GetDB()
		cw := m.CarWash{ Owner : owner.UUID, ID : uint(washId) }
		query := db.Where(&cw).Find(&cw)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You do not own this wash")
			return
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		query = db.Where(&m.ServiceCost{ Service : uint(servId) }).Delete(&m.ServiceCost{})
		if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		if err := db.Where(&service).Delete(&service).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, service)
	default:
		respondWithError(ctx, fasthttp.StatusUnauthorized, "You are not authorized")
	}
}
