package main

import (
	"fmt"
	"log"

	"os"

	c "carwashes/constants"
	contr "carwashes/controllers"
	"carwashes/database"

	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const PUT = "PUT"

func main() {

	log.Println("turaQwash api server. Version 0.0.0")
	log.Println("Loading environment variables.")
	e := godotenv.Load(".env") //Load .env file
	if e != nil {
		log.Fatal(e)
	}

	port := os.Getenv("app_port")
	if port == "" {
		log.Fatal("$app_port not set")
	} else {
		log.Println(fmt.Sprintf("$app_port: %s", port))
	}

	database.Connect() //Load database

	router := router.New()
	router.Handle(GET, "/", contr.Hello)
	router.Handle(POST, c.RegisterURI, contr.Register)
	router.Handle(POST, c.LoginURI, contr.LoginOwner)

	router.Handle(GET, c.Managers, contr.GetManagers)
	router.Handle(GET, c.Manager, contr.GetManager)
	router.Handle(PUT, c.Manager, contr.PutManager)
	router.Handle(DELETE, c.Manager, contr.DeleteManager)
	
	router.Handle(POST, c.Workers, contr.AddWorker)
	router.Handle(GET, c.Workers, contr.GetWorkers)
	router.Handle(GET, c.Worker, contr.GetWorker)
	router.Handle(PUT, c.Worker, contr.ChangeWorker)
	router.Handle(DELETE, c.Worker, contr.DeleteWorker)

	//router.Handle(GET, c.ReadQRCode, contr.ReadQRCode)
	router.Handle(POST, c.PayForWash, contr.PayForWash)
	router.Handle(POST, c.WashServicesURI, contr.InsertWashService)
	router.Handle(PUT, c.WashServiceURI, contr.ChangeWashService)
	router.Handle(GET, c.WashServicesURI, contr.GetWashServices)
	router.Handle(GET, c.WashServiceURI, contr.GetWashService)
	router.Handle(DELETE, c.WashServiceURI, contr.DeleteWashService)

	router.Handle(POST, c.WashCarTypesURI, contr.AddCarType)
	router.Handle(GET, c.WashCarTypesURI, contr.GetCarTypes)
	router.Handle(DELETE, c.WashCarTypeURI, contr.DeleteCarType)
	router.Handle(PUT, c.WashCarTypeURI, contr.ChangeCarType)

	router.Handle(POST, c.WashesURI, contr.InsertWash)
	router.Handle(GET, c.WashesURI, contr.GetWashes)
	router.Handle(GET, c.WashURI, contr.GetWash)
	router.Handle(PUT, c.WashURI, contr.UpdateWash)
	router.Handle(DELETE, c.WashURI, contr.DeleteWash)
	router.Handle(GET, c.ClientsURI, contr.GetClients)

	router.Handle(PUT, c.BookingURI, contr.ChangeBooking)
	router.Handle(POST, c.BookingsURI, contr.InsertBooking)
	router.Handle(GET, c.BookingsURI, contr.GetBookings)
	router.Handle(GET, c.BookingURI, contr.GetBooking)
	router.Handle(DELETE, c.BookingURI, contr.DeleteBooking)

	router.Handle(GET, c.ServiceStatus, contr.GetStatus)

	router.Handle(GET, c.AdminWashOwnersURI, contr.AdminGetOwners)
	router.Handle(GET, c.AdminWashOwnerURI, contr.AdminGetOwner)
	router.Handle(POST, c.AdminWashesURI, contr.AdminInsertWash)
	router.Handle(GET, c.AdminWashesURI, contr.AdminGetWashes)
	router.Handle(GET, c.AdminWashURI, contr.AdminGetWash)
	router.Handle(GET, c.AdminPaymentHistoriesURI, contr.AdminGetPaymentHistory)
	router.Handle(PUT, c.AdminPaymentHistoryURI, contr.AdminProceedMonthlyPayment)

	router.Handle(GET, c.StatisticsURI, contr.GetStats)

	log.Println("Maintaining app")
	log.Fatal(fasthttp.ListenAndServe(":"+port, router.Handler))
}
