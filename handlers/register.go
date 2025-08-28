package handlers

import (
	"fmt"
	"net/http"

	"my_app/helper"
)

type RegisterRequestApi struct {
	Role string `json:"role"`
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var apiResponse helper.ApiResponse
	data, role, err := helper.CheckRegisterData(r.Body)
	if err != nil {
		fmt.Println("here")
		fmt.Println(err)
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: err.Error(),
		}
		apiResponse.Sent(w)
		return
	}

	apiResponse = helper.ApiResponse{
		Code: http.StatusOK,
		Data: "you're registred seccuffully",
	}

	if role == "clinique" {

		cliniqueData := data.(helper.CliniqueRegisterApi)
		if h.DB.CheckIfThisCliniqueNameExist(cliniqueData.UserName) {
			apiResponse = helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "this clinique name already used",
			}
		} else if h.DB.IsEmailTakenByDifferentOwner(cliniqueData.Email, cliniqueData.OwnerName) {
			apiResponse = helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "this email already used by another owner of a clinique",
			}
		} else if h.DB.IsNumberTakenByDifferentOwner(cliniqueData.Number, cliniqueData.OwnerName) {
			apiResponse = helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "this number already used by another owner of a clinique",
			}
		} else {
			err := h.DB.InsertNewClinique(cliniqueData.OwnerName, cliniqueData.UserName, cliniqueData.Email, cliniqueData.Number, cliniqueData.Password, cliniqueData.City)
			if err != nil {
				apiResponse = helper.ApiResponse{
					Code: http.StatusInternalServerError,
					Data: "error inserting user in database, try again later",
				}
			}
		}

	} else {
		patientData := data.(helper.PatientRegisterApi)
		if !helper.ValidateSurgeryDate(patientData.SurgeryDate, "2000-01-01") {
			apiResponse = helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "invalid data, it should be at least tomorrow",
			}
		} else {
			err := h.DB.InsertNewPatient(patientData.FirstName, patientData.LastName, patientData.WhatsappNumber1, patientData.WhatsappNumber2, patientData.Email, patientData.CardId, patientData.City, patientData.SurgeryDate, 2, patientData.Age)
			if err != nil {
				apiResponse = helper.ApiResponse{
					Code: http.StatusInternalServerError,
					Data: "error inserting user in database, try again later",
				}
			}
		}
	}

	apiResponse.Sent(w)
}
