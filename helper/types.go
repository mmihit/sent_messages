package helper

type ApiResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type Profile struct {
	Id   int
	Name string
	Role string
	City string
}

type LoginResponseData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}