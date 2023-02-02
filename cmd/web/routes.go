package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotan519/keijiban/internal/config"
	"github.com/kotan519/keijiban/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	//パニックを吸収しスタックトレースを表示する(panic抑制)
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Route("/auth", func(mux chi.Router){
		mux.Use(Auth)
		// Threadのリスト画面
		mux.Get("/threadlist", handlers.Repo.ThreadListScreen)
		mux.Post("/threadlist", handlers.Repo.PostThreadNumber)
		// Threadの内容画面
		mux.Get("/thread", handlers.Repo.ThreadScreen)
		mux.Post("/thread", handlers.Repo.PostWriteComment)

		// Threadの書き込み
		mux.Get("/write-thread-tokumei", handlers.Repo.WriteThread)
		mux.Post("/write-thread-tokumei", handlers.Repo.PostWriteThreadsData)
	})
	

	mux.Route("/user", func(mux chi.Router){
		mux.Get("/login", handlers.Repo.ShowLogin)
		mux.Post("/login", handlers.Repo.PostShowLogin)
		mux.Get("/logout", handlers.Repo.Logout)

		mux.Get("/signup", handlers.Repo.Signup)
		mux.Post("/signup", handlers.Repo.PostSignup)
	})
	

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
