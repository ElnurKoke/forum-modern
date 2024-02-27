package handler

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func (h *Handler) postPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id == 0 || err != nil {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	userValue := r.Context().Value("user")
	if userValue == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	post, err := h.Service.ServicePostIR.GetPostId(id)
	if err != nil {
		h.ErrorPage(w, models.ErrPostNotFound.Error(), http.StatusNotFound)
		log.Println(err.Error())
		return
	}
	comments, err := h.Service.GetCommentsByIdPost(id)
	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	categories, err := h.Service.ServicePostIR.GetCategories()
	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		sort.Sort(ByCreatedAtCom(comments))
		model := models.Info{
			User:        user,
			Post:        post,
			Comment:     comments,
			AllCategory: categories,
		}
		if err := h.Temp.ExecuteTemplate(w, "post.html", model); err != nil {
			log.Println(err.Error())
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, "Bad request", http.StatusBadRequest)
			return
		}
		commentText := r.FormValue("text")

		commentid, err := h.Service.CommentServiceIR.CreateComment(id, user.Id, commentText)
		if err != nil {
			h.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = h.Service.ServiceMsgIR.CreateMassage(models.Message{
			PostId: post.Id, CommentId: commentid, ReactAuthor: user.Username, Author: post.Author, Message: "cc",
		}, commentText); err != nil {
			h.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.URL.Path+fmt.Sprintf("/?id=%d", id), http.StatusSeeOther)
	default:
		h.ErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
