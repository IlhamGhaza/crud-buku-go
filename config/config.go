package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env, menggunakan variabel environment sistem jika ada.")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Port database tidak valid: %v", err)
	}

	initialConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	tempDB, err := sql.Open("postgres", initialConnStr)
	if err != nil {
		log.Fatalf("Gagal terhubung ke server PostgreSQL: %v", err)
	}
	defer tempDB.Close()

	rows, err := tempDB.Query(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName))
	if err != nil {
		log.Fatalf("Gagal memeriksa keberadaan database: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {

		log.Printf("Database '%s' tidak ditemukan, mencoba membuat...", dbName)
		_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("Gagal membuat database '%s': %v", dbName, err)
		}
		log.Printf("Database '%s' berhasil dibuat.", dbName)
	} else {
		log.Printf("Database '%s' sudah ada.", dbName)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}

	DB = database
	log.Println("Berhasil terhubung ke database PostgreSQL!")

	createTable()
}

func createTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		year INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Gagal membuat tabel 'books': %v", err)
	}
	log.Println("Tabel 'books' siap digunakan.")
}
