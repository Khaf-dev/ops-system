package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Seed() {
	log.Println("Seeding initial data...")

	// untuk hash password
	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	DB.Exec(`
	INSERT INTO users (name, email, password_hash, role) 
	VALUES ('Kaffa', 'rifyatkaffa@gmail.com', ?, 'Admin')
	ON CONFLICT (email) DO NOTHING
	`, string(password))

	DB.Exec(`
	INSERT INTO request_types (name) 
	VALUES ('Operasional Kantor'), ('Perizinan'), ('Maintenance'), ('Perjalanan Dinas')
	ON CONFLICT (name) DO NOTHING
	`)

	DB.Exec(`
	INSERT INTO activities (name) 
	VALUES ('Pembelian Barang'), ('Perbaikan Unit'), ('Transportasi')
	ON CONFLICT (name) DO NOTHING
	`)

	log.Println("Seeding Complete!")
}
