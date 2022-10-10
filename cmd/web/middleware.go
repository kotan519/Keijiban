package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf すべてのPOSTリクエストにCSRF保護機能を追加
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad リクエストごとにセッションをロードし、保存します。
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
