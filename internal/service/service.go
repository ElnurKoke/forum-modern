package service

import "forum/internal/storage"

type Service struct {
	Auth
	ServicePostIR
	User
	CommentServiceIR
	EmotionServiceIR
	ServiceMsgIR
}

func NewService(storages *storage.Storage) *Service {
	return &Service{
		Auth:             NewAuthService(storages),
		ServicePostIR:    NewPostService(storages.PostIR),
		User:             NewUserService(storages),
		CommentServiceIR: newCommentServ(storages.CommentIR),
		EmotionServiceIR: NewEmotionService(storages.ReactionIR),
		ServiceMsgIR:     NewServiceMsg(storages.NotificationIR),
	}
}
