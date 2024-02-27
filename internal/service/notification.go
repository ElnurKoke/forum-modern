package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/storage"
	"time"
)

type ServiceMsgIR interface {
	CreateMassage(post models.Message, comm string) error
	GetMessagesByAuthor(author string) ([]models.Message, error)
	GetMessagesByReactAuthor(rauthor string) ([]models.Message, error)
}

type MsgService struct {
	storage storage.NotificationIR
}

func NewServiceMsg(NotificationIR storage.NotificationIR) ServiceMsgIR {
	return &MsgService{
		storage: NotificationIR,
	}
}

func (m *MsgService) CreateMassage(mes models.Message, comm string) error {
	if mes.Message == "cl" {
		mes.Message = fmt.Sprintf(" %s liked comment: \"%s\" . Which was created by \"%s\"", mes.ReactAuthor, comm, mes.Author)
	} else if mes.Message == "cd" {
		mes.Message = fmt.Sprintf(" %s disliked comment: \"%s\". Which was created by \"%s\"", mes.ReactAuthor, comm, mes.Author)
	} else if mes.Message == "pl" {
		mes.Message = fmt.Sprintf(" %s loved post: \"%s\" . Which was created by \"%s\"", mes.ReactAuthor, comm, mes.Author)
	} else if mes.Message == "pd" {
		mes.Message = fmt.Sprintf("Oh no! %s disliked post: \"%s\". Which was created by \"%s\"", mes.ReactAuthor, comm, mes.Author)
	} else if mes.Message == "cc" {
		mes.Message = fmt.Sprintf(" %s commented on post: \"%s\". Which was created by \"%s\"", mes.ReactAuthor, comm, mes.Author)
	}
	exists, err := m.storage.MessageExists(mes.Author, mes.Message, mes.PostId, mes.CommentId)
	if err != nil {
		return err
	}
	if exists {
		return m.storage.UpdateMessageCreationTime(mes.Author, mes.Message, time.Now(), mes.PostId, mes.CommentId)
	} else {
		return m.storage.CreateMassage(mes)
	}
}

func (m *MsgService) GetMessagesByAuthor(author string) ([]models.Message, error) {
	return m.storage.GetMessagesByAuthor(author)
}

func (m *MsgService) GetMessagesByReactAuthor(rauthor string) ([]models.Message, error) {
	return m.storage.GetMessagesByReactAuthor(rauthor)
}
