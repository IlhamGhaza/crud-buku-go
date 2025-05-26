package routes

import (
	"crud-buku-go/controllers"
	"fmt"
	"log"
	"net/http"

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
	// Menggunakan fmt.Fprintf untuk menulis multiple lines dengan lebih mudah
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Selamat Datang di API CRUD Buku</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f7f6;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            text-align: center;
        }
        .container {
            background-color: #ffffff;
            padding: 40px 60px;
            border-radius: 10px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 90%;
        }
        h1 {
            color: #2c3e50;
            font-size: 2.5em;
            margin-bottom: 20px;
        }
        p {
            font-size: 1.1em;
            line-height: 1.6;
            margin-bottom: 15px;
        }
        code {
            background-color: #e8f0fe;
            color: #1967d2;
            padding: 3px 6px;
            border-radius: 4px;
            font-family: 'Courier New', Courier, monospace;
        }
        .footer {
            margin-top: 30px;
            font-size: 0.9em;
            color: #7f8c8d;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Selamat Datang di API CRUD Buku!</h1>
        <p>Server API ini berjalan dengan baik dan siap melayani permintaan Anda.</p>
        <p>Anda dapat mulai berinteraksi dengan data buku melalui endpoint di <code>/api/books</code>.</p>
        <p class="footer">Dikembangkan dengan Go.</p>
    </div>
</body>
</html>`)
}
