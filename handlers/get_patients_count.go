package handlers

import (
	"net/http"

	"my_app/auth"
	"my_app/helper"
)

func (h *Handlers) GetPatientsCount(w http.ResponseWriter, r *http.Request) {
	var ResponseApi helper.ApiResponse
	user := r.Context().Value(auth.UserContextKey).(auth.UserInfo)

	count, err := h.DB.GetPatientsCount(user.UserID)
	if err != nil {
		ResponseApi = helper.ApiResponse{
			Code: http.StatusInternalServerError,
			Data: "internal server error",
		}
		ResponseApi.Sent(w)
		return
	}
	ResponseApi = helper.ApiResponse{
		Code: http.StatusOK,
		Data: count,
	}
	ResponseApi.Sent(w)
}
