package helper

/************ Profiles ***************/

type Profile struct {
	Id       int
	UserName string
	Role     string
	City     string
	Data     interface{}
}

type Clinique struct {
	Id        int
	OwnerName string
	UserName  string
	Email     string
	Number    string
	City      int
}

type Patient struct {
	Id             int
	FirstName      string
	LastName       string
	WhatsappNumber string
	Age            int
	CardId         string
	City           string
	SurgeryDate    string
}

/************ Response Api ***************/

type ApiResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type LoginResponseData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type CliniqueRegisterApi struct {
	OwnerName string `json:"owner_name" validate:"required,min=3,max=50,ownername"`
	UserName  string `json:"user_name" validate:"required,min=3,max=20,ownername"`
	Email     string `json:"email" validate:"required,email"`
	Number    string `json:"number" validate:"required,numeric,min=3,max=10"`
	Password  string `json:"password" validate:"required,bcrypt_hash"`
	City      string `json:"city" validate:"required"`
}

type PatientRegisterApi struct {
	FirstName       string `json:"first_name" validate:"required,alpha,min=3,max=50"`
	LastName        string `json:"last_name" validate:"required,alpha,min=3,max=50"`
	WhatsappNumber1 string `json:"whatsapp_number1" validate:"required,numeric,min=3,max=10"`
	WhatsappNumber2 string `json:"whatsapp_number2" validate:"required,numeric,min=3,max=10"`
	Email           string `json:"email" validate:"required,email"`
	Age             int    `json:"age" validate:"required,gte=0,lte=120"`
	CardId          string `json:"card_id" validate:"required,alphanum,min=5,max=20"`
	City            string `json:"city" validate:"required,alpha,min=3,max=50"`
	SurgeryDate     string `json:"surgery_date"`
}
