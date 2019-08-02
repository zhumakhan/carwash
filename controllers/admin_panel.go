package controllers

import (
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func AdminGetOwners(ctx *fasthttp.RequestCtx) {
	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
	var ownercpy m.WashOwners
	db := database.GetDB()
	var owners []m.Owner
	db.Find(&owners)
	ownercpy.Owners = owners
	respondWithJSON(ctx, fasthttp.StatusOK, ownercpy)
}

func AdminGetOwner(ctx *fasthttp.RequestCtx) {
	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
	var uuid = ctx.UserValue(c.UUIDPathVar).(string)
	var washes []m.CarWash
	var owner m.Owner
	database.GetDB().Where(&m.CarWash{Owner: uuid}).Find(&washes)
	database.GetDB().Where(&m.Owner{UUID: uuid}).Find(&owner)
	var ownerfull m.OwnerWithWashes
	ownerfull.Owner = owner
	ownerfull.Washes = washes
	respondWithJSON(ctx, fasthttp.StatusOK, ownerfull)
}

func AdminInsertWash(ctx *fasthttp.RequestCtx) {
	wash := m.CarWash{}
	if err := (&wash).UnmarshalJSON(ctx.PostBody()); err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	err := getWashRegistrationErrors(wash)
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}

	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
	}
	wash.CreatedAt = time.Now().UTC().Truncate(time.Second)
	query := string(ctx.QueryArgs().Peek("months"))
	var months int64
	if query != "" {
		months, err = strconv.ParseInt(query, 10, 64)
	} else {
		months = 0
	}
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	wash.PaidUntil = wash.CreatedAt.AddDate(0, int(months), 0)
	db := database.GetDB().Create(&wash)
	if db.Error != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, "Unknown Error")
		return
	}
	//create archive booking table for wash
	m.WashId = wash.ID
	if db = database.GetPastDB().CreateTable(&m.PastBooking{}); db.Error != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, db.Error.Error())
		return
	}
	respondWithJSON(ctx, fasthttp.StatusOK, wash)
}

func AdminGetWashes(ctx *fasthttp.RequestCtx) {
	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
	var washes []m.CarWash
	var owners []m.Owner
	var ids []string

	database.GetDB().Find(&washes)
	database.GetDB().Find(&owners)
	for i := 0; i < len(washes); i++ {
		for j := 0; j < len(owners); j++ {
			if washes[i].Owner == owners[j].UUID {
				ids = append(ids, washes[i].Owner)
				washes[i].Owner = owners[j].FirstName + " " + owners[j].SecondName
				break
			}
		}
	}
	var nameWashes m.CarWashesWithOwnerName
	nameWashes.Washes = washes
	nameWashes.OwnerUUID = ids
	respondWithJSON(ctx, fasthttp.StatusOK, nameWashes)
}

func AdminGetWash(ctx *fasthttp.RequestCtx) {
	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
	}
	washId64, err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	washID := uint(washId64)
	var wash m.CarWash
	if err := database.GetDB().Where(&m.CarWash{ID: washID}).Find(&wash); err.Error != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
		return
	}
	respondWithJSON(ctx, fasthttp.StatusOK, wash)
}

func AdminGetPaymentHistory(ctx *fasthttp.RequestCtx) {
	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
	}
	washId64, err := strconv.ParseUint(ctx.UserValue(c.WashIdPathVar).(string), 10, 64)
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	washID := uint(washId64)
	var payments []m.MonthlyPaymentHistory
	if err := database.GetDB().Where(&m.MonthlyPaymentHistory{CarWash: washID}).Find(&payments); err.Error != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
		return
	}
	var wrapper m.Histories
	wrapper.History = payments
	respondWithJSON(ctx, fasthttp.StatusOK, wrapper)
}

func AdminProceedMonthlyPayment(ctx *fasthttp.RequestCtx) {
	pass := string(ctx.Request.Header.Peek("Password"))
	if pass != "qwertyui" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
	}
	paymentID := ctx.UserValue(c.UUIDPathVar).(string)
	newPayment := &m.MonthlyPaymentHistory{}
	if err := newPayment.UnmarshalJSON(ctx.PostBody()); err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	var payment m.MonthlyPaymentHistory
	if err := db.Where(&m.MonthlyPaymentHistory{UUID: paymentID}).Find(&payment); err.Error != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
		return
	}
	payment.Status = newPayment.Status
	if err := db.Save(&payment).Error; err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(ctx, fasthttp.StatusOK, payment)
}

// func AddAdmin(ctx *fasthttp.RequestCtx) {
// 	var admin m.Admin
// 	var owner m.Owner
// 	db := database.GetDB()
// 	admin.FirstName = "Allmight"
// 	admin.SecondName = "Hero"
// 	admin.Phone = "+7(777)666-66-66"
// 	admin.Password = "12345"
// 	admin.ConfirmPassword = "12345"
// 	admin.UUID = uuid.Must(uuid.NewV4(), nil).String()
// 	admin.CreatedAt = time.Now().UTC().Truncate(time.Second)
// 	db = db.Create(admin)
// 	admin.AccessToken, admin.RefreshToken = utils.GenerateTokens(admin.UUID)
// 	owner.AccessToken = admin.AccessToken
// 	respondWithJSON(ctx, fasthttp.StatusOK, owner)
// }
