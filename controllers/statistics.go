package controllers

import (
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"carwashes/utils"

	"github.com/valyala/fasthttp"
)

func GetStats(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch u := user.(type) {
	case m.Owner:
		if u.Role == 1 {
			washID := ctx.UserValue(c.WashIdPathVar).(string)
			start := string(ctx.QueryArgs().Peek("start"))
			end := string(ctx.QueryArgs().Peek("end"))
			// start, err := time.Parse("2006-01-02T15:04:05Z", startStr)
			// if err != nil {
			// 	respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			// 	return
			// }
			// end, err := time.Parse("2006-01-02T15:04:05Z", endStr)
			// if err != nil {
			// 	respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			// 	return
			// }
			db := database.GetPastDB()
			var bookings []m.PastBooking
			var query string
			if start == "" {
				if end != "" {
					query = "SELECT * FROM past_bookings_" + washID + " WHERE updated_at <= STR_TO_DATE('" + end + "', '%Y-%m-%dT%TZ')"
				} else {
					query = "SELECT * FROM past_bookings_" + washID
				}
			} else {
				if end != "" {
					query = "SELECT * FROM past_bookings_" + washID + " WHERE updated_at <= STR_TO_DATE('" + end + "', '%Y-%m-%dT%TZ') AND updated_at >= STR_TO_DATE('" + start + "', '%Y-%m-%dT%H:%i:%sZ')"
				} else {
					query = "SELECT * FROM past_bookings_" + washID + " WHERE updated_at >= STR_TO_DATE('" + start + "', '%Y-%m-%dT%TZ')"
				}
			}
			if err := db.Raw(query).Scan(&bookings); err.Error != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error.Error())
				return
			}
			var bookingsWrap PastBookingsResponse
			bookingsWrap.Bookings = bookings
			respondWithJSON(ctx, fasthttp.StatusOK, bookingsWrap)
		} else {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You have no permission")
			return
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
}
