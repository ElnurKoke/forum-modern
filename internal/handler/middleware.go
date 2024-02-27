package handler

import (
	"context"
	"forum/internal/models"
	"net/http"
	"time"
)

func (h *Handler) middleWareGetUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		c, err := r.Cookie("token")
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", models.User{})))
			return
		}
		user, err = h.Service.GetUserByToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", models.User{})))
			return
		}
		if user.ExpiresAt.Before(time.Now()) {
			if err := h.Service.DeleteToken(c.Value); err != nil {
				h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		user.IsAuth = true
		ticker := time.NewTicker(1 * time.Second)
		if r.Method != http.MethodGet {
			// Создаем новый таймер

			defer ticker.Stop() // Важно остановить таймер в конце функции, чтобы избежать утечки памяти

			// Цикл для обработки событий таймера
			<-ticker.C // Ждем события таймера

			// Продолжаем выполнение запроса
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
		} else {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
		}
	}
}
