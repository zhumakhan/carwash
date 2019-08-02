
package models

import (
	"time"
	"strconv"
)
var WashId uint = 1
const (
	PastBookingTablePrefix = "past_bookings_"
	PastServicesTablePrefix = "past_booking_services_"
)

//easyjson:json
type TopicPush struct {
	To           string `json:"to"`
	Notification struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"notification"`
	TimeToLive int `json:"time_to_live"`
}

//easyjson:json
type TopicSubscribe struct {
	To                 string   `json:"to"`
	RegistrationTokens []string `json:"registration_tokens"`
}

//easyjson:json
type OwnerError struct {
	FirstName       string `json:",omitempty"`
	SecondName      string `json:",omitempty"`
	Phone           string `json:",omitempty"`
	Email           string `json:",omitempty"`
	Password        string `json:",omitempty"`
	ConfirmPassword string `json:",omitempty"`
}

//easyjson:json
type Client struct {
	UUID string `gorm:"PRIMARY_KEY;TYPE:BINARY(36)"`
	FirstName string `gorm:"NOT NULL"`
	SecondName string
	MiddleName string
	Phone string `gorm:"NOT NULL;SIZE:12;UNIQUE"`
}

//easyjson:json
type Owner struct {
	UUID string `gorm:"PRIMARY_KEY;TYPE:BINARY(36)"`
	WashID uint `gorm:"NULL"`//in case if owner is manager <=> not director, WashID is set.
	FirstName string `gorm:"NOT NULL"`
	SecondName string `gorm:"NOT NULL"`
	MiddleName string `gorm:"NOT NULL"`
	Phone string `gorm:"NOT NULL;SIZE:12;UNIQUE"`
	Password string `gorm:"NOT NULL" json:",omitempty"`
	ConfirmPassword string `gorm:"-" json:",omitempty"`
	Email string `gorm:"NOT NULL"`
	AccessToken string `gorm:"-"`
	RefreshToken string `gorm:"-"`
	CreatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:",omitempty"`
	Role    uint  	`gorm:"NULL; DEFAULT:1" schema:"-" json:",omitempty"`//1 highedt rank, 2 second highest rank
}

type Admin struct {
	UUID string `gorm:"PRIMARY_KEY;TYPE:BINARY(36)"`
	FirstName string `gorm:"NOT NULL"`
	SecondName string `gorm:"NOT NULL"`
	MiddleName string `gorm:"NOT NULL"`
	Phone string `gorm:"NOT NULL;SIZE:12;UNIQUE"`
	Password string `gorm:"NOT NULL" json:",omitempty"`
	ConfirmPassword string `gorm:"-" json:",omitempty"`
	Email string `gorm:"NOT NULL"`
	AccessToken string `gorm:"-"`
	RefreshToken string `gorm:"-"`
	CreatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:",omitempty"`
}
type User struct {
	UUID string `gorm:"PRIMARY_KEY;TYPE:BINARY(36)"`
	FirstName string `gorm:"NOT NULL"`
	SecondName string `gorm:"NOT NULL"`
	MiddleName string `gorm:"NOT NULL"`
	Phone string `gorm:"NOT NULL;SIZE:12;UNIQUE"`
	Password string `gorm:"NOT NULL" json:",omitempty"`
	ConfirmPassword string `gorm:"-" json:",omitempty"`
	Email string `gorm:"NOT NULL"`
	AccessToken string `gorm:"-"`
	RefreshToken string `gorm:"-"`	
	CreatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:",omitempty"`
}
type Worker struct{
	ID   uint `gorm:"PRIMARY_KEY; AUTO_INCREMENT" json:",omitempty"`
	Name string `gorm:"NOT NULL; SIZE:50"`
	Phone string `gorm:"NOT NULL;SIZE:20;"`
	WashID uint  `gorm:"NOT NULL"`
}
type Workers struct{
	Workers []Worker `json:""`
}
//easyjson:json
type CarWash struct {
	ID uint `gorm:"PRIMARY_KEY; AUTO_INCREMENT" json:",omitempty"`
	Name string `gorm:"NOT NULL" json:",omitempty"`
	Address string `gorm:"NOT NULL" json:",omitempty"`
	Longitude float64 `gorm:"NOT NULL" json:",omitempty"`
	Latitude float64 `gorm:"NOT NULL" json:",omitempty"`
	Owner string `gorm:"NOT NULL;TYPE:BINARY(36)" json:",omitempty"`
	Photo string `gorm:"SIZE:4096" json:",omitempty"`
	PaidUntil time.Time  `gorm:"TYPE:DATETIME" schema:"-" json:"-"`  
	CreatedAt time.Time  `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:"-"`
	Status   uint `gorm:"NOT NULL; DEFAULT:0" json:",omitempty"`//TODO meaning of status
	Services []Service `gorm:"foreignkey:CarWash;association_foreignkey:ID" json:",omitempty"`
	CarTypes []CarType `gorm:"foreignkey:CarWashID;association_foreignkey:ID" json:",omitempty"`
	QueueSize uint `gorm:"-"`
}

//easyjson:json
type CarType struct {
	ID uint `gorm:"PRIMARY_KEY; AUTO_INCREMENT"`
	Name string `gorm:"NOT NULL"`
	CarWashID uint `gorm:"NOT NULL"`
}

//easyjson:json
type Service struct {
	ID uint `gorm:"PRIMARY_KEY; AUTO_INCREMENT"`
	CarWash uint `gorm:"NOT NULL"`
	Name string `gorm:"NOT NULL"`
	Description string
	Costs []ServiceCost `gorm:"foreignkey:Service"`
}

//easyjson:json
type ServiceCost struct {
	Service uint `gorm:"NOT NULL;PRIMARY_KEY"`
	CarModel string `gorm:"NOT NULL;PRIMARY_KEY"`
	Cost uint `gorm:"NOT NULL"`
	Duration uint `gorm:"NOT NULL"`
}
//easyjson:json
type BookingService struct {
	ServiceID uint `gorm:"NOT NULL;PRIMARY_KEY"`
	Service Service `gorm:"foreignkey:ID;association_foreignkey:ServiceID"`
	Booking string `gorm:"NOT NULL;TYPE:BINARY(36);PRIMARY_KEY"`
	Cost uint `gorm:"NOT NULL"`
	Duration uint `gorm:"NOT NULL"` // duration of service to complete in minutes
}

//easyjson:json
type Booking struct {
	UUID string `gorm:"PRIMARY_KEY; NOT NULL;TYPE:BINARY(36)"`
	CarWash uint `gorm:"NOT NULL"`
	ClientUUID string `gorm:"TYPE:BINARY(36);NOT NULL"`
	Client Client `gorm:"foreignkey:UUID;association_foreignkey:ClientUUID"`
	CarNumber string `gorm:"SIZE:12;NOT NULL"`
	Cost uint `gorm:"NOT NULL"`
	CreatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL"`
	UpdatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL"`
	PaymentStatus uint `gorm:"NOT NULL; DEFAULT:0;"`//0 = not payed, 1 = payed
	Status uint `gorm:"NOT NULL; DEFAULT:0;"`//0 = in queue, 1 = waiting for owner
	Order uint `gorm:"NOT NULL"`
	CarModel string `gorm:"NOT NULL"`
	Vehicle string
	BookingServices []BookingService `gorm:"foreignkey:Booking;association_foreignkey:UUID"`
	ServiceCosts []ServiceCost `gorm:"-" json:"-"`
	WokerID uint `gorm:"NOT NULL"`
	RemoteBooked uint //0 for false, 1 for true
}
type PastBooking struct {
	ID uint `gorm:"PRIMARY_KEY; AUTO_INCREMENT"`
	CarWash uint `gorm:"NOT NULL"`
	ClientUUID string `gorm:"type:binary(36);NOT NULL"`
	ClientFirstName string `gorm:"NOT NULL"`
	ClientSecondName string `gorm:"NOT NULL"`
	ClientMiddleName string `gorm:"NOT NULL"`
	ClientPhone string `gorm:"NOT NULL"`
	CarNumber string `gorm:"SIZE:12;NOT NULL"`
	Cost uint `gorm:"NOT NULL"`
	CreatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL" schema:"-"`
	UpdatedAt time.Time `gorm:"TYPE:DATETIME;NOT NULL"`
	CarModel string `gorm:"NOT NULL"`
	BookingServices string
}
func (p PastBooking)TableName() string{
	return PastBookingTablePrefix + strconv.FormatUint(uint64(WashId), 10)
}
//easyjson:json
type CarWashes struct{
	Washes []CarWash
}

//easyjson:json
type WashOwners struct {
	Owners []Owner
}

//easyjson:json
type OwnerWithWashes struct {
	Owner  Owner
	Washes []CarWash
}

//easyjson:json
type MonthlyPaymentHistory struct {
	UUID string `gorm:"NOT NULL;PRIMARY_KEY;TYPE:BINARY(36)"`
	CarWash uint `gorm:"NOT NULL"`
	Month time.Time `gorm:"TYPE:DATETIME;NOT NULL"`
	Amount uint `gorm:"NOT NULL"`
	Status uint `gorm:"NOT NULL;DEFAULT:0;"`
}

//easyjson:json
type CarWashesWithOwnerName struct {
	OwnerUUID []string
	Washes []CarWash //Здесь вместо UUID Owner-а будет его имя
}

//easyjson:json
type Histories struct {
	History []MonthlyPaymentHistory
}