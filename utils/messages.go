package utils

import (
	"database/sql"
)

type Message struct {
	ID        int
	Content   string
	CreatedAt string
	Sender    string
}

func AddMessage(db *sql.DB, senderID, content string) error {
	_, err := db.Exec("INSERT INTO messages(sender_id, content) VALUES (?, ?)", senderID, content)
	return err
}

func GetAllMessages(db *sql.DB) ([]Message, error) {
	rows, err := db.Query(`SELECT m.id, m.content, m.created_at, u.username FROM messages m JOIN users u ON m.sender_id = u.id ORDER BY m.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt, &msg.Sender); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func DeleteMessage(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM messages WHERE id=?", id)
	return err
}
