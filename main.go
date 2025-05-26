package main

import (
	"crud-buku-go/config"
	"crud-buku-go/models"
	"crud-buku-go/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env. Pastikan variabel environment sudah diatur.")
	}

	config.ConnectDB()
	defer func() {
		if config.DB != nil {
			config.DB.Close()
			log.Println("Koneksi database ditutup.")
		}
	}()

	models.SeedData()

	router := routes.SetupRoutes()

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Port default
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%s", appPort)

	log.Printf("🚀 Server berjalan di http://localhost:%s", appPort)
	log.Printf("🌐 Server juga dapat diakses di LAN pada alamat IP mesin Anda dengan port %s", appPort)
	log.Printf("Tekan CTRL+C untuk menghentikan server.")

	err = http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}