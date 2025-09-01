package handlers

import (
	"net/http"

	"my_app/auth"
	"my_app/helper"
)

func (h *Handlers) GetAllPatients(w http.ResponseWriter, r *http.Request) {
	var ResponseApi helper.ApiResponse
	user := r.Context().Value(auth.UserContextKey).(auth.UserInfo)

	patients, err := h.DB.GetAllPatientByCliniqueID(user.UserID)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() != "internal server error" {
			code = http.StatusBadRequest
		}
		ResponseApi = helper.ApiResponse{
			Code: code,
			Data: err.Error(),
		}
		ResponseApi.Sent(w)
		return
	}

	ResponseApi = helper.ApiResponse{
		Code: http.StatusOK,
		Data: patients,
	}
	ResponseApi.Sent(w)
}
