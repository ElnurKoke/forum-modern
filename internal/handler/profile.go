package handler

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/storage"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	asks, err := h.Service.CommunicationServiceIR.GetAllAsks(user.Rol)
	if err != nil {
		h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		models.ErrLog.Println(err)
		return
	}
	model := models.ProfileInfo{
		User:        user,
		ProfileUser: profileUser,
		Posts:       posts,
		Askeds:      asks,
	}
	if user.Rol == "king" {
		alluser, err := h.Service.User.GetAllUser(user.Id)
		if err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			models.ErrLog.Println(err)
			return
		}
		model.AllUsers = alluser
	} else if user.Rol == "moderator" {
		waitPosts, err := h.Service.ServicePostIR.GetAllWaitPosts()
		if err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			models.ErrLog.Println(err)
			return
		}
		model.WaitPosts = waitPosts
	}
	switch r.Method {
	case http.MethodGet:

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
		if r.FormValue("form") == "username" { //---------------------------------------------------------------------name edit
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
		} else if r.FormValue("form") == "email" { //------------------------------------------------------------------email edit
		} else if r.FormValue("form") == "crPost" { //------------------------------------------------------------------email edit
			if user.Rol != "moderator" {
				h.ErrorPage(w, "your role is not suitable", http.StatusBadRequest)
				return
			}
			res := r.Form.Get("isCrPost")
			inf := strings.Split(res, ",")
			action := inf[0]
			strId := inf[1]
			post_id, err := strconv.Atoi(strId)
			if err != nil {
				models.ErrLog.Println(" Error strconv.Atoi: ", post_id)
				return
			}
			if err := h.Service.CommunicationServiceIR.ConfirmPost(post_id, action); err != nil {
				h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			link := fmt.Sprintf("/profile/?id=%d", user.Id)
			http.Redirect(w, r, link, http.StatusSeeOther)
		} else if r.FormValue("form") == "role" { //---------------------------------------------------------------------role ask
			res := r.Form.Get("isLevelUp")
			if res == "isLevelUp" {
				if err := h.Service.CommunicationServiceIR.AskRole(models.Communication{FromUserId: user.Id, OldRole: user.Rol}); err != nil {
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
				info := models.ProfileInfo{
					Error:       "Your request for a role upgrade has been sent",
					User:        user,
					ProfileUser: profileUser,
					Posts:       posts,
				}
				if err := h.Temp.ExecuteTemplate(w, "profile.html", info); err != nil {
					h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				models.ErrLog.Println("such control with a role is not provided")
				return
			}
		} else if r.FormValue("form") == "bio" { //---------------------------------------------------------------------bio edit
		} else if r.FormValue("form") == "ava" { //---------------------------------------------------------------------ava edit
		} else if r.FormValue("form") == "changeRole" { //--------------------------------------------------------------rol change
			if user.Rol != "king" {
				h.ErrorPage(w, "your role is not suitable", http.StatusBadRequest)
				return
			}
			res := r.Form.Get("isLevel")
			inf := strings.Split(res, ",")
			strId := inf[2]
			oldRole := inf[1]
			updown := inf[0]
			var newRole string
			if updown == "up" {
				newRole = storage.UpRole(oldRole)
			} else {
				newRole = storage.DownRole(oldRole)
			}

			id, err := strconv.Atoi(strId)
			if err != nil {
				models.ErrLog.Println(" Error strconv.Atoi: ", id)
				return
			}

			if err := h.Service.CommunicationServiceIR.UpUserRole(id, newRole); err != nil {
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
			if err := h.Service.ServiceMsgIR.CreateMassageUpRole(models.Message{FromUserId: user.Id, ToUserId: id, Message: "upRole"}); err != nil {
				h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			link := fmt.Sprintf("/profile/?id=%d", user.Id)
			http.Redirect(w, r, link, http.StatusSeeOther)
		} else if r.FormValue("form") == "roleUp" { //---------------------------------------------------------------------up role
			res := r.Form.Get("isLevelUp")
			if strings.Contains(res, "accept") {
				id, err := strconv.Atoi(res[6:])
				if err != nil {
					models.ErrLog.Println(" Error strconv.Atoi: ", id)
					return
				}
				if err := h.Service.CommunicationServiceIR.UpUserRole(id, user.Rol); err != nil {
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
				if err := h.Service.CommunicationServiceIR.DeleteAskRole(id); err != nil {
					h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if err := h.Service.ServiceMsgIR.CreateMassageUpRole(models.Message{FromUserId: user.Id, ToUserId: id, Message: "upRole"}); err != nil {
					h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				link := fmt.Sprintf("/profile/?id=%d", user.Id)
				http.Redirect(w, r, link, http.StatusSeeOther)
			} else if strings.Contains(res, "refuse") {
				id, err := strconv.Atoi(res[6:])
				if err != nil {
					models.ErrLog.Println(" Error strconv.Atoi: ", id)
					return
				}
				if err := h.Service.CommunicationServiceIR.DeleteAskRole(id); err != nil {
					h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				link := fmt.Sprintf("/profile/?id=%d", user.Id)
				http.Redirect(w, r, link, http.StatusSeeOther)
			} else {
				models.ErrLog.Println("such control with a role is not provided")
				return
			}
		} else {
			models.ErrLog.Println("this post request has not been looked at yet")
			h.ErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	default:
		h.ErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
