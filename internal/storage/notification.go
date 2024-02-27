package storage

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"time"
)

type NotificationIR interface {
	CreateMassage(post models.Message) error
	GetMessagesByAuthor(author string) ([]models.Message, error)
	MessageExists(author string, message string, postid, commentid int) (bool, error)
	UpdateMessageCreationTime(author string, message string, createdAt time.Time, postid, commentid int) error
	GetMessagesByReactAuthor(rauthor string) ([]models.Message, error)
}

type NotificationStorage struct {
	db *sql.DB
}

func NewNotificationStorage(db *sql.DB) NotificationIR {
	return &NotificationStorage{
		db: db,
	}
}

func (ns *NotificationStorage) CreateMassage(mes models.Message) error {
	mes.CreateAt = time.Now()
	query := `INSERT INTO notification(post_id, comment_id, author, reactauthor, message, created_at) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := ns.db.Exec(query, mes.PostId, mes.CommentId, mes.Author, mes.ReactAuthor, mes.Message, mes.CreateAt)
	if err != nil {
		return err
	}
	return nil
}

func (ns *NotificationStorage) GetMessagesByAuthor(author string) ([]models.Message, error) {
	query := "SELECT * FROM notification WHERE author = ?"
	rows, err := ns.db.Query(query, author)
	if err != nil {
		return nil, fmt.Errorf("Error querying the database: %v", err)
	}

	var messages []models.Message

	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.Id, &message.PostId, &message.CommentId, &message.Author, &message.ReactAuthor, &message.Message, &message.Active, &message.CreateAt)
		if err != nil {
			return nil, fmt.Errorf("Error scanning rows: %v", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return messages, nil
}

func (ns *NotificationStorage) GetMessagesByReactAuthor(rauthor string) ([]models.Message, error) {
	query := "SELECT * FROM notification WHERE reactauthor = ?"
	rows, err := ns.db.Query(query, rauthor)
	if err != nil {
		return nil, fmt.Errorf("Error querying the database: %v", err)
	}

	var messages []models.Message

	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.Id, &message.PostId, &message.CommentId, &message.Author, &message.ReactAuthor, &message.Message, &message.Active, &message.CreateAt)
		if err != nil {
			return nil, fmt.Errorf("Error scanning rows: %v", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return messages, nil
}

func (ns *NotificationStorage) MessageExists(author string, message string, postid, commentid int) (bool, error) {
	query := `SELECT COUNT(*) FROM notification WHERE author = $2 AND message = $3 AND post_id = $4 AND comment_id = $5;`
	var count int
	err := ns.db.QueryRow(query, author, message, postid, commentid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ns *NotificationStorage) UpdateMessageCreationTime(author string, message string, createdAt time.Time, postid, commentid int) error {
	updateQuery := `UPDATE notification SET created_at = $1 WHERE author = $2 AND message = $3 AND post_id = $4 AND comment_id = $5;`
	result, err := ns.db.Exec(updateQuery, createdAt, author, message, postid, commentid)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Record with author %s and message %s does not exist", author, message)
	}

	return nil
}
