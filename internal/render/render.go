package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"github.com/justinas/nosurf"
	"github.com/kotan519/keijiban/internal/config"
	"github.com/kotan519/keijiban/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	//tmplで指定したファイルをtへ
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	// bufへtを実行(埋め込み)
	_ = t.Execute(buf, td)

	// render the template
	// bufをwに書き込み
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the named *.page.tmpl from ./templates
	// 一致するパス名を返すGlob	pages→ページのパス名(複数)配列 ["./templates/home.page.tmpl", "./templates/about.page.tmpl"]
	pages, _ := filepath.Glob("./templates/*.page.tmpl")

	// range through all this files ending with *.page.tmpl
	// ページのパス名(1つずつ)配列要素をpageへ
	for _, page := range pages {
		// パスの最後の要素を返す(name = home.page.tmpl)
		name := filepath.Base(page)
		// template.New(名前)→PaeseFiles(ページ)に名前を入れ込む
		ts, _ := template.New(name).ParseFiles(page)

		// 一致するパス名を返すGlob
		matches, _ := filepath.Glob("./templates/*.layout.tmpl")

		if len(matches) > 0 {
			// ファイルの取り込み
			// tsはlayout.tmplを読み込んだもの
			ts, _ = ts.ParseGlob("./templates/*.layout.tmpl")
		}

		myCache[name] = ts
	}
	return myCache, nil
}
