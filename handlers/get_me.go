package handlers

import (
	"fmt"
	"net/http"

	"my_app/auth"
	"my_app/helper"
)

func (h *Handlers) GetMe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i'm here here here")
	var apiResponse helper.ApiResponse
	userId := r.Context().Value(auth.UserContextKey).(auth.UserInfo).UserID
	profile, err := h.DB.GetCliniqueInfo(userId)
	if err != nil {
		fmt.Println(err)
		apiResponse = helper.ApiResponse{
			Code: http.StatusInternalServerError,
			Data: "Something wrong, Please reopen the app",
		}
	} else {
		apiResponse = helper.ApiResponse{
			Code: http.StatusOK,
			Data: profile,
		}
	}
	apiResponse.Sent(w)
}
