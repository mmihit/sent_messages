package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"my_app/auth"
	"my_app/helper"

	"github.com/gorilla/mux"
)

func (h Handlers) DeletePatient(w http.ResponseWriter, r *http.Request) {
	var ResponseApi helper.ApiResponse

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		ResponseApi = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "invalid patient id",
		}
		ResponseApi.Sent(w)
		return
	}

	user := r.Context().Value(auth.UserContextKey).(auth.UserInfo)

	err = h.DB.DeletePatient(id, user.UserID)
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
		Data: "patient removed succeffully",
	}
	ResponseApi.Sent(w)
}
