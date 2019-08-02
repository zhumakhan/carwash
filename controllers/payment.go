package controllers
import (
	c "carwashes/constants"
	"carwashes/database"
	m "carwashes/models"
	"github.com/valyala/fasthttp"
	"os"
	"errors"
	"strconv"
	jwt "github.com/dgrijalva/jwt-go"
	"fmt"
	"carwashes/utils"	
	"encoding/json"
	"time"
	"bytes"
	"net/http"
	"io/ioutil"
//	"strings"
)
const	payment_register = "http://127.0.0.1:8081/payment/v1/register-or-update"
const 	createOperation = "http://127.0.0.1:8081/payment/v1/operations"
const   post_link       = "http://127.0.0.1:8081/payment/v1/booking?api=wash"
type UserCredit struct {
		UserUUID     string    `json:"UserUUID"`
		Amount       float64    `json:"Amount"`
		Description  string    `json:"Description"`
		EarnDate     time.Time `json:"EarnDate"`
		PayoutDate   time.Time `json:"PayoutDate"`
		IncomeStatus int       `json:"IncomeStatus"`
	}
type UserDebt    struct {
		UserUUID    string `json:"UserUUID"`
		Amount      float64    `json:"Amount"`
		Description string `json:"Description"`
	}
type BankCredit struct {
		Amount      float64 `json:"Amount"`
		Description string  `json:"Description"`
	}
type CompanyCredit struct {
		Amount      float64    `json:"Amount"`
		Description string `json:"Description"`
	} 
type Operation struct {
	TypeID      int    `json:"TypeID"`
	Description string `json:"Description"`
	UserDebt    UserDebt `json:"UserDebt"`
	UserCredits []UserCredit `json:"UserCredits"`
	BankCredit BankCredit `json:"BankCredit"`
	CompanyCredit CompanyCredit `json:"CompanyCredit"`
}
func PayForWash(ctx *fasthttp.RequestCtx){
	var ownerUUID string
	var cost float64
	user := utils.Authorize(ctx)
	switch user.(type) {
	case m.User:
		ownerUUID = string(ctx.FormValue("owner"))
		temp, err := strconv.ParseFloat(string(ctx.FormValue("cost")), 64)
		if err != nil{
			respondWithError(ctx, fasthttp.StatusBadRequest, "Values are not set")
			return	
		}
		cost = temp//decalred not used problem
	default:
		respondWithError(ctx, fasthttp.StatusBadRequest, "You are not authorized")
		return
	}
	res, err := processWashPayment(ownerUUID, user.(m.User), cost)
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return	
	}
	ctx.Response.Header.SetStatusCode( fasthttp.StatusOK)
	wb, err := json.Marshal(map[string]interface{}{"message":"OK.","cost":cost, "post-link":post_link,"OrderBase64":(res["OrderBase64"].(string))})
	if err != nil {
		fmt.Print(err)
		ctx.SetContentType("text/plain")
		ctx.WriteString(err.Error())
	} else {
		ctx.SetContentType("application/json")
		ctx.Write(wb)
	}
}
func processWashPayment(ownerUUID string, user m.User, cost float64) (map[string]interface{}, error) {
	operation := Operation{
		TypeID:1,
		Description:"Pay from wash",
		UserDebt:UserDebt{
			UserUUID:user.UUID,
			Amount:cost,
			Description:"Pay from wash",
		},
		BankCredit:BankCredit{
			Amount: cost,
			Description:"Pay from wash",
		},
		CompanyCredit:CompanyCredit{
			Amount:cost,
			Description:"Pay from wash",
		},
	}
	operation.UserCredits = make([]UserCredit,1)
	operation.UserCredits[0] = UserCredit{
			UserUUID :ownerUUID,
			Amount	:cost,
			Description :"Pay from wash",
			EarnDate    : time.Now().UTC(),
			PayoutDate   :time.Now().AddDate(100,0,0).UTC(),
			IncomeStatus :2,
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader([]byte("user=" + user.UUID + "&phone=" + user.Phone))
	var resp *http.Response
	resp, err = http.Post(payment_register, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}
	owner := m.Owner { UUID: ownerUUID }
	err = database.GetDB().Table("washes.owners").Take(&owner).Error
	if err != nil {
		return nil, err
	}
	body = bytes.NewReader([]byte("user=" + ownerUUID + "&phone=" + owner.Phone))
	resp, err = http.Post(payment_register, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}

	body = bytes.NewReader(operationJson)
	resp, err = http.Post(createOperation, "application/json", body)
	if err != nil {
		return nil, err
	}

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	res := &map[string]interface{}{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func ReadQRCode(ctx *fasthttp.RequestCtx){
	var wash m.CarWash
	QRCode:= ctx.UserValue(c.QRCode).(string)
	if len(QRCode) == 0{
		respondWithError(ctx, fasthttp.StatusBadRequest, "QRCode is not passed")
		return
	}
	ID,err := getID(QRCode)
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	washId64,err := strconv.ParseUint(ID, 10, 64)
	washId := uint(washId64)
	
	if err != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, err.Error())
		return
	}
	if database.GetDB().Where(&m.CarWash{ID : washId}).Take(&wash).Error != nil{
		respondWithError(ctx, fasthttp.StatusBadRequest, "Can't find wash")
		return	
	}
	respondWithJSON(ctx, fasthttp.StatusOK, wash)
}

//for QR code
func generateToken(value string) string{
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["UUID"] = value
	claims["exp"] = "2100-01-01T00:00:01"

	token.Claims = claims
	tokenString, _:= token.SignedString([]byte(os.Getenv("access_token_password")))
	return tokenString
}
func getID(tokenStr string) (string, error){
	var ID = ""
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("access_token_password")), nil
	})

	if !token.Valid {
		return "", errors.New("QR code is not valid")

	}
	for key, val := range claims {
		if key == "UUID"{
			ID=fmt.Sprintf("%v",val)
		}
	}
	if len(ID) > 0 {
		return ID, nil
	}
	return ID, errors.New("Cannot read QR code")

}