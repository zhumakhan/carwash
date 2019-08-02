package controllers

import (
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"carwashes/utils"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"strconv"
//	"os"
	"errors"
	"fmt"
	"time"
//	"strings"
)


func Register(ctx *fasthttp.RequestCtx) {
	//Reading $User struct from form
	args := ctx.QueryArgs()
	owner := m.Owner{}
	err := (&owner).UnmarshalJSON(ctx.PostBody())
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	owner.UUID = uuid.Must(uuid.NewV4(), nil).String()
	owner.CreatedAt = time.Now().UTC().Truncate(time.Second)
	// Checking registering user for errors
	// If has errors, respond $UserError struct with errors description
	if string(args.Peek("role")) == "manager"{
		user := utils.Authorize(ctx)
		switch  user.(type){
		case m.Owner:
			owner.Role = 2
			id64, err := strconv.ParseUint(string(args.Peek("wash")), 10, 64)
			if err != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return	
			}
			owner.WashID = uint(id64)
		default:
			respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
			return
		}
	} 
	ownerE := getUserRegistrationErrors(owner)
	if !isOwnerEmpty(ownerE) {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, ownerE)
		return
	}
	// If everything is ok, add to database
	// If gets errors with adding to database, respond as "unknown error"
	/*
	if db.Where(&m.CarWash{Owner : owner.UUID}).Take(&wash).RecordNotFound(){
		respondWithError(ctx, fasthttp.StatusBadRequest, fmt.Sprintf("Wash with id %d does not exist!"))
		return
	}
	*/
	fmt.Println(owner)
	db = db.Create(&owner)
	if db.Error != nil || db.RowsAffected == 0 {
		respondWithError(ctx, fasthttp.StatusBadRequest, "Unknown error.")
		return
	}
	// Create access and refresh tokens
	// Respond $User struct and success message
	owner.AccessToken, owner.RefreshToken = utils.GenerateTokens(owner.UUID)
//	Subscribe(strings.Fields(owner.DeviceToken))
	err = utils.MailWelcome(owner.Email, owner.FirstName)
	owner.Password = ""
	owner.ConfirmPassword = ""
	respondWithJSON(ctx, fasthttp.StatusOK, owner)
}
func GetManagers(ctx *fasthttp.RequestCtx){
	user := utils.Authorize(ctx)
	var washes []m.CarWash
	var managers []m.Owner
	switch  u := user.(type){
	case m.Owner:
		if  u.Role == 1{
			database.GetDB().Where(&m.CarWash{Owner : u.UUID}).Find(&washes)
			var washIDs []uint
			for i := 0; i < len(washes); i++{
				washIDs = append(washIDs, washes[i].ID)
			}
			database.GetDB().Where("wash_id IN (?)", washIDs).Find(&managers)
			washOwners := m.WashOwners{Owners:managers}
			respondWithJSON(ctx, fasthttp.StatusOK, washOwners)
			return
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")	
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
	}
}
func GetManager(ctx *fasthttp.RequestCtx){
	managerUUID := ctx.UserValue(c.UUIDPathVar)
	user := utils.Authorize(ctx)
	switch u := user.(type){
	case m.Owner:
		if u.Role == 1{
			manager := m.Owner{UUID : managerUUID.(string)}
			err := database.GetDB().Where(&manager).Take(&manager).Error
			if err != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest,err.Error())
			}else{
				wash := m.CarWash{ID : manager.WashID}
				database.GetDB().Where(&wash).Take(&wash)
				if wash.Owner != u.UUID{
					respondWithError(ctx, fasthttp.StatusBadRequest, "Manager not found")
					return	
				}
				respondWithJSON(ctx, fasthttp.StatusOK,manager)
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")	
	}
}
func PutManager(ctx *fasthttp.RequestCtx){
	user := utils.Authorize(ctx)
	switch u := user.(type){
	case m.Owner:
		if  u.Role == 1{
			updatedManager := m.Owner{}
			err := (&updatedManager).UnmarshalJSON(ctx.PostBody())
			if err != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return	
			}
			updatedManager.UUID = ctx.UserValue(c.UUIDPathVar).(string)
			wash := m.CarWash{ID : updatedManager.WashID}
			database.GetDB().Where(&wash).Take(&wash)
			if wash.Owner != u.UUID{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Your not owner of that wash")
				return	
			}
			err = database.GetDB().Save(&updatedManager).Error
			if err != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
				return
			}
			respondWithJSON(ctx, fasthttp.StatusOK, updatedManager)	
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized")	
	}	
}
func DeleteManager(ctx *fasthttp.RequestCtx){
	user := utils.Authorize(ctx)
	switch u := user.(type){
	case m.Owner:
		if u.Role == 1{
			manager := m.Owner{UUID : ctx.UserValue(c.UUIDPathVar).(string)}
			database.GetDB().Where(&manager).Take(&manager)
			wash := m.CarWash{ID : manager.WashID}
			database.GetDB().Where(&wash).Take(&wash)
			if wash.Owner != u.UUID{
				respondWithError(ctx, fasthttp.StatusBadRequest, "Manager not found")
				return	
			}
			err := database.GetDB().Where(&manager).Delete(&m.Owner{}).Error
			if err != nil{
				respondWithError(ctx, fasthttp.StatusBadRequest,err.Error())
				return
			}else{
				respondWithMessage(ctx, fasthttp.StatusOK,"Deleted")
				return
			}
		}else{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized1")
			return
		}
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "Your are not authorized2")	
	}	
}

//TODO:save with photo
func WashRegister(ctx *fasthttp.RequestCtx){
	wash := m.CarWash{}
	err := (&wash).UnmarshalJSON(ctx.PostBody())
	if err != nil {
		respondWithError(ctx,fasthttp.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	now := time.Now().UTC().Truncate(time.Second)
	wash.CreatedAt = now;
	// Checking registering user for errors
	// If has errors, respond $UserError struct with errors description
	err = getWashRegistrationErrors(wash)
	if  err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	// If everything is ok, add to database
	// If gets errors with adding to database, respond as "unknown error"
	db = db.Create(&wash)
	if db.Error != nil || db.RowsAffected == 0 {
		respondWithError(ctx, fasthttp.StatusBadRequest, "Unknown error.")
		return
	}
	respondWithJSON(ctx, fasthttp.StatusOK, wash)
}
// $ROOT/api/v1/login
func LoginOwner(ctx *fasthttp.RequestCtx) {
	//Parsing form
	owner := m.Owner{}
	err := (&owner).UnmarshalJSON(ctx.PostBody())
	if err != nil {
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	// Checking login user for errors
	// If has errors, respond $UserError struct with errors description
	ownerE := getUserLoginErrors(owner)
	if !isOwnerEmpty(ownerE) {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, ownerE)
		return
	}
	// Checking login user for existence in database and password correctness
	if db.Where(&m.Owner{Phone: owner.Phone, Password: owner.Password}).Take(&owner).RecordNotFound() {
		respondWithError(ctx, fasthttp.StatusBadRequest, "Wrong phone number or password.")
		return
	}
	// Create access and refresh tokens
	// Respond $User struct and success message
	//add device token
	owner.AccessToken, owner.RefreshToken = utils.GenerateTokens(owner.UUID)
	db.Save(&owner)
	//Add user to topic 'all'
	//Subscribe(strings.Fields(deviceToken))
	respondWithJSON(ctx, fasthttp.StatusOK, owner)
}

/*
// $ROOT/api/v1/refresh
func Refresh(ctx *fasthttp.RequestCtx) {
	//http://127.0.0.1:8080/path?refresh=sdkjalsjdlajdlkjk
	refreshTokenStr := ctx.QueryArgs().Peek("refresh") // getting refresh token as  query args
	//Checking is token valid

	isOk, tc := utils.ValidateToken(string(refreshTokenStr), os.Getenv("refresh_token_password"))
	if !isOk || database.GetDB().Where(&m.Owner{UUID: tc.UUID}).Take(&m.Owner{}).RecordNotFound() {
		respondWithError(ctx, fasthttp.StatusBadRequest, c.TokenInvalid)
		return
	}

	//Creating access token with 1 day duration
	now := time.Now().UTC()
	issuedAt := now.Unix()
	accessExpiresAt := now.AddDate(0, 0, 1).Unix()
	accessTC := &utils.TokenClaim{UUID: tc.UUID, IssuedAt: issuedAt, ExpiresAt: accessExpiresAt}
	accessTokenStr, _ := utils.GenerateToken(accessTC, os.Getenv("access_token_password"))
	// Responding access token
	respondWithMessage(ctx, fasthttp.StatusOK, accessTokenStr)
}
*/
// Checks user for empty fields, correctness of email and phone number, uniqueness in database
// Sets error texts in $UserError structure fields
func getUserRegistrationErrors(owner m.Owner) m.OwnerError {
	ownerE := m.OwnerError{}
	db := database.GetDB()

	if owner.FirstName == "" {
		ownerE.FirstName =   c.Empty
	}
	if owner.Phone == "" {
		ownerE.Phone =  c.Empty
	}
	if owner.Password == "" {
		ownerE.Password =   c.Empty
	}
	if owner.ConfirmPassword == "" {
		ownerE.ConfirmPassword =   c.Empty
	}

	if owner.Email != "" && !utils.IsEmail(owner.Email) {
		ownerE.Email = c.WrongEmail
	}

	if ownerE.Phone == "" && !utils.IsPhoneNumber(owner.Phone) {
		ownerE.Phone = c.WrongPhone
	}

	if ownerE.Phone == "" && !db.Where(&m.Owner{Phone: owner.Phone}).Take(&m.Owner{}).RecordNotFound() {
		ownerE.Phone = c.AlreadyExist
	}
	if owner.Email != "" && ownerE.Email == "" && !db.Where(&m.Owner{Email: owner.Email}).Take(&m.Owner{}).RecordNotFound() {
		ownerE.Email = c.AlreadyExist
	}

	if ownerE.ConfirmPassword == "" && ownerE.Password == "" && owner.Password != owner.ConfirmPassword {
		ownerE.ConfirmPassword = c.WrongConfirmPass
	}
	return ownerE
}

// Checks phone number and password for existence in database and correctness
func getUserLoginErrors(owner m.Owner) m.OwnerError {
	ownerE := m.OwnerError{}

	if owner.Phone == "" {
		ownerE.Phone =  c.Empty
	}
	if owner.Password == "" {
		ownerE.Password =  c.Empty
	}

	if ownerE.Phone == "" && !utils.IsPhoneNumber(owner.Phone) {
		ownerE.Phone = c.WrongPhone
	}

	return ownerE
}
func getWashRegistrationErrors(wash m.CarWash) error{
	var err error
	if len(wash.Name) == 0{
		err = errors.New("name is empty")
	}
	if len(wash.Address) == 0{
		err = errors.New("Address is empty")
	}
	return err

}

// Checks if user doesn't have errors
// It means user is valid for registration and login
func  isOwnerEmpty(ue m.OwnerError) bool {
	if ue.FirstName != "" {
		return false
	}
	if ue.SecondName != "" {
		return false
	}
	if ue.Phone != "" {
		return false
	}
	if ue.Email != "" {
		return false
	}
	if ue.Password != "" {
		return false
	}
	if ue.ConfirmPassword != "" {
		return false
	}
	return true
}
