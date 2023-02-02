package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kotan519/keijiban/internal/config"
	"github.com/kotan519/keijiban/internal/driver"
	"github.com/kotan519/keijiban/internal/forms"
	"github.com/kotan519/keijiban/internal/helpers"
	"github.com/kotan519/keijiban/internal/models"
	"github.com/kotan519/keijiban/internal/render"
	"github.com/kotan519/keijiban/internal/repository"
	"github.com/kotan519/keijiban/internal/repository/dbrepo"
	"golang.org/x/crypto/bcrypt"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// NewRepo create the new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}


// 匿名掲示板記入ハンドラー

func (m *Repository) WriteThread(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "write-thread-tokumei.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// 書いたスレッドをデータベースにインサート
func (m *Repository) PostWriteThreadsData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	threadsdata := models.TokumeiPostData{
		Title:      r.Form.Get("name"),
		Text:       r.Form.Get("body"),
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("name", "body")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["threadsdata"] = threadsdata

		render.RenderTemplate(w, r, "write-thread-tokumei.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertData(threadsdata)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "threadsdata", threadsdata)

	http.Redirect(w, r, "/auth/threadlist", http.StatusSeeOther)
}

func (m *Repository) ThreadListScreen(w http.ResponseWriter, r *http.Request){
	threadsdata, err := m.DB.GetThreadList()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["threadsdata"] = threadsdata

	render.RenderTemplate(w, r, "threads-list.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// スレッド番号を受け取る
func (m *Repository) PostThreadNumber(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	threadnumber, err := strconv.Atoi(r.Form.Get("number"))

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "番号を入力してください")
		return
	}

	

	number := models.TokumeiPostDataNumber{
		ThreadNumber: threadnumber,
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("number")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["number"] = number
		render.RenderTemplate(w, r, "threads-list.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "number", number)

	http.Redirect(w, r, "/auth/thread", http.StatusSeeOther)
}

func (m *Repository) ThreadScreen(w http.ResponseWriter, r *http.Request){
	threadnumberdata, ok := m.App.Session.Get(r.Context(), "number").(models.TokumeiPostDataNumber)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "セッションを受け取れませんでした")
		http.Redirect(w, r, "/auth/threadlist", http.StatusTemporaryRedirect)
		return
	}

	var emptyCommentData models.TokumeiPostData
	comdata := make(map[string]interface{})
	comdata["commentsdata"] = emptyCommentData

	threadnumber := threadnumberdata.ThreadNumber
	
	threaddata, err1 := m.DB.GetThreadData(threadnumber)

	commentdata, err2 := m.DB.GetCommentData(threadnumber)

	for _, value := range commentdata {
		threaddata = append(threaddata, value)
	}


	if err1!= nil {
        helpers.ServerError(w, err1)
	}

	if err2!= nil {
        helpers.ServerError(w, err2)
	}

	data := make(map[string]interface{})
	data["threaddata"] = threaddata

	render.RenderTemplate(w, r, "thread.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostWriteComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	threadId, _ := strconv.Atoi(r.Form.Get("thread_id"))

	commentdata := models.TokumeiPostData{
		ThreadID: 	threadId,
		Title:      r.Form.Get("name"),
		Text:       r.Form.Get("body"),
	}

	form := forms.New(r.PostForm)

	form.Required("name")
	form.Required("body")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["commentdata"] = commentdata

		render.RenderTemplate(w, r, "thread.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertCommentData(commentdata)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "commentdata", commentdata)

	http.Redirect(w, r, "/auth/thread", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.EmailKosen("email")

	if !form.Valid() {
		render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "ログインできませんでした")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "ログインしました")
	http.Redirect(w, r, "/auth/threadlist", http.StatusSeeOther)

}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) Signup (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signup.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostSignup (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	userdata := models.User{
		UserName: 	r.Form.Get("username"),
		Email: 	email,
		Password: 	string(hashedPassword),
		AccressLevel: 1,
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("username", "email", "password")
	form.EmailKosen("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["userdata"] = userdata

		render.RenderTemplate(w, r, "signup.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertUserData(userdata)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "userdata", userdata)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}


