package storage

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"time"
)

type CommunicationIR interface {
	CreateCommunication(mes models.Communication) error
	AskRole(mes models.Communication) error
	GetTimeAskRole(id int) (time.Time, error)
	GetAllAsks(role string) ([]models.Communication, error)
	UpUserRole(id int, newrole string) error
	ConfirmPost(id int, action string) error
	DeleteAskRole(id int) error
}

type CommunicationStore struct {
	db *sql.DB
}

func NewCommunicationStore(db *sql.DB) CommunicationIR {
	return &CommunicationStore{
		db: db,
	}
}

func (c *CommunicationStore) CreateCommunication(mes models.Communication) error {
	query := `INSERT INTO communication(post_id, comment_id, to_user_id , from_user_id, message, created_at) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := c.db.Exec(query)
	if err != nil {
		models.ErrLog.Println(err)
		return err
	}
	return nil
}

func (c *CommunicationStore) UpUserRole(id int, newrole string) error {
	query := `UPDATE user SET rol = ? WHERE id= ?;`
	_, err := c.db.Exec(query, newrole, id)
	if err != nil {
		models.ErrLog.Println(err)
		return err
	}
	return nil
}

func (c *CommunicationStore) ConfirmPost(id int, action string) error {
	var query string
	var err error
	if action == "accept" {
		query = `UPDATE post SET status = "done" WHERE id= ?;`
		_, err = c.db.Exec(query, id)
	} else if action == "delete" {
		query = `DELETE FROM post WHERE id = $1;`
		_, err = c.db.Exec(query, id)
	} else if action == "forking" {
		query = `UPDATE post SET status = 'done' WHERE id = (SELECT MAX(id) FROM post);`
		_, err = c.db.Exec(query)
	}
	if err != nil {
		models.ErrLog.Println(err)
		return err
	}
	return nil
}

func (c *CommunicationStore) DeleteAskRole(id int) error {
	query := `DELETE FROM askrole WHERE from_user_id = $1;`
	_, err := c.db.Exec(query, id)
	if err != nil {
		models.ErrLog.Println(err)
		return err
	}
	return nil
}
func (c *CommunicationStore) AskRole(mes models.Communication) error {
	query := `INSERT INTO askrole(from_user_id, old_role, new_role , for_whom_role) VALUES ($1, $2, $3, $4);`
	mes.NewRole, mes.ForWhomRole = UpRole(mes.OldRole), UpRole(mes.OldRole)
	_, err := c.db.Exec(query, mes.FromUserId, mes.OldRole, mes.NewRole, mes.ForWhomRole)
	if err != nil {
		models.ErrLog.Println(err)
		return err
	}
	return nil
}

func UpRole(old string) string {
	switch old {
	case "user":
		return "moderator"
	case "moderator":
		return "admin"
	case "admin":
		return "admin"
	case "king":
		return "king"
	default:
		return "havent role"
	}
}
func DownRole(old string) string {
	switch old {
	case "user":
		return "user"
	case "moderator":
		return "user"
	case "admin":
		return "moderator"
	case "king":
		return "king"
	default:
		return "havent role"
	}
}

func (c *CommunicationStore) GetTimeAskRole(id int) (time.Time, error) {
	query := `SELECT created_at FROM askrole WHERE from_user_id = $1 ORDER BY created_at ASC LIMIT 1;`
	row := c.db.QueryRow(query, id)
	var timer time.Time
	err := row.Scan(&timer)
	if err != nil && err != sql.ErrNoRows {
		return timer, err
	}
	if err == nil {
		return timer, nil
	}

	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	return sevenDaysAgo, nil
}

func (c *CommunicationStore) GetAllAsks(role string) ([]models.Communication, error) {
	askMessages := []models.Communication{}
	query := `
		SELECT a.from_user_id,u.username, a.old_role, a.new_role , a.for_whom_role, a.created_at
		FROM askrole a
		LEFT JOIN user u ON a.from_user_id = u.id
		WHERE a.for_whom_role=$1;`
	row, err := c.db.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("storage: get all askMessage: %w", err)
	}
	for row.Next() {
		var communication models.Communication
		if err := row.Scan(&communication.FromUserId, &communication.FromUserName, &communication.OldRole,
			&communication.NewRole, &communication.ForWhomRole, &communication.CreatedAt); err != nil {
			return nil, fmt.Errorf("storage: get all askMessage: %w", err)
		}
		askMessages = append(askMessages, communication)
	}

	return askMessages, nil
}
