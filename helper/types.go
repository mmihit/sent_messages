package helper

/************ Profiles ***************/

type Profile struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	City     string `json:"city"`
}

type Clinique struct {
	ID        int
	OwnerName string
	UserName  string
	Email     string
	Number    string
	City      int
}

type Patient struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	CardID         string `json:"card_id"`
	Age            int    `json:"age"`
	JJStentRemoval string `json:"JJ_stent_removal"`
}

type PatientApi struct {
	ID               int    `json:"id"`
	CliniqueID       int    `json:"clinique_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	WhatsappNumber1  string `json:"whatsapp_number1"`
	WhatsappNumber2  string `json:"whatsapp_number2"`
	Email            string `json:"email"`
	Age              int    `json:"age"`
	CardID           string `json:"card_id"`
	City             string `json:"city"`
	JJStentInsertion string `json:"insertion_stent_JJ"`
	JJStentRemoval   string `json:"removal_stent_JJ"`
	Diagnostic       string `json:"diagnostic"`
	CreatedAt        string `json:"created_at"`
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
	FirstName        string `json:"first_name" validate:"required,alpha,min=3,max=50"`
	LastName         string `json:"last_name" validate:"required,alpha,min=3,max=50"`
	WhatsappNumber1  string `json:"whatsapp_number1" validate:"required,numeric,min=3,max=10"`
	WhatsappNumber2  string `json:"whatsapp_number2" validate:"numeric,min=3,max=10"`
	Email            string `json:"email" validate:"required,email"`
	Age              int    `json:"age" validate:"required,gte=0,lte=120"`
	CardId           string `json:"card_id" validate:"required,alphanum,min=5,max=20"`
	City             string `json:"city" validate:"required,alpha,min=3,max=50"`
	InsertionStentJJ string `json:"insertion_stent_JJ" validate:"required"`
	RemovalStentJJ   string `json:"removal_stent_JJ" validate:"required,future_date"`
	Diagnostic       string `json:"diagnostic" validate:"max=1500"`
}
