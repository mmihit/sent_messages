package handlers

import (
	"encoding/json"
	"net/http"

	"my_app/auth"
	"my_app/helper"
)

type LoginRequestApi struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var bodyRequest LoginRequestApi
	var apiResponse helper.ApiResponse
	var loginData helper.LoginResponseData

	err := json.NewDecoder(r.Body).Decode(&bodyRequest)
	if err != nil {
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "invalid body request",
		}
		apiResponse.Sent(w)
		return
	}

	// password_hashed, err := helper.HashPassword(bodyRequest.Password)
	// if err != nil {
	// 	apiResponse = helper.ApiResponse{
	// 		Code: http.StatusBadRequest,
	// 		Data: "error hashing password",
	// 	}
	// 	apiResponse.Sent(w)
	// 	return
	// }

	apiResponse, Profile := h.DB.CheckLoginInfo(bodyRequest.Login, bodyRequest.Password)
	if apiResponse.Code == http.StatusOK {
		token, err := auth.CreateToken(Profile)
		if err != nil {
			apiResponse = helper.ApiResponse{
				Code: http.StatusInternalServerError,
				Data: "error generating token, please try again",
			}
		} else {
			loginData = helper.LoginResponseData{
				Message: "connecting succeffully",
				Token:   token,
			}
			apiResponse.Data = loginData
		}
	}

	apiResponse.Sent(w)
}
