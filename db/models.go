package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"my_app/helper"
)

func (DataBase *DataBase) InsertNewClinique(ownerName, user_name, email, number, password, city string) error {
	_, err := DataBase.DB.Exec(`INSERT INTO cliniques (owner_name, user_name, email, number, password_hash, is_admin, city)
		VALUES (?,?,?,?,?,0,?)`, ownerName, user_name, email, number, password, city)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (DataBase *DataBase) InsertNewPatient(first_name, last_name, whatsappNumber1, whatsappNumber2, email, cardId, city, surgeryDate string, cliniqueId, age int) error {
	_, err := DataBase.DB.Exec(`INSERT INTO patients (clinique_id, first_name, last_name, whatsapp_number1, whatsapp_number2, email, age, card_id, city, surgery_date)
		VALUES (?,?,?,?,?,?,?,?,?,?)`, cliniqueId, strings.TrimSpace(first_name), strings.TrimSpace(last_name), whatsappNumber1, whatsappNumber2, strings.TrimSpace(email), age, cardId, strings.TrimSpace(city), surgeryDate)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (Database *DataBase) CheckLoginInfo(login, password string) (helper.ApiResponse, helper.Profile) {
	var Profile helper.Profile
	if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
		return helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "Please fill all field",
		}, Profile
	}

	query := `
    SELECT id, user_name, city, is_admin
    FROM cliniques
    WHERE (user_name = ? COLLATE NOCASE OR email = ? COLLATE NOCASE)
      AND password = ?;
	`

	var is_admin bool

	row := Database.DB.QueryRow(query, login, login, password)
	err := row.Scan(&Profile.Id, &Profile.UserName, &Profile.City, &is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "please enter the correct data",
			}, Profile
		} else {
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

func (DataBase *DataBase) CheckIfThisCliniqueNameExist(userName string) bool {
	var exists bool
	err := DataBase.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM cliniques WHERE user_name = ? COLLATE NOCASE);`, userName).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return true
	}

	fmt.Println("exists: ", exists)

	return exists
}

func (DataBase *DataBase) IsEmailTakenByDifferentOwner(email, ownerName string) bool {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM cliniques
        WHERE email = ? COLLATE NOCASE AND owner_name != ? COLLATE NOCASE 
    `

	err := DataBase.DB.QueryRow(query, email, ownerName).Scan(&count)
	if err != nil {
		fmt.Printf("Database error in IsEmailTakenByDifferentOwner: %v\n", err)
		return true
	}

	return count > 0
}

func (DataBase *DataBase) IsNumberTakenByDifferentOwner(number, ownerName string) bool {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM cliniques
        WHERE number = ? COLLATE NOCASE AND owner_name != ? COLLATE NOCASE 
    `

	err := DataBase.DB.QueryRow(query, number, ownerName).Scan(&count)
	if err != nil {
		fmt.Printf("Database error in IsNumberTakenByDifferentOwner: %v\n", err)
		return true
	}

	return count > 0
}
