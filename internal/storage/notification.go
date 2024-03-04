package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"time"
)

type NotificationIR interface {
	CreateMassageComment(mes models.Message) error
	CreateMassagePost(mes models.Message) error
	GetMessagesByAuthorId(id int) ([]models.Message, error)
	// MessageExists(author string, message string, postid, commentid int) (bool, error)
	// UpdateMessageCreationTime(author string, message string, createdAt time.Time, postid, commentid int) error
	GetMessagesByReactAuthorId(id int) ([]models.Message, error)
}

type NotificationStorage struct {
	db *sql.DB
}

func NewNotificationStorage(db *sql.DB) NotificationIR {
	return &NotificationStorage{
		db: db,
	}
}

func (ns *NotificationStorage) CreateMassageComment(mes models.Message) error {
	mes.CreateAt = time.Now()
	query := `INSERT INTO notification(post_id, comment_id, to_user_id , from_user_id, message, created_at) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := ns.db.Exec(query, mes.PostId, mes.CommentId, mes.ToUserId, mes.FromUserId, mes.Message, mes.CreateAt)
	if err != nil {
		return err
	}
	return nil
}

func (ns *NotificationStorage) CreateMassagePost(mes models.Message) error {
	mes.CreateAt = time.Now()
	query := `INSERT INTO notification(post_id,  to_user_id , from_user_id,  message, created_at) VALUES ($1, $2, $3, $4, $5);`
	_, err := ns.db.Exec(query, mes.PostId, mes.ToUserId, mes.FromUserId, mes.Message, mes.CreateAt)
	if err != nil {
		return err
	}
	return nil
}

func (ns *NotificationStorage) GetMessagesByAuthorId(id int) ([]models.Message, error) {
	query := `SELECT n.id,
			n.from_user_id ,
			n.to_user_id ,
			u.username,
			u.imageURL,
			n.post_id ,
			p.imageURL,
			n.comment_id ,
			n.message ,
			n.activity,
			n.created_at
		FROM notification n
		LEFT JOIN user u
		ON u.id = n.from_user_id  
		LEFT JOIN post p
		ON p.id = n.post_id 
		WHERE n.to_user_id = ?`
	rows, err := ns.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Error querying the database: %v", err)
	}

	var messages []models.Message

	for rows.Next() {
		var message models.Message
		var commentID sql.NullInt64
		err := rows.Scan(&message.Id,
			&message.ToUserId,
			&message.FromUserId,
			&message.FromUserName,
			&message.AvaImage,
			&message.PostId,
			&message.PostImage,
			&commentID,
			&message.Message,
			&message.Active,
			&message.CreateAt)
		if err != nil {
			models.ErrLog.Printf("Error scanning rows: %v", err)
			return nil, errors.New("storage sql cant scan notification table")
		}
		if commentID.Valid {
			message.CommentId = int(commentID.Int64)
		}
		message.Message = ConvertMessageAuthor(message.Message)
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return messages, nil
}

func (ns *NotificationStorage) GetMessagesByReactAuthorId(id int) ([]models.Message, error) {
	query := `SELECT n.id,
			n.from_user_id ,
			n.to_user_id ,
			u.username,
			u.imageURL,
			n.post_id INT,
			p.imageURL,
			n.comment_id ,
			n.message ,
			n.activity,
			n.created_at
		FROM notification n
		LEFT JOIN user u
		ON u.id = n.to_user_id  
		LEFT JOIN post p
		ON p.id = n.post_id 
		WHERE n.from_user_id = ?`
	rows, err := ns.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Error querying the database: %v", err)
	}

	var messages []models.Message

	for rows.Next() {
		var message models.Message
		var commentID sql.NullInt64
		err := rows.Scan(&message.Id,
			&message.ToUserId,
			&message.FromUserId,
			&message.FromUserName,
			&message.AvaImage,
			&message.PostId,
			&message.PostImage,
			&commentID,
			&message.Message,
			&message.Active,
			&message.CreateAt)
		if err != nil {
			models.ErrLog.Printf("Error scanning rows: %v", err)
			return nil, errors.New("storage sql cant scan notification table")
		}
		if commentID.Valid {
			message.CommentId = int(commentID.Int64)
		}
		message.Message = ConvertMessageAction(message.Message)
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return messages, nil
}

func ConvertMessageAction(mes string) string {
	switch mes {
	case "pl":
		return "You liked post"
	case "pd":
		return "You disliked post"
	case "cl":
		return "You liked comment"
	case "cd":
		return "You disliked comment"
	case "cc":
		return "You create comment"
	default:
		return "not have code message"
	}
}

func ConvertMessageAuthor(mes string) string {
	switch mes {
	case "pl":
		return "user liked your post"
	case "pd":
		return "user disliked your post"
	case "cl":
		return "user liked your comment"
	case "cd":
		return "user disliked your comment"
	case "cc":
		return "user create comment on your post"
	default:
		return "not have code message"
	}
}

// func (ns *NotificationStorage) MessageExists(author string, message string, postid, commentid int) (bool, error) {
// 	query := `SELECT COUNT(*) FROM notification WHERE author = $2 AND message = $3 AND post_id = $4 AND comment_id = $5;`
// 	var count int
// 	err := ns.db.QueryRow(query, author, message, postid, commentid).Scan(&count)
// 	if err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

// func (ns *NotificationStorage) UpdateMessageCreationTime(author string, message string, createdAt time.Time, postid, commentid int) error {
// 	updateQuery := `UPDATE notification SET created_at = $1 WHERE author = $2 AND message = $3 AND post_id = $4 AND comment_id = $5;`
// 	result, err := ns.db.Exec(updateQuery, createdAt, author, message, postid, commentid)
// 	if err != nil {
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("Record with author %s and message %s does not exist", author, message)
// 	}

// 	return nil
// }
