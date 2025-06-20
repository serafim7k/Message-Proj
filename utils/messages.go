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
	rows, err := db.Query(`SELECT messages.id, messages.content, messages.created_at, users.username FROM messages JOIN users ON messages.sender_id = users.id ORDER BY messages.created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Content, &m.CreatedAt, &m.Sender); err == nil {
			messages = append(messages, m)
		}
	}
	return messages, nil
}

func DeleteMessage(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}
