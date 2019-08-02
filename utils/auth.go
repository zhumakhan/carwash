package utils

import (
		"github.com/valyala/fasthttp"
		m "carwashes/models"
		"os"
		"carwashes/database"
		"time"
)

func Authorize(ctx *fasthttp.RequestCtx) interface{} {
	access := ctx.Request.Header.Peek("Authorization")

	err, tc := ValidateToken(string(access), os.Getenv("access_token_password"))
	if err != nil {
		return nil
	}
	if err == nil {
		owner := m.Owner{UUID: tc.UUID}
		user  := m.User{UUID : tc.UUID}
		admin := m.Admin{UUID: tc.UUID}
		if !database.GetDB().Take(&owner).RecordNotFound() {
			return owner
		}
		if !database.GetDB().Table("turaQshare.users").Where(&m.User{UUID : user.UUID}).Scan(&user).RecordNotFound() {
			return user
		}
		if !database.GetDB().Where(&admin).Take(&admin).RecordNotFound() {
			return admin
		}
		return nil
	}
	return nil
}

func AuthorizeWash(wash m.CarWash) bool{
	present := time.Now().UTC().Truncate(time.Second)

	if present.Before(wash.PaidUntil) {
		return true
	}
	return false
}
