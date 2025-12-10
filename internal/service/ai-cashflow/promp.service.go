package aicashflow

import (
	"encoding/json"
	"fmt"
)

func (s *Service) promptUserTransactionInput(input string) (string, error) {
	categoryList, err := s.rp.Category.GetAllCategories()
	if err != nil {
		return "", err
	}

	categoriesJSON, err := json.Marshal(categoryList)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`Kamu adalah asisten keuangan. Parse input transaksi berikut dan kembalikan HANYA JSON tanpa penjelasan tambahan.

	Input: "%s"

	Categories yang tersedia:
	%s

	Instruksi:
	1. Tentukan apakah ini EXPENSE (Pengeluaran) atau INCOME (Pemasukan)
	2. Ekstrak jumlah uang (amount) sebagai integer (e.g., 50000)
	3. Pilih category_id yang paling sesuai dari daftar categories
	4. Buat 'message' yang ramah dalam Bahasa Indonesia. Formatnya harus: 
	- EXPENSE: "Pengeluaran [nama item] Rp. [amount dengan titik pemisah ribuan] dengan category [nama kategori]"
	- INCOME: "Pemasukan [nama item] Rp. [amount dengan titik pemisah ribuan] dengan category [nama kategori]"
	(e.g., "Pengeluaran Nasi Padang Rp. 50.000 dengan category Makanan")

	Response format (HANYA JSON ini, tanpa markdown, kode blok, atau teks lain):
	{
	"message": "string dengan format yang diminta",
	"req_payload": {
		"type": "EXPENSE atau INCOME",
		"amount": integer,
		"category_id": integer,
		"description": "string deskripsi singkat"
	}
	}`, input, string(categoriesJSON)), nil
}
