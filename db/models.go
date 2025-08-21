package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"my_app/helper"
)

func (Database *DataBase) CheckLoginInfo(login, password string) (helper.ApiResponse, helper.Profile) {
	var Profile helper.Profile
	if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
		return helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "Please fill all field",
		}, Profile
	}

	query := `
    SELECT id, name, city, is_admin
    FROM users
    WHERE (name = ? COLLATE NOCASE OR email = ? COLLATE NOCASE)
      AND password = ?;
	`

	var is_admin bool

	row := Database.DB.QueryRow(query, login, login, password)
	err := row.Scan(&Profile.Id, &Profile.Name, &Profile.City, &is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "please enter the correct data",
			}, Profile
		} else {
			fmt.Println(err)
			return helper.ApiResponse{
				Code: http.StatusInternalServerError,
				Data: "error scaning data from database",
			}, Profile
		}
	}
	if is_admin {
		Profile.Role = "Admin"
	} else {
		Profile.Role = "Clinique"
	}

	return helper.ApiResponse{
		Code: http.StatusOK,
		Data: "connecting succeffully",
	}, Profile
}
