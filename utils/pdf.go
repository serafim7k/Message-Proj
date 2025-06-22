package utils

import (
	"database/sql"
	"os"
)

type PDF struct {
	ID        int
	Filename  string
	CreatedAt string
	Sender    string
}

func AddPDF(db *sql.DB, senderID, filename string) error {
	_, err := db.Exec("INSERT INTO pdfs(sender_id, filename) VALUES (?, ?)", senderID, filename)
	return err
}

func GetAllPDFs(db *sql.DB) ([]PDF, error) {
	rows, err := db.Query(`SELECT p.id, p.filename, p.created_at, u.username FROM pdfs p JOIN users u ON p.sender_id = u.id ORDER BY p.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pdfs []PDF
	for rows.Next() {
		var p PDF
		if err := rows.Scan(&p.ID, &p.Filename, &p.CreatedAt, &p.Sender); err == nil {
			pdfs = append(pdfs, p)
		}
	}
	return pdfs, nil
}

func DeletePDF(db *sql.DB, pdfID string) error {
	var filename string
	err := db.QueryRow("SELECT filename FROM pdfs WHERE id=?", pdfID).Scan(&filename)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM pdfs WHERE id=?", pdfID)
	if err != nil {
		return err
	}
	if filename != "" {
		_ = os.Remove("uploads/pdfs/" + filename)
	}
	return nil
}
