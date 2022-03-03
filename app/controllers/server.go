package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sampleapp/app/models"
	"sampleapp/config"
	"strconv"
)

//layoutの共通化
//第１引数はレスポンスライター、第２引数はHTMLに入力したいデータ、第３引数はHTMLファイルを渡す
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	//引数で渡されたファイルを配列に格納する
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}


	templates := template.Must(template.ParseFiles(files...))
	//第２引数はHTMLのdefineで呼ばれたものを明示的に宣言する
	templates.ExecuteTemplate(w, "layout", data)
}

//URLの正規表現のパターン
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")
//URLからIDを取得する処理
func ParseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//validPathとURLがマッチしたものをスライスで取得
		q := validPath.FindStringSubmatch(r.URL.Path)
		//URLが間違っていたらNotFoundを返す
		if q == nil {
			http.NotFound(w, r)
			return
		}

		//Pathをint型で受け取る
		qi, err := strconv.Atoi(q[2])
		//エラーがあったらNotFoundを返す
		if err != nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, qi)
	}
}

//cookieを取得する関数
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		//UUIDをもとにセッションを探す
		sess = models.Session{UUID: cookie.Value}
		//セッションが存在しなかったらエラーを返す
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")
		}
	}

	return sess, err
}

//サーバーの立ち上げ
func StartMainServer() error {
	//CSSとJSを読み込む
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static", files))
	//URLの登録
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	http.HandleFunc("/todos/edit/", ParseURL(todoEdit))
	http.HandleFunc("/todos/update/", ParseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", ParseURL(todoDelete))
	return http.ListenAndServe(":" + config.Config.Port, nil)
}