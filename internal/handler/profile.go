package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) profilePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile/" {
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
	if !user.IsAuth {
		h.ErrorPage(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	profileUser := models.User{}

	posts, err := h.Service.ServicePostIR.GetMyPost(user.Id)
	if err != nil {
		h.ErrorPage(w, models.ErrPostNotFound.Error(), http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:

		model := models.ProfileInfo{
			User:        user,
			ProfileUser: profileUser,
			Posts:       posts,
		}
		if err := h.Temp.ExecuteTemplate(w, "profile.html", model); err != nil {
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
		username := r.FormValue("username")

		if err := h.Service.User.UpdateUserName(user.Id, username); err != nil {
			info := models.ProfileInfo{
				Error:       err.Error(),
				User:        user,
				ProfileUser: profileUser,
				Posts:       posts,
			}
			if err := h.Temp.ExecuteTemplate(w, "profile.html", info); err != nil {
				h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}
		user.Username = username
		info := models.ProfileInfo{
			Error:       "You have successfully update name",
			User:        user,
			ProfileUser: profileUser,
			Posts:       posts,
		}
		if err := h.Temp.ExecuteTemplate(w, "profile.html", info); err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		h.ErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
