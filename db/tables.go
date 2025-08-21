package db

var Tables = []string{
	`CREATE TABLE IF NOT EXISTS cliniques (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(100) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
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
        whatsapp_number VARCHAR(20) NOT NULL,
        age INTEGER,
        card_id VARCHAR(50) UNIQUE,
        city VARCHAR(100),
        surgery_date TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (clinique_id) REFERENCES cliniques(id) ON DELETE CASCADE
    );`,
}
