package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/kotan519/keijiban/internal/config"
)

var app *config.AppConfig


func NewHelpers(a *config.AppConfig) {
	app = a
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s¥n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exits := app.Session.Exists(r.Context(), "user_id")
	return exits
}