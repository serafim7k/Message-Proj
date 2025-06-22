package utils

import (
	"database/sql"
	"os"
)

type Music struct {
	ID        int
	Filename  string
	CreatedAt string
	Sender    string
}

func AddMusic(db *sql.DB, senderID, filename string) error {
	_, err := db.Exec("INSERT INTO music(sender_id, filename) VALUES (?, ?)", senderID, filename)
	return err
}

func GetAllMusic(db *sql.DB) ([]Music, error) {
	rows, err := db.Query(`SELECT m.id, m.filename, m.created_at, u.username FROM music m JOIN users u ON m.sender_id = u.id ORDER BY m.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var music []Music
	for rows.Next() {
		var m Music
		if err := rows.Scan(&m.ID, &m.Filename, &m.CreatedAt, &m.Sender); err == nil {
			music = append(music, m)
		}
	}
	return music, nil
}

func DeleteMusic(db *sql.DB, musicID string) error {
	var filename string
	err := db.QueryRow("SELECT filename FROM music WHERE id=?", musicID).Scan(&filename)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM music WHERE id=?", musicID)
	if err != nil {
		return err
	}
	if filename != "" {
		_ = os.Remove("uploads/music/" + filename)
	}
	return nil
}
