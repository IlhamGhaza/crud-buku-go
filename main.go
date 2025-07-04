package main

import (
	"crud-buku-go/config"
	"crud-buku-go/models"
	"crud-buku-go/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	_ "crud-buku-go/docs"
)

// @title CRUD Buku API
// @version 1.0
// @description This is a simple CRUD API for managing books.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http
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

	models.SeedData() //

	router := routes.SetupRoutes() //

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Port default
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%s", appPort)

	log.Printf("📚 Dokumentasi API (Swagger UI) tersedia di http://localhost:%s/api/doc/", appPort)
	log.Printf("🚀 Server berjalan di http://localhost:%s", appPort)
	log.Printf("🌐 Server juga dapat diakses di LAN pada alamat IP mesin Anda dengan port %s", appPort)
	log.Printf("Tekan CTRL+C untuk menghentikan server.")

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server berjalan di %s", serverAddr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}


//http://localhost:8080/api/doc/