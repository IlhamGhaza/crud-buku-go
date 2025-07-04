package models

import (
	"crud-buku-go/config"
	"database/sql"
	"errors"
	"log"
	"time"
)

// @Description Struktur data untuk buku
// Book merepresentasikan struktur data buku
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllBooks mengambil semua buku dari database
func GetAllBooks() ([]Book, error) {
	rows, err := config.DB.Query("SELECT id, title, author, year, created_at, updated_at FROM books ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetBookByID mengambil satu buku berdasarkan ID
func GetBookByID(id int) (Book, error) {
	var book Book
	row := config.DB.QueryRow("SELECT id, title, author, year, created_at, updated_at FROM books WHERE id = $1", id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return book, errors.New("buku tidak ditemukan")
		}
		return book, err
	}
	return book, nil
}

// CreateBook menambahkan buku baru ke database
func CreateBook(book *Book) error {
	// Menggunakan goroutine untuk logging (contoh sederhana)
	// Dalam aplikasi nyata, ini bisa untuk tugas background yang lebih kompleks
	go func(b *Book) {
		log.Printf("Goroutine: Memulai proses pembuatan buku: %s", b.Title)
		// Simulasi pekerjaan tambahan
		time.Sleep(100 * time.Millisecond)
		log.Printf("Goroutine: Selesai proses pembuatan buku: %s", b.Title)
	}(book)

	query := `INSERT INTO books (title, author, year, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err := config.DB.QueryRow(query, book.Title, book.Author, book.Year, time.Now(), time.Now()).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBook memperbarui data buku di database
func UpdateBook(id int, book *Book) error {
	// Menggunakan goroutine untuk logging pembaruan
	go func(bookID int, b *Book) {
		log.Printf("Goroutine: Memulai proses pembaruan buku ID %d: %s", bookID, b.Title)
		time.Sleep(50 * time.Millisecond)
		log.Printf("Goroutine: Selesai proses pembaruan buku ID %d", bookID)
	}(id, book)

	query := `UPDATE books SET title = $1, author = $2, year = $3, updated_at = $4
	          WHERE id = $5 RETURNING updated_at`
	err := config.DB.QueryRow(query, book.Title, book.Author, book.Year, time.Now(), id).Scan(&book.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("buku tidak ditemukan untuk diperbarui")
		}
		return err
	}
	book.ID = id // Pastikan ID tetap
	return nil
}

// SearchBooks mencari buku berdasarkan query
func SearchBooks(query string) ([]Book, error) {
	searchQuery := "%" + query + "%"
	
	// First try full-text search
	rows, err := config.DB.Query(`
		SELECT id, title, author, year, created_at, updated_at 
		FROM books 
		WHERE search_vector @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(search_vector, plainto_tsquery('english', $1)) DESC
	`, query)
	
	if err != nil {
		// Fallback to LIKE search if full-text search fails
		rows, err = config.DB.Query(`
			SELECT id, title, author, year, created_at, updated_at 
			FROM books 
			WHERE LOWER(title) LIKE LOWER($1) 
			   OR LOWER(author) LIKE LOWER($1)
			   OR year::TEXT LIKE $1
			ORDER BY 
				CASE 
					WHEN LOWER(title) = LOWER($1) THEN 1
					WHEN LOWER(title) LIKE LOWER($1) || '%' THEN 2
					WHEN LOWER(title) LIKE '%' || LOWER($1) || '%' THEN 3
					ELSE 4
				END,
			title
		`, searchQuery)
		
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// DeleteBook menghapus buku dari database
func DeleteBook(id int) error {
	// Menggunakan goroutine untuk logging penghapusan
	go func(bookID int) {
		log.Printf("Goroutine: Memulai proses penghapusan buku ID %d", bookID)
		time.Sleep(50 * time.Millisecond)
		log.Printf("Goroutine: Selesai proses penghapusan buku ID %d", bookID)
	}(id)

	result, err := config.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("buku tidak ditemukan untuk dihapus")
	}
	return nil
}

// SeedData mengisi data dummy ke tabel buku jika kosong
func SeedData() {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		log.Fatalf("Gagal menghitung data buku: %v", err)
	}

	if count > 0 {
		log.Println("Data buku sudah ada, tidak perlu seeding.")
		return
	}

	log.Println("Memulai seeding data buku...")

	dummyBooks := []Book{
		{Title: "Laskar Pelangi", Author: "Andrea Hirata", Year: 2005},
		{Title: "Bumi Manusia", Author: "Pramoedya Ananta Toer", Year: 1980},
		{Title: "Negeri 5 Menara", Author: "Ahmad Fuadi", Year: 2009},
		{Title: "Ayat-Ayat Cinta", Author: "Habiburrahman El Shirazy", Year: 2004},
		{Title: "Sang Pemimpi", Author: "Andrea Hirata", Year: 2006},
		{Title: "Perahu Kertas", Author: "Dee Lestari", Year: 2009},
		{Title: "Ronggeng Dukuh Paruk", Author: "Ahmad Tohari", Year: 1982},
		{Title: "Supernova: Ksatria, Puteri, dan Bintang Jatuh", Author: "Dee Lestari", Year: 2001},
		{Title: "Ketika Cinta Bertasbih", Author: "Habiburrahman El Shirazy", Year: 2007},
		{Title: "5 cm", Author: "Donny Dhirgantoro", Year: 2005},
		{Title: "Pulang", Author: "Leila S. Chudori", Year: 2012},
		{Title: "Cantik Itu Luka", Author: "Eka Kurniawan", Year: 2002},
		{Title: "Saman", Author: "Ayu Utami", Year: 1998},
		{Title: "Gadis Kretek", Author: "Ratih Kumala", Year: 2012},
		{Title: "Laut Bercerita", Author: "Leila S. Chudori", Year: 2017},
		{Title: "Filosofi Kopi", Author: "Dee Lestari", Year: 2006},
		{Title: "Edensor", Author: "Andrea Hirata", Year: 2007},
		{Title: "Amba", Author: "Laksmi Pamuntjak", Year: 2012},
		{Title: "Orang-Orang Biasa", Author: "Andrea Hirata", Year: 2019},
		{Title: "Aroma Karsa", Author: "Dee Lestari", Year: 2018},
		{Title: "Sirkus Pohon", Author: "Andrea Hirata", Year: 2017},
	}

	for _, book := range dummyBooks {
		// Kita panggil CreateBook agar goroutine di dalamnya juga tereksekusi
		// untuk setiap data dummy, meskipun ini hanya contoh sederhana.
		err := CreateBook(&book) // Perhatikan, CreateBook mengembalikan ID, dll.
		if err != nil {
			log.Printf("Gagal seeding buku '%s': %v", book.Title, err)
		} else {
			log.Printf("Berhasil seeding buku: %s (ID: %d)", book.Title, book.ID)
		}
	}
	log.Println("Seeding data buku selesai.")
}
