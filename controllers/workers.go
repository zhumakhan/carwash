package controllers
import(
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"carwashes/utils"
	"github.com/valyala/fasthttp"
//	"os"
	"fmt"
	"strconv"
)
func AddWorker(ctx *fasthttp.RequestCtx){
	worker := m.Worker{}
	if err := (&worker).UnmarshalJSON(ctx.PostBody()); err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	if len(worker.Name) == 0 || len(worker.Phone)  == 0 || worker.WashID == 0{
		respondWithMessage(ctx, fasthttp.StatusBadRequest, "Some fields are empty!")
		return
	}
	fmt.Println(worker)
	user := utils.Authorize(ctx)
	switch  u := user.(type){
	case m.Owner:
		if u.Role == 2{ 
			if u.WashID != worker.WashID{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your can't add worker")
				return
			}
		}else if u.Role == 1{
			if database.GetDB().Where(&m.CarWash{ID : worker.WashID, Owner : u.UUID}).Take(&m.CarWash{}).Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your can't add worker")
				return	
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your can't add worker")
				return
		}
		db := database.GetDB().Create(&worker)
		if db.Error != nil || db.RowsAffected == 0 {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Unknown error.")
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, worker)	
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
	}
}
func ChangeWorker(ctx *fasthttp.RequestCtx){
	worker := m.Worker{}
	if err := (&worker).UnmarshalJSON(ctx.PostBody()); err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(worker)
	user := utils.Authorize(ctx)
	if len(worker.Name) == 0 || len(worker.Phone)  == 0 || worker.WashID == 0{
		respondWithMessage(ctx, fasthttp.StatusBadRequest, "Some fields are empty!")
		return
	}
	switch  u := user.(type){
	case m.Owner:
		if u.Role == 2{ 
			if u.WashID != worker.WashID{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your can't add worker")
				return
			}
		}else if u.Role == 1{
			if database.GetDB().Where(&m.CarWash{ID : worker.WashID,Owner : u.UUID}).Take(&m.CarWash{}).Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your can't add worker")
				return	
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your can't add worker")
				return
		}
		idd,err := strconv.ParseUint(ctx.UserValue(c.WorkerID).(string), 10, 64)
		worker.ID = uint(idd)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return	
		}
		if database.GetDB().Save(&worker).Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Error occured while updating")
		}
		respondWithJSON(ctx, fasthttp.StatusOK, worker)	
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
	}	
}
func GetWorkers(ctx *fasthttp.RequestCtx){
	args := ctx.QueryArgs()
	user := utils.Authorize(ctx)
	washID,err := strconv.ParseUint(string(args.Peek("wash")), 10, 64)
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	switch  u := user.(type){
	case m.Owner:
		var workers m.Workers
		if u.Role == 1{
			if database.GetDB().Where(&m.CarWash{ID : uint(washID), Owner : u.UUID}).Take(&m.CarWash{}).Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not owner")
				return	
			}
		}else if u.Role == 2{
			if u.WashID != uint(washID){
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not owner")
				return
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
			return
		}
		database.GetDB().Where("wash_id = ?",uint(washID)).Find(&workers.Workers)
		respondWithJSON(ctx, fasthttp.StatusOK, workers)	
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
	}		
}
func GetWorker(ctx *fasthttp.RequestCtx){
	user := utils.Authorize(ctx)
	switch  u := user.(type){
	case m.Owner:
		idd,err := strconv.ParseUint(ctx.UserValue(c.WorkerID).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return	
		}
		worker:= m.Worker{ID : uint(idd)}
		db := database.GetDB().Where(&worker).Find(&worker)
		if db.Error != nil || db.RowsAffected == 0 {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Not found")
			return
		}
		if u.Role == 1{
			if database.GetDB().Where(&m.CarWash{ID : worker.WashID, Owner : u.UUID}).Take(&m.CarWash{}).Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
				return	
			}
		}else if u.Role == 2{
			if u.WashID == worker.WashID{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
				return	
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, worker)	
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
	}		
}
func DeleteWorker(ctx *fasthttp.RequestCtx){
	user := utils.Authorize(ctx)
	switch  u := user.(type){
	case m.Owner:
		idd,err := strconv.ParseUint(ctx.UserValue(c.WorkerID).(string), 10, 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return	
		}
		worker := m.Worker{ID : uint(idd)}
		if database.GetDB().Where(&worker).Find(&worker).Error != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Cannot find")
			return
		}
		if u.Role == 1{
			if database.GetDB().Where(&m.CarWash{ID : worker.WashID, Owner : u.UUID}).Take(&m.CarWash{}).Error != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your cannot delete")
				return	
			}
		}else if u.Role == 2{
			if u.WashID == worker.WashID{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your cannot delete")
				return	
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
			return
		}
		db := database.GetDB().Where(&worker).Delete(&worker)
		if db.Error != nil || db.RowsAffected == 0 {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Cannot delete")
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, worker)	
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
	}
}