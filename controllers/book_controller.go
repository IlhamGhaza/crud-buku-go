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
func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, books)
}

// GetBookHandler menghandle request untuk mendapatkan satu buku berdasarkan ID
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