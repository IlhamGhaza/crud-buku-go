package routes

import (
	"crud-buku-go/controllers"
	_ "crud-buku-go/docs"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	router.PathPrefix("/api/doc/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/", homeHandler).Methods("GET")

	// Book routes
	bookRouter := router.PathPrefix("/api/books").Subrouter()
	bookRouter.HandleFunc("", controllers.GetBooksHandler).Methods("GET")
	bookRouter.HandleFunc("", controllers.CreateBookHandler).Methods("POST")
	bookRouter.HandleFunc("/search", controllers.SearchBooksHandler).Methods("GET")
	bookRouter.HandleFunc("/{id}", controllers.GetBookHandler).Methods("GET")
	bookRouter.HandleFunc("/{id}", controllers.UpdateBookHandler).Methods("PUT")
	bookRouter.HandleFunc("/{id}", controllers.PatchBookHandler).Methods("PATCH")
	bookRouter.HandleFunc("/{id}", controllers.DeleteBookHandler).Methods("DELETE")

	log.Println("Rute Swagger UI telah diinisialisasi di /api/doc/")
	log.Println("Rute API telah diinisialisasi.")
	return router
}

// logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)

	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
        <p>Anda dapat mulai berinteraksi dengan data buku melalui endpoint di <code>/api/doc/</code>.</p>
        <p class="footer">Dikembangkan dengan Go.</p>
    </div>
</body>
</html>`)
}
