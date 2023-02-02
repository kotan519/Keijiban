package main

import (
	"net/http"

	"github.com/kotan519/keijiban/internal/helpers"
)

// SessionLoad リクエストごとにセッションをロードし、保存します。
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "先にログインしてください")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}) 
}
