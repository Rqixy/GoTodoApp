package controllers

import (
	"log"
	"net/http"
	"sampleapp/app/models"
)

//DB内のメールアドレスとフォームから取得したメールアドレスを比較する
func checkEmail(w http.ResponseWriter, r *http.Request, user models.User) (valid bool) {
	valid = true
	
	//登録されているすべてのUserを取得
	modelsUsers, err := models.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	//DBにある全てのUserのメールアドレスと
	//フォームから送られてきたメールアドレスを比較して、
	//存在していたら、メッセージを返す
	for _, modelsUser := range modelsUsers {
		if user.Email == modelsUser.Email {
			valid = false
			break
		}
	}
	
	return valid
}

//新規登録のルート作成
func signup(w http.ResponseWriter, r *http.Request) {
	//SignUpPageに移動した時に生成される
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}
	} else if r.Method == "POST" {
		//postメソッドが送信された時に実行される
		//入力フォームからの解析
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		//Userのstructを作成
		user := models.User{
			//postから受け取った値をUserのstructに入れる
			Name: r.PostFormValue("name"),
			Email: r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		//DB内のメールアドレスとフォームから取得したメールアドレスを比較する
		if ok := checkEmail(w, r, user); !ok {
			generateHTML(w, user, "layout", "public_navbar", "signup")
			return
		}

		err = user.CreateUser()
		if err != nil {
			log.Fatal(err)
		}

		//TopPageに戻る
		http.Redirect(w, r, "/", 302)
	}
}

//ログインのルート作成
func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "todos", 302)
	}
}

//ログイン確認処理
func authenticate(w http.ResponseWriter, r *http.Request) {
	//postメソッドが送信された時に実行される
	//入力フォームからの解析
	err := r.ParseForm()
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}
	//Userが存在して、入力したパスワードがあっていたらログイン成功
	if user.Password == models.Encrypt(r.PostFormValue("password")) {
		//セッションの作成
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

		//作成されたセッションをもとにしてcookieを作成
		cookie := http.Cookie{
			Name: "_cookie",
			Value: session.UUID,
			HttpOnly: true,	//xss対策
		}

		//cookieに保存
		http.SetCookie(w, &cookie)

		//ログイン成功
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

//ログアウト処理
func logout(w http.ResponseWriter, r *http.Request) {
	//ブラウザからcookieを取得
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	//cookieがあったら
	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}
	http.Redirect(w, r, "/", 302)
}