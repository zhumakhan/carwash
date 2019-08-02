package database

import (
	"carwashes/models"
)

func drop() {
	db.DropTableIfExists(&models.BookingService{})
	db.DropTableIfExists(&models.ServiceCost{})
	db.DropTableIfExists(&models.Service{})
	db.DropTableIfExists(&models.Booking{})
	db.DropTableIfExists(&models.CarWash{})
	db.DropTableIfExists(&models.Owner{})
	db.DropTableIfExists(&models.Client{})
	db.DropTableIfExists(&models.Admin{})
	db.DropTableIfExists(&models.Worker{})
}

func migrate() {
	db.AutoMigrate(&models.Client{})
	db.AutoMigrate(&models.Owner{})
	db.AutoMigrate(&models.CarWash{})
	db.AutoMigrate(&models.Service{})
	db.AutoMigrate(&models.BookingService{})
	db.AutoMigrate(&models.Booking{})
	db.AutoMigrate(&models.ServiceCost{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.CarType{})
	db.AutoMigrate(&models.MonthlyPaymentHistory{})
	db.AutoMigrate(&models.Worker{})

	dbpast.AutoMigrate(&models.PastBooking{})
}

func foreignkey() {
	db.Model(&models.ServiceCost{}).AddForeignKey("service", "services(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CarWash{}).AddForeignKey("owner", "owners(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&models.Service{}).AddForeignKey("car_wash", "car_washes(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Booking{}).AddForeignKey("car_wash", "car_washes(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Booking{}).AddForeignKey("client_uuid", "clients(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&models.BookingService{}).AddForeignKey("booking", "bookings(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&models.BookingService{}).AddForeignKey("service_id", "services(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CarType{}).AddForeignKey("car_wash_id", "car_washes(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MonthlyPaymentHistory{}).AddForeignKey("car_wash", "car_washes(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Worker{}).AddForeignKey("wash_id", "car_washes(id)", "RESTRICT", "RESTRICT")
}

func index() {
	db.Model(&models.Booking{}).AddIndex("idx_bookings_car_wash", "car_wash")
	dbpast.Model(&models.PastBooking{}).AddIndex("idx_past_bookings_car_wash", "car_wash")
	dbpast.Model(&models.PastBooking{}).AddIndex("idx_past_bookings_updated_at", "updated_at")

	db.Exec("ALTER TABLE `washes`.`clients` ADD FULLTEXT (`first_name`);")
	db.Exec("ALTER TABLE `washes`.`clients` ADD FULLTEXT (`phone`);")
	db.Exec("ALTER TABLE `washes`.`clients` ADD FULLTEXT (`first_name`, `phone`);")
}
/*
func insertSome() {
	db.Create(&models.Client{ UUID : "123456789012345678901234567890123456" })
	db.Create(&models.Owner{ UUID : "123456789012345678901234567890123456", Phone : "1234", Password : "1234" })
	db.Create(&models.CarWash{ Owner : "123456789012345678901234567890123456" })
	db.Create(&models.Booking{ UUID: "123456789012345678901234567890123456", CarWash : 1, ClientUUID: "123456789012345678901234567890123456" })
	db.Create(&models.Service{ CarWash : 1 })
	db.Create(&models.ServiceCost{ Service : 1 })
	db.Create(&models.BookingService{ Booking : "123456789012345678901234567890123456", Service : 1 })
}
*/
