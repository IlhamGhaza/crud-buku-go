package routes

import (
	"fmt"
	"log"
	"net/http"
	"crud-buku-go/controllers"

	"github.com/gorilla/mux"
)

// SetupRoutes mengkonfigurasi semua rute untuk aplikasi
func SetupRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Middleware untuk logging setiap request (contoh sederhana)
	router.Use(loggingMiddleware)

	// Rute untuk halaman utama sederhana
	router.HandleFunc("/", homeHandler).Methods("GET")

	// Rute untuk Buku
	router.HandleFunc("/api/books", controllers.GetBooksHandler).Methods("GET")
	router.HandleFunc("/api/books", controllers.CreateBookHandler).Methods("POST")
	router.HandleFunc("/api/books/{id}", controllers.GetBookHandler).Methods("GET")
	router.HandleFunc("/api/books/{id}", controllers.UpdateBookHandler).Methods("PUT")
	router.HandleFunc("/api/books/{id}", controllers.PatchBookHandler).Methods("PATCH")
	router.HandleFunc("/api/books/{id}", controllers.DeleteBookHandler).Methods("DELETE")

	log.Println("Rute API telah diinisialisasi.")
	return router
}

// loggingMiddleware adalah contoh middleware sederhana untuk mencatat request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.RequestURI)
		// Panggil handler berikutnya dalam chain
		next.ServeHTTP(w, r)
		// Kode setelah handler (jika ada)
		// log.Println("Selesai memproses request") // Bisa ditambahkan jika perlu
	})
}

// homeHandler menangani request ke rute root ("/") dan menampilkan halaman HTML sederhana.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>Selamat Datang di API CRUD Buku!</h1>")
	fmt.Fprintln(w, "<p>Server API berjalan dengan baik.</p>")
	fmt.Fprintln(w, "<p>Anda dapat mengakses endpoint buku di <code>/api/books</code>.</p>")
}