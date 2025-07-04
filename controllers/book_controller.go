package controllers

import (
	"crud-buku-go/models"
	"crud-buku-go/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetBooksHandler menghandle request untuk mendapatkan semua buku
// @Summary Mendapatkan semua buku
// @Description Mengambil daftar semua buku dari database.
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book "Daftar semua buku"
// @Router /books [get]
func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, books)
}

// GetBookHandler menghandle request untuk mendapatkan satu buku berdasarkan ID
// @Summary Mendapatkan buku berdasarkan ID
// @Description Mengambil detail buku berdasarkan ID.
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID Buku"
// @Success 200 {object} models.Book "Detail buku"
// @Failure 400 {object} map[string]string "ID buku tidak valid"
// @Failure 404 {object} map[string]string "Buku tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server internal"
// @Router /books/{id} [get]
func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	book, err := models.GetBookByID(id)
	if err != nil {
		if err.Error() == "buku tidak ditemukan" {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, book)
}

// CreateBookHandler menghandle request untuk membuat buku baru
// @Summary Membuat buku baru
// @Description Menambahkan buku baru ke database.
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Data buku baru"
// @Success 201 {object} models.Book "Buku berhasil dibuat"
// @Failure 400 {object} map[string]string "Payload request tidak valid atau data buku tidak lengkap"
// @Failure 500 {object} map[string]string "Kesalahan server internal"
// @Router /books [post]
func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Payload request tidak valid")
		return
	}
	defer r.Body.Close()

	// Validasi dasar
	if book.Title == "" || book.Author == "" || book.Year == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Judul, Penulis, dan Tahun tidak boleh kosong")
		return
	}


	if err := models.CreateBook(&book); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, book)
}

// UpdateBookHandler menghandle request untuk memperbarui buku
// @Summary Memperbarui buku
// @Description Memperbarui data buku berdasarkan ID.
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID Buku"
// @Param book body models.Book true "Data buku yang diperbarui"
// @Success 200 {object} models.Book "Buku berhasil diperbarui"
// @Failure 400 {object} map[string]string "ID buku tidak valid atau payload request tidak valid"
// @Failure 404 {object} map[string]string "Buku tidak ditemukan untuk diperbarui"
// @Failure 500 {object} map[string]string "Kesalahan server internal"
// @Router /books/{id} [put]
func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	var book models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Payload request tidak valid")
		return
	}
	defer r.Body.Close()

	// Validasi dasar
	if book.Title == "" || book.Author == "" || book.Year == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Judul, Penulis, dan Tahun tidak boleh kosong")
		return
	}


	if err := models.UpdateBook(id, &book); err != nil {
		if err.Error() == "buku tidak ditemukan untuk diperbarui" {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// Ambil data buku yang sudah terupdate untuk response
	updatedBook, err := models.GetBookByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Gagal mengambil data buku setelah update")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updatedBook)
}

// DeleteBookHandler menghandle request untuk menghapus buku
// @Summary Menghapus buku
// @Description Menghapus buku berdasarkan ID.
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID Buku"
// @Success 200 {object} map[string]string "Pesan sukses penghapusan"
// @Failure 400 {object} map[string]string "ID buku tidak valid"
// @Failure 404 {object} map[string]string "Buku tidak ditemukan untuk dihapus"
// @Failure 500 {object} map[string]string "Kesalahan server internal"
// @Router /books/{id} [delete]
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	if err := models.DeleteBook(id); err != nil {
		if err.Error() == "buku tidak ditemukan untuk dihapus" {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Buku berhasil dihapus"})
}

// PatchBookHandler menghandle request untuk memperbarui sebagian data buku
// @Summary Memperbarui sebagian data buku
// @Description Memperbarui sebagian data buku (judul, penulis, atau tahun) berdasarkan ID.
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID Buku"
// @Param book body models.Book true "Data buku yang akan diperbarui (hanya field yang ingin diubah)"
// @Success 200 {object} models.Book "Buku berhasil diperbarui (sebagian)"
// @Failure 400 {object} map[string]string "ID buku tidak valid atau payload request tidak valid"
// @Failure 404 {object} map[string]string "Buku tidak ditemukan untuk diperbarui"
// @Failure 500 {object} map[string]string "Kesalahan server internal"
// @Router /books/{id} [patch]
func PatchBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	// Ambil buku yang ada dari database
	existingBook, err := models.GetBookByID(id)
	if err != nil {
		if err.Error() == "buku tidak ditemukan" { // Sesuai dengan error dari GetBookByID
			utils.RespondWithError(w, http.StatusNotFound, "Buku tidak ditemukan untuk diperbarui")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var payloadBook models.Book // Untuk menampung data dari request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payloadBook); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Payload request tidak valid")
		return
	}
	defer r.Body.Close()

	// Terapkan pembaruan parsial
	// Hanya perbarui field jika ada di payload dan tidak kosong/nol
	updated := false
	if payloadBook.Title != "" {
		existingBook.Title = payloadBook.Title
		updated = true
	}
	if payloadBook.Author != "" {
		existingBook.Author = payloadBook.Author
		updated = true
	}
	if payloadBook.Year != 0 { // Asumsi tahun tidak boleh 0 jika diisi, konsisten dengan validasi lain
		existingBook.Year = payloadBook.Year
		updated = true
	}

	if !updated {
		// Jika tidak ada field yang valid untuk diupdate dalam payload, kembalikan data yang ada.
		utils.RespondWithJSON(w, http.StatusOK, existingBook)
		return
	}

	// Perbarui buku di database menggunakan fungsi UpdateBook yang ada
	// Asumsi: models.UpdateBook akan memperbarui semua field dari objek existingBook yang diteruskan.
	if err := models.UpdateBook(id, &existingBook); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Gagal memperbarui buku: "+err.Error())
		return
	}

	// Ambil data buku yang sudah terupdate untuk response (untuk memastikan konsistensi)
	finalUpdatedBook, err := models.GetBookByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Gagal mengambil data buku setelah update")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, finalUpdatedBook)
}

// SearchBooksHandler handles book search requests
// @Summary Search books
// @Description Search books by title, author, or year
// @Tags books
// @Produce json
// @Param q query string true "Search query (can be title, author, or year)" 
// @Success 200 {array} models.Book "List of matching books"
// @Failure 400 {object} map[string]string "Search query is required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books/search [get]
func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Search query parameter 'q' is required")
		return
	}

	books, err := models.SearchBooks(query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, books)
}