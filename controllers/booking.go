package controllers

import (
	"github.com/valyala/fasthttp"
	"carwashes/database"
	c "carwashes/constants"
	m "carwashes/models"
	"strconv"
	"carwashes/utils"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	"time"
)

//easyjson:json
type BookingsResponse struct {
	Bookings []m.Booking
}

//easyjson:json
type PastBookingsResponse struct {
	Bookings []m.PastBooking
}

func GetBookings(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	washStr := string(args.Peek("wash"))
	user := utils.Authorize(ctx)
	var washId uint64
	var err error
	if washStr != "" {
		washId, err = strconv.ParseUint(washStr, 10, 64)
	}
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	if string(args.Peek("past")) == "true" {
		switch u := user.(type) {
		case m.Owner:
			if u.Role == 1{
				if database.GetDB().Where(m.CarWash{ ID : uint(washId), Owner : u.UUID }).Find(&m.CarWash{}).RecordNotFound() {
					respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
					return
				}
			}else if u.Role == 2{
				if u.WashID  != uint(washId) || database.GetDB().Where(m.CarWash{ ID : uint(washId)}).Find(&m.CarWash{}).RecordNotFound() {
					respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
					return
				}
			}else{
				respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
				return
			}
		case m.Client:
			if u.UUID != string(args.Peek("client")) {
				respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
				return
			}
		default:
			respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
			return
		}
		bookings := PastBookingsResponse{}
		booking := m.PastBooking{
			CarWash: uint(washId),
			ClientUUID: string(args.Peek("client")),
		}
		db := database.GetPastDB().Order("updated_at desc")
		start := args.GetUintOrZero("start")
		if start != 0 {
			db = db.Offset(start)
		}
		end := args.GetUintOrZero("end")
		if end != 0 {
			db = db.Limit(end - start)
		}

		startDay := string(args.Peek("startDay"))
		if startDay != "" {
			startTime, err := time.Parse(time.RFC3339, startDay)
			if err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
			db = db.Where("updated_at > ?", startTime)
		}
		endDay := string(args.Peek("endDay"))
		if endDay != "" {
			endTime, err := time.Parse(time.RFC3339, endDay)
			if err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
			db = db.Where("updated_at < ?", endTime)
		}

		if booking.CarWash == 0 && booking.ClientUUID == "" {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Specify wash")
			return
		}
		m.WashId = uint(washId)
		if err := db.Table("past_bookings_" + strconv.FormatUint(washId, 10)).
			Where(&booking).
			Find(&(bookings.Bookings)).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, bookings)
	} else {
		switch u := user.(type) {
		case m.Owner:
			if u.Role == 1{
				if database.GetDB().Where(m.CarWash{ ID : uint(washId), Owner : u.UUID }).Find(&m.CarWash{}).RecordNotFound() {
					respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
					return
				}
			}else if u.Role == 2{
				if u.WashID  != uint(washId) || database.GetDB().Where(m.CarWash{ ID : uint(washId)}).Find(&m.CarWash{}).RecordNotFound() {
					respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
					return
				}
			}else{
				respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
				return
			}
		}
		db := database.GetDB().Order("order")
		bookings := BookingsResponse{}
		booking := m.Booking{
			CarWash: uint(washId),
			ClientUUID: string(args.Peek("client")),
		}
		if err := db.Preload("Client").Preload("BookingServices").Preload("BookingServices.Service").Where(&booking).Find(&(bookings.Bookings)).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		switch u := user.(type) {
		case m.Client:
			for _, b := range bookings.Bookings {
				if b.Client.UUID != u.UUID {
					 b.Client.UUID = ""
					 b.ClientUUID = ""
					 b.CarNumber = ""
				}
			}
		}

		respondWithJSON(ctx, fasthttp.StatusOK, bookings)
	}
}


func GetBooking(ctx *fasthttp.RequestCtx) {
	bookingUUID := ctx.UserValue(c.UUIDPathVar)

	db := database.GetDB()
	booking := m.Booking{
		UUID : bookingUUID.(string),
	}

	if booking.UUID == "" {
		respondWithError(ctx, fasthttp.StatusBadRequest, "Specify booking UUID")
		return
	}

	if err := db.Where(&booking).Preload("Client").Preload("BookingServices").Preload("BookingServices.Service").Take(&booking).Error; err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	user := utils.Authorize(ctx)
	switch u := user.(type) {
	case m.Owner:
		if u.Role == 1{
			if database.GetDB().Where(m.CarWash{ ID : booking.CarWash, Owner : u.UUID }).Find(&m.CarWash{}).RecordNotFound() {
				respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
				return
			}
		}else if u.Role == 2{
			if u.WashID != booking.CarWash || database.GetDB().Where(m.CarWash{ ID : booking.CarWash}).Find(&m.CarWash{}).RecordNotFound() {
				respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
				return
			}
		}else{
			respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
			return
		}
	case m.Client:
		if booking.Client.UUID != u.UUID {
			respondWithError(ctx, fasthttp.StatusUnauthorized, "Not authorized")
			return
		}
	}
	respondWithJSON(ctx, fasthttp.StatusOK, booking)
}

func DeleteBooking(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch owner := user.(type) {
	case m.Owner:
		bookingUUID := ctx.UserValue(c.UUIDPathVar)
		past := ctx.QueryArgs().Peek("past")

		db := database.GetDB()
		booking := m.Booking{
			UUID : bookingUUID.(string),
		}

		query := db.Where(&booking).Take(&booking)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "No such booking")
			return
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		var cw m.CarWash
		if owner.Role == 1{
			cw = m.CarWash{ ID : booking.CarWash, Owner : owner.UUID}	
		}else if owner.Role ==2{
			if owner.WashID == booking.CarWash{
				cw = m.CarWash{ ID : booking.CarWash}
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "You do not have enough privileges to delete bookings")
			return
		}
		query = db.Where(&cw).Take(&cw)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You are not the wash owner")
			return
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		if !utils.AuthorizeWash(cw) {
			respondWithError(ctx, fasthttp.StatusBadRequest, "This wash is inactive")
			return
		}

		if string(past) == "true" {
			query = db.Model(&booking).Preload("Client").Preload("BookingServices").Preload("BookingServices.Service").Where(&booking).Take(&booking)
			if query.Error != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
				return
			}
			booking.Order -= 1
			if err := db.Save(&booking).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
			servicesStr := ""
			for _, bs := range booking.BookingServices {
				servicesStr += bs.Service.Name + "; "
			}
			pastBooking := m.PastBooking{
				CarWash : booking.CarWash,
				ClientUUID : booking.ClientUUID,
				ClientFirstName : booking.Client.FirstName,
				ClientSecondName : booking.Client.SecondName,
				ClientMiddleName : booking.Client.MiddleName,
				ClientPhone : booking.Client.Phone,
				CarNumber : booking.CarNumber,
				Cost : booking.Cost,
				CreatedAt : booking.CreatedAt,
				UpdatedAt : booking.UpdatedAt,
				CarModel : booking.CarModel,
				BookingServices : servicesStr,
			}
			if err := database.GetPastDB().Create(&pastBooking).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
			//MonthlyPaymentHistory Start
			var payment m.MonthlyPaymentHistory
			year, mth, _ := booking.CreatedAt.Date()
			period := time.Date(year, mth, 1, 0, 0, 0, 0, time.UTC)
			if err := db.Where(&m.MonthlyPaymentHistory{CarWash: booking.CarWash, Month: period}).Find(&payment).Error; err != nil {
				payment.UUID = uuid.Must(uuid.NewV4(), nil).String()
				payment.CarWash = booking.CarWash
				payment.Amount = booking.Cost
				payment.Month = period
				if err := db.Create(&payment).Error; err != nil {
					respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
					return
				}
			} else {
				payment.Amount += booking.Cost
				if err := db.Save(&payment).Error; err != nil {
					respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
					return
				}
			}
		}

		bookingService := m.BookingService{ Booking: bookingUUID.(string) }
		bookingServices := []m.BookingService{}
		if err := db.Where(&bookingService).Find(&bookingServices).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		for _, bs := range bookingServices {
			if err := db.Where(&bs).Delete(&bs).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
		}

		if err := db.Delete(&booking).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		query = db.Model(&m.Booking{}).Where(m.Booking{ CarWash : booking.CarWash}).Order("`order` ASC").Where("`order` > ?", booking.Order).Update("order", gorm.Expr("`order` - ?", 1))
		if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, booking)
	default:
		respondWithError(ctx, fasthttp.StatusUnauthorized, "You are not authorized")
	}
}

func InsertBooking(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch owner := user.(type) {
	case m.Owner:
		newbooking := &m.Booking{}
		if err := newbooking.UnmarshalJSON(ctx.PostBody()); err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		washId, err := strconv.ParseUint(string(ctx.QueryArgs().Peek("wash")), 10, 64)
		if err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		if washId == 0 {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Wrong Wash")
			return
		}
		if newbooking.CarModel == "" {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Please include car model")
			return
		}
		db := database.GetDB()
		var wash m.CarWash
		if owner.Role == 1{
			wash = m.CarWash{ID : uint(washId),Owner : owner.UUID}	
		}else if owner.Role ==2{
			if owner.WashID == uint(washId){
				wash = m.CarWash{ID : uint(washId)}
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "You do not have enough privileges to delete bookings")
			return
		}
		query := db.Where(&wash).Take(&wash)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You are not the owner")
			return
		}
		if !utils.AuthorizeWash(wash) {
			respondWithError(ctx, fasthttp.StatusBadRequest, "This wash is inactive")
			return
		}
		query = db.Where(&newbooking.Client).Take(&newbooking.Client)
		if query.RecordNotFound() {
			newbooking.Client.UUID = uuid.Must(uuid.NewV4(), nil).String()
			if err := db.Create(&newbooking.Client).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
		}
		newbooking.ClientUUID = newbooking.Client.UUID
		newbooking.Client = m.Client{}
		if !db.Where(&m.Booking { CarNumber : newbooking.CarNumber}).Take(&m.Booking {}).RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "already exist")
			return
		}
		newbooking.CarWash = uint(washId)
		var order struct {
			Order uint
		}
		query = db.Table("bookings").Select("`order`").Where("car_wash = ?",uint(washId)).Order("`order` DESC").Limit(1).Scan(&order)
		if query.RecordNotFound() {
			order.Order = 0
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
			return
		}
		newbooking.Order = order.Order + 1
		if (newbooking.Cost == 0) {
			for i, s := range newbooking.BookingServices {
				if s.ServiceID == 0 {
					respondWithError(ctx, fasthttp.StatusBadRequest, "Wrong Service")
					return
				}
				sc := m.ServiceCost{ Service : s.ServiceID, CarModel : newbooking.CarModel }

				if err := db.Where(&sc).Take(&sc).Error; err != nil {
					respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
					return
				}
				newbooking.BookingServices[i].Cost = sc.Cost
				newbooking.BookingServices[i].Duration = sc.Duration
				newbooking.ServiceCosts = append(newbooking.ServiceCosts, sc)
				newbooking.Cost += sc.Cost
			}
		}
		newbooking.UUID = uuid.Must(uuid.NewV4(), nil).String()
		if err := db.Create(&newbooking).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(ctx, fasthttp.StatusOK, newbooking)
	default:
		respondWithError(ctx, fasthttp.StatusUnauthorized, "You are not authorized")
	}
}

func ChangeBooking(ctx *fasthttp.RequestCtx) {
	user := utils.Authorize(ctx)
	switch owner := user.(type) {
	case m.Owner:
		newbooking := &m.Booking{}
		if err := newbooking.UnmarshalJSON(ctx.PostBody()); err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		bookingUUID := ctx.UserValue(c.UUIDPathVar).(string)
		if newbooking.CarModel == "" {
			respondWithError(ctx, fasthttp.StatusBadRequest, "Please include car model")
			return
		}
		db := database.GetDB()
		booking := m.Booking{ UUID : bookingUUID }
		if err := db.Where(&booking).Take(&booking).Error; err != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
			return
		}
		var wash m.CarWash
		if owner.Role == 1{
			wash = m.CarWash{ID : booking.CarWash,Owner : owner.UUID}	
		}else if owner.Role ==2{
			if owner.WashID == booking.CarWash{
				wash = m.CarWash{ID : booking.CarWash}
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "You do not have enough privileges to delete bookings")
			return
		}
		query := db.Where(&wash).Take(&wash)
		if query.RecordNotFound() {
			respondWithError(ctx, fasthttp.StatusBadRequest, "You are not the owner")
			return
		}
		if !utils.AuthorizeWash(wash) {
			respondWithError(ctx, fasthttp.StatusBadRequest, "This wash is inactive")
			return
		}
		query = db.Where(&newbooking.Client).Take(&newbooking.Client)
		if query.RecordNotFound() {
			newbooking.Client.UUID = uuid.Must(uuid.NewV4(), nil).String()
			if err := db.Create(&newbooking.Client).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
		} else if query.Error != nil {
			respondWithError(ctx, fasthttp.StatusBadRequest, query.Error.Error())
		}
		oldClientUUID := booking.ClientUUID
		booking.ClientUUID = newbooking.Client.UUID
		booking.Client = newbooking.Client
		db.Where(&m.BookingService{ Booking : booking.UUID }).Find(&booking.BookingServices)
		for _, s := range booking.BookingServices {
			if err := db.Where(&s).Delete(&s).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
		}
		booking.BookingServices = nil
		booking.Cost = 0
		for _, s := range newbooking.BookingServices {
			if s.ServiceID == 0 {
				respondWithError(ctx, fasthttp.StatusBadRequest, "Wrong Service")
				return
			}
			sc := m.ServiceCost{ Service : s.ServiceID, CarModel : newbooking.CarModel }

			if err := db.Where(&sc).Take(&sc).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
			s.Booking = booking.UUID
			s.Cost = sc.Cost
			s.Duration = sc.Duration
			booking.BookingServices = append(booking.BookingServices, s)
			booking.ServiceCosts = append(booking.ServiceCosts, sc)
			booking.Cost += sc.Cost
			if err := db.FirstOrCreate(&s).Error; err != nil {
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
		}
		if (newbooking.Cost != 0) {
			booking.Cost = newbooking.Cost
		}
		booking.CarNumber = newbooking.CarNumber
		booking.PaymentStatus = newbooking.PaymentStatus
		booking.CarModel = newbooking.CarModel
		db.Save(&booking)
		db.Where(&m.Client{ UUID : oldClientUUID }).Delete(&m.Client{ UUID : oldClientUUID })
		respondWithJSON(ctx, fasthttp.StatusOK, booking)
	default:
		respondWithError(ctx, fasthttp.StatusUnauthorized, "You are not authorized")
	}
}
