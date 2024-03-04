package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	Id             int
	Email          string
	Username       string
	Password       string
	RepeatPassword string
	ExpiresAt      time.Time
	IsAuth         bool
	ImageBack      string
	ImageURL       string
	Rol            string
	Bio            string
	Created_at     time.Time
	Updated_at     time.Time
}

type GoogleLoginUserData struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      string
	Photo     string
	Verified  bool
	Provider  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GithubUserData struct {
	Login  string `json:"login"`
	ID     int    `json:"id"`
	NodeID string `json:"node_id"`
}

type Post struct {
	Id          int
	Title       string
	Description string
	Image       string
	Category    []string
	UserId      int
	Author      string
	Likes       int
	Dislikes    int
	CreateAt    time.Time
}

type Message struct {
	Id            int
	PostId        int
	CommentId     int
	FromUserId    int
	ToUserId      int
	Author        string
	ReactAuthor   string
	Message       string
	Active        int
	FromUserName  string
	AvaImage      string
	PostImage     string
	FromUserImage string
	CreateAt      time.Time
}

type Comment struct {
	Id         int
	PostId     int
	UserId     int
	Creator    string
	Text       string
	Likes      int
	Dislikes   int
	IsAuth     bool
	Created_at time.Time
}

type Like struct {
	UserID       int
	PostID       int
	Islike       int
	CommentID    int
	CountLike    int
	Countdislike int
}

type Category struct {
	Name string
}

type Communication struct {
	Id            int
	FromUserId    int
	FromUserName  string
	ForWhomRole   string
	OldRole       string
	NewRole       string
	AboutUserId   int
	AboutUserName string
	PostId        int
	PostImage     string
	CommentId     int
	CommentText   string
	Message       string
	CreatedAt     time.Time
}
