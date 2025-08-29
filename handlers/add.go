// handlers.go
package handlers

import (
	"fmt"
	"net/http"

	"my_app/helper"
)

// AddClinique handles clinic registration
func (h *Handlers) AddClinique(w http.ResponseWriter, r *http.Request) {
	var apiResponse helper.ApiResponse

	cliniqueData, err := helper.CheckCliniqueRegisterData(r.Body)
	if err != nil {
		fmt.Println("Validation error:", err)
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: err.Error(),
		}
		apiResponse.Sent(w)
		return
	}

	// Check if username is already taken
	if h.DB.CheckIfThisCliniqueNameExist(cliniqueData.UserName) {
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "this clinique name already used",
		}
		apiResponse.Sent(w)
		return
	}

	// Check if email is taken by different owner
	if h.DB.IsEmailTakenByDifferentOwner(cliniqueData.Email, cliniqueData.OwnerName) {
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "this email already used by another owner of a clinique",
		}
		apiResponse.Sent(w)
		return
	}

	// Check if number is taken by different owner
	if h.DB.IsNumberTakenByDifferentOwner(cliniqueData.Number, cliniqueData.OwnerName) {
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "this number already used by another owner of a clinique",
		}
		apiResponse.Sent(w)
		return
	}

	// Insert new clinique
	err = h.DB.InsertNewClinique(cliniqueData.OwnerName, cliniqueData.UserName, cliniqueData.Email, cliniqueData.Number, cliniqueData.Password, cliniqueData.City)
	if err != nil {
		fmt.Println("Database error:", err)
		apiResponse = helper.ApiResponse{
			Code: http.StatusInternalServerError,
			Data: "error inserting clinique in database, try again later",
		}
		apiResponse.Sent(w)
		return
	}

	// Success response
	apiResponse = helper.ApiResponse{
		Code: http.StatusCreated,
		Data: "clinique registered successfully",
	}
	apiResponse.Sent(w)
}

// AddPatient handles patient surgery registration
func (h *Handlers) AddPatient(w http.ResponseWriter, r *http.Request) {
	var apiResponse helper.ApiResponse

	patientData, err := helper.CheckPatientRegisterData(r.Body)
	if err != nil {
		fmt.Println("Validation error:", err)
		apiResponse = helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: err.Error(),
		}
		apiResponse.Sent(w)
		return
	}

	// Insert new patient
	err = h.DB.InsertNewPatient(
		patientData.FirstName,
		patientData.LastName,
		patientData.WhatsappNumber1,
		patientData.WhatsappNumber2,
		patientData.Email,
		patientData.CardId,
		patientData.City,
		patientData.SurgeryDate,
		2,
		patientData.Age,
	)
	if err != nil {
		fmt.Println("Database error:", err)
		apiResponse = helper.ApiResponse{
			Code: http.StatusInternalServerError,
			Data: "error registering patient surgery, try again later",
		}
		apiResponse.Sent(w)
		return
	}

	// Success response
	apiResponse = helper.ApiResponse{
		Code: http.StatusCreated,
		Data: "patient surgery registered successfully",
	}
	apiResponse.Sent(w)
}