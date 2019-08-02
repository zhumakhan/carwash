package constants

const (
	Root = "."
	ApiURI = "/wash/v1"

	RegisterURI = ApiURI + "/register"
	LoginURI    = ApiURI + "/login"

	ReadQRCode      = ApiURI+"/QRCode/:"+QRCode
	PayForWash      = ApiURI+"/pay"
	CheckAccessTokenURI = ApiURI + "/check-access"
	RefreshURI	= ApiURI + "/refresh"

	Managers    = ApiURI + "/managers" 
	Manager     = Managers + "/:" + UUIDPathVar

	WashesURI 	= ApiURI + "/washes"
	WashURI 	= WashesURI + "/:" + WashIdPathVar

	Workers     = ApiURI + "/workers"
	Worker      = Workers + "/:" + WorkerID
	
	WashServicesURI = WashURI + "/services"
	WashServiceURI  = WashServicesURI + "/:" + ServiceIdPathVar

	WashCarTypesURI = WashURI + "/car-types"
	WashCarTypeURI = WashCarTypesURI + "/:" + CarTypeIdPathVar

	ClientsURI = ApiURI + "/clients"

	BookingsURI = ApiURI + "/bookings"
	BookingURI  = BookingsURI + "/:" + UUIDPathVar

	StatisticsURI = ApiURI + "/statistics/:" + WashIdPathVar


 	AdminNotifyAllURI = ApiURI + "/notify-all-wash-admins" + "/{" + NotificationMsg + "}"

	AdminPanelURI          = ApiURI + "/admin-panel"
	AdminWashesURI         = AdminPanelURI + "/washes"
	AdminWashURI           = AdminWashesURI + "/:" + WashIdPathVar
	AdminWashOwnersURI     = AdminPanelURI + "/owners"
	AdminWashOwnerURI      = AdminWashOwnersURI + "/:" + UUIDPathVar
	AdminPaymentHistoriesURI = AdminPanelURI + "/payments/:" + WashIdPathVar
	AdminPaymentHistoryURI = AdminPanelURI + "/payments/:" + UUIDPathVar

	PostBookingURI = ApiURI + "/post-link/bookings/{uuid}"

	Service = ApiURI + "/service"
	ServiceStatus = Service + "/status"
)

//path variables
const (
	WashIdPathVar   = "wash-id"
	UUIDPathVar     = "uuid"
	ServiceIdPathVar = "service-id"
	CabinIdPathVar = "cabin-id"
	CarTypeIdPathVar = "car-type-id"
	NotificationMsg = "msg"
	QRCode			= "qr-code"
	WorkerID        = "worker-id"
)
