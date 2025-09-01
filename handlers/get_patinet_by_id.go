package handlers

import (
	"net/http"
	"strconv"

	"my_app/auth"
	"my_app/helper"

	"github.com/gorilla/mux"
)

func (h *Handlers) GetPatientById(w http.ResponseWriter, r *http.Request) {
	var ResponseApi helper.ApiResponse

	// date should be like this: YYYY-MM-DD
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ResponseApi = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "invalid patient id",
		}
		ResponseApi.Sent(w)
		return
	}

	user := r.Context().Value(auth.UserContextKey).(auth.UserInfo)

	patients, err := h.DB.GetPatientInfo(id, user.UserID)
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
