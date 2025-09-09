package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"my_app/helper"

	"golang.org/x/crypto/bcrypt"
)

/************************** INSERTING **************************/

func (DataBase *DataBase) InsertNewClinique(ownerName, user_name, email, number, password, city string) error {
	_, err := DataBase.DB.Exec(`INSERT INTO cliniques (owner_name, user_name, email, number, password_hash, is_admin, city)
		VALUES (?,?,?,?,?,0,?)`, ownerName, user_name, email, number, password, city)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (DataBase *DataBase) InsertNewPatient(first_name, last_name, whatsappNumber1, whatsappNumber2, email, cardId, city, insertion, removal, diagnostic string, cliniqueId, age int) error {
	/* check if this surgery date already exists */
	var exists bool
	err := DataBase.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM patients
	 WHERE clinique_id = ? 
	 	AND first_name = ? COLLATE NOCASE
		AND last_name = ? COLLATE NOCASE
		AND card_id = ? COLLATE NOCASE);`, cliniqueId, first_name, last_name, cardId).Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}
	if exists {
		return errors.New("already exists")
	}

	/* insert new surgery date */
	_, err = DataBase.DB.Exec(`INSERT INTO patients (clinique_id, first_name, last_name, whatsapp_number1, whatsapp_number2, email, age, card_id, city, jj_stent_insertion, jj_stent_removal, diagnostic)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, cliniqueId, strings.TrimSpace(first_name), strings.TrimSpace(last_name), whatsappNumber1, whatsappNumber2, strings.TrimSpace(email), age, cardId, strings.TrimSpace(city), insertion, removal, diagnostic)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

/************************** CHECKING **************************/

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
    WHERE user_name = ? COLLATE NOCASE OR email = ? COLLATE NOCASE;
	`

	var is_admin bool

	row := Database.DB.QueryRow(query, login, login, password)
	err := row.Scan(&Profile.ID, &Profile.UserName, &Profile.City, &is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ApiResponse{
				Code: http.StatusBadRequest,
				Data: "invalid login",
			}, Profile
		} else {
			fmt.Println(err)
			return helper.ApiResponse{
				Code: http.StatusInternalServerError,
				Data: "error scaning data from database",
			}, Profile
		}
	}

	err = Database.ComparePasswords(Profile.ID, password)
	if err != nil {
		return helper.ApiResponse{
			Code: http.StatusBadRequest,
			Data: "invalid password",
		}, Profile
	}

	if is_admin {
		Profile.Role = "admin"
	} else {
		Profile.Role = "clinique"
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

func (DataBase *DataBase) IsValidRole(id int, role string) bool {
	is_admin := strings.ToLower(role) == "admin"
	var exists bool
	err := DataBase.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM cliniques WHERE id = ? AND is_admin = ?);`, id, is_admin).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (DataBase *DataBase) ComparePasswords(id int, password string) error {
	var password_stored string
	err := DataBase.DB.QueryRow("SELECT password_hash FROM cliniques WHERE id = ?", id).Scan(&password_stored)
	if err != nil {
		fmt.Println(err)
		return errors.New("error scaning the password from database")
	}
	err = bcrypt.CompareHashAndPassword([]byte(password_stored), []byte(password))
	return err
}

/***************************** GET **********************************/

func (d *DataBase) GetPatientsFromRemovalDate(date string, cliniqueID int) ([]helper.PatientApi, error) {
	rows, err := d.DB.Query(`
		SELECT id, clinique_id, first_name, last_name, whatsapp_number1, whatsapp_number2,
		       email, age, card_id, city, jj_stent_insertion, jj_stent_removal,
		       diagnostic, created_at
		FROM patients WHERE jj_stent_removal = ? AND clinique_id = ?;
	`, date, cliniqueID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("no patient on this day")
		}
		return nil, errors.New("internal server error")
	}
	defer rows.Close()

	var patients []helper.PatientApi
	for rows.Next() {
		var p helper.PatientApi
		err := rows.Scan(
			&p.ID, &p.CliniqueID, &p.FirstName, &p.LastName,
			&p.WhatsappNumber1, &p.WhatsappNumber2, &p.Email, &p.Age,
			&p.CardID, &p.City, &p.JJStentInsertion, &p.JJStentRemoval,
			&p.Diagnostic, &p.CreatedAt,
		)
		if err != nil {
			fmt.Println(err)
			if err == sql.ErrNoRows {
				return nil, errors.New("invalid JJ Stent removal date")
			}
			return nil, errors.New("internal server error")
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (d *DataBase) GetPatientsFromScheduler(date string) ([]helper.Patient, error) {
	rows, err := d.DB.Query(`
		SELECT id, first_name, last_name, age, card_id, whatsapp_number1, whatsapp_number2
		FROM patients WHERE jj_stent_removal = ?;
	`, date)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("no patient on this day")
		}
		return nil, errors.New("internal server error")
	}
	defer rows.Close()

	var patients []helper.Patient
	for rows.Next() {
		var p helper.Patient
		err := rows.Scan(
			&p.ID, &p.FirstName, &p.LastName, &p.Age, &p.CardID, &p.WhatsappNumber1, &p.WhatsappNumber2,
		)
		if err != nil {
			fmt.Println(err)
			if err == sql.ErrNoRows {
				return nil, errors.New("invalid JJ Stent removal date")
			}
			return nil, errors.New("internal server error")
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (d *DataBase) GetPatientInfo(patientID, cliniqueID int) (*helper.PatientApi, error) {
	row := d.DB.QueryRow(`
		SELECT id, clinique_id, first_name, last_name, whatsapp_number1, whatsapp_number2,
		       email, age, card_id, city, jj_stent_insertion, jj_stent_removal,
		       diagnostic, created_at
		FROM patients WHERE id = ? AND clinique_id = ?;
	`, patientID, cliniqueID)

	var p helper.PatientApi
	err := row.Scan(
		&p.ID, &p.CliniqueID, &p.FirstName, &p.LastName,
		&p.WhatsappNumber1, &p.WhatsappNumber2, &p.Email, &p.Age,
		&p.CardID, &p.City, &p.JJStentInsertion, &p.JJStentRemoval,
		&p.Diagnostic, &p.CreatedAt,
	)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("patient not found")
		}
		return nil, errors.New("internal server error")
	}

	return &p, nil
}

func (d *DataBase) GetAllPatientByCliniqueID(cliniqueID int) ([]helper.Patient, error) {
	rows, err := d.DB.Query(`SELECT id, first_name, last_name, age, card_id, jj_stent_removal
		FROM patients WHERE clinique_id = ?`, cliniqueID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, errors.New("no patient right now")
		}
		return nil, errors.New("internal server error")
	}

	var patients []helper.Patient
	for rows.Next() {
		var p helper.Patient
		err := rows.Scan(
			&p.ID, &p.FirstName, &p.LastName, &p.Age,
			&p.CardID, &p.JJStentRemoval,
		)
		if err != nil {
			fmt.Println(err)
			if err == sql.ErrNoRows {
				return nil, errors.New("no patient right now")
			}
			return nil, errors.New("internal server error")
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (d *DataBase) GetPatientsCount(cliniqueID int) (int, error) {
	var count int
	err := d.DB.QueryRow("SELECT COUNT(*) FROM patients WHERE clinique_id = ?", cliniqueID).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return count, errors.New("internal server error")
	}
	return count, nil
}

/****************************** DELETE ********************************/

func (d *DataBase) DeletePatient(patientID, cliniqueId int) error {
	_, err := d.DB.Exec(`DELETE FROM patients WHERE id = ? AND clinique_id = ?;`, patientID, cliniqueId)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (d *DataBase) DeleteClinique(cliniqueID int) error {
	_, err := d.DB.Exec(`DELETE FROM cliniques WHERE id = ?;`, cliniqueID)
	return err
}
