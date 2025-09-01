package db

var Tables = []string{
	`CREATE TABLE IF NOT EXISTS cliniques (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        owner_name VARCHAR(100) NOT NULL,
        user_name VARCHAR(100) UNIQUE NOT NULL,
        email VARCHAR(100) NOT NULL,
        number VARCHAR(20) NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        is_admin BOOLEAN,
        city VARCHAR(100),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`,

	`CREATE TABLE IF NOT EXISTS patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinique_id INTEGER NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    whatsapp_number1 VARCHAR(20) NOT NULL,
    whatsapp_number2 VARCHAR(20),
    email VARCHAR(100) NOT NULL,
    age INTEGER NOT NULL,
    card_id VARCHAR(50),
    city VARCHAR(100),
    jj_stent_insertion DATE NOT NULL,
    jj_stent_removal DATE NOT NULL,
    diagnostic TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (clinique_id) REFERENCES cliniques(id) ON DELETE CASCADE
    );`,

	`CREATE TABLE IF NOT EXISTS patients_archive (
    id INTEGER,
    clinique_id INTEGER,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    whatsapp_number1 VARCHAR(20),
    whatsapp_number2 VARCHAR(20),
    email VARCHAR(100),
    age INTEGER,
    card_id VARCHAR(50),
    city VARCHAR(100),
    jj_stent_insertion DATE,
    jj_stent_removal DATE,
    diagnostic TEXT,
    created_at DATETIME,
    FOREIGN KEY (clinique_id) REFERENCES cliniques(id) ON DELETE CASCADE
    );`,

	`CREATE TRIGGER IF NOT EXISTS archive_patient_before_delete
    BEFORE DELETE ON patients
    FOR EACH ROW
    BEGIN
        INSERT INTO patients_archive (
            id,
            clinique_id,
            first_name,
            last_name,
            whatsapp_number1,
            whatsapp_number2,
            email,
            age,
            card_id,
            city,
            jj_stent_insertion,
            jj_stent_removal,
            diagnostic,
            created_at
        )
        VALUES (
            OLD.id,
            OLD.clinique_id,
            OLD.first_name,
            OLD.last_name,
            OLD.whatsapp_number1,
            OLD.whatsapp_number2,
            OLD.email,
            OLD.age,
            OLD.card_id,
            OLD.city,
            OLD.jj_stent_insertion,
            OLD.jj_stent_removal,
            OLD.diagnostic,
            OLD.created_at
        );
    END;`,
}
