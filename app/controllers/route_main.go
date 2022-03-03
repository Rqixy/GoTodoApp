package controllers

import (
	"log"
	"net/http"
	"sampleapp/app/models"
)

//トップページにアクセスするルートの作成
func top(w http.ResponseWriter, r *http.Request) {	//Handlerとして登録する関数
	//Pathの引数を渡したファイルを解析
	// t, _ := template.ParseFiles("app/views/templates/top.html")
	//第１引数にResponseWriter、第２引数に表示するデータを渡す
	// t.Execute(w, "Hello")
	// t.Execute(w, nil)

	//セッションを取得
	_, err := session(w, r)
	//もしログインセッションが存在していないならtopに移動する
	if err != nil {
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}
}

//Todoリストのページにアクセスするルートの作成
func index(w http.ResponseWriter, r *http.Request) {
	//セッションを取得
	sess, err := session(w, r)
	//もしログインセッションが存在していなかったらTopに移動する
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		//セッション情報からログインしているUser情報を取得する
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		
		//Userが作成したTodoリストを取得
		todos, _ := user.GetTodosByUser()

		//UserのTodosにtodosを渡す
		user.Todos = todos
		generateHTML(w, user, "layout", "private_navbar", "index")
	}
}

//Todoリストを作成するページにアクセスするルート作成
func todoNew(w http.ResponseWriter, r *http.Request) {
	//セッションを取得
	_, err := session(w, r)
	//もしログインセッションが存在していなかったらLoginページに移動する
	if err != nil {
			http.Redirect(w, r, "/login", 302)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

//作成したテキストをほぞんする
func todoSave(w http.ResponseWriter, r *http.Request) {
	//セッションを取得
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		//フォームの値を取得
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		//セッションからUserを取得する
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}

		//フォームから入力したテキストを取得
		content := r.PostFormValue("content")

		//Todoリストを作成する
		if err := user.CreateTodo(content); err != nil {
			log.Println(err)
		}

		//todosページにリダイレクトする
		http.Redirect(w, r, "/todos", 302)


	}
}

func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	//セッションを取得
	sess, err := session(w, r)
	//もしログインセッションが存在していなかったらLoginページに移動する
	if err != nil {
			http.Redirect(w, r, "/login", 302)
	} else {
		//セッションからUserを取得する
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		//引数のIDからTodoを取得する
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		// HTML生成
		generateHTML(w, t, "layout", "private_navbar", "todo_edit")
	}
}

func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	//セッションを取得
	sess, err := session(w, r)
	//もしログインセッションが存在していなかったらLoginページに移動する
	if err != nil {
			http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		//セッションからUserを取得する
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		
		//フォームから入力したテキストを取得
		content := r.PostFormValue("content")
		
		//Todoのストラクトを作成
		t := &models.Todo{ID: id, Content: content, User_ID: user.ID}
		if err := t.UpdateTodo(); err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/todos", 302)
	}
}

func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	//セッションを取得
	sess, err := session(w, r)
	//もしログインセッションが存在していなかったらLoginページに移動する
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		//セッションからUserを取得する
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		//idからtodoを取得
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}

		//Todoを削除する
		if err := t.DeleteTodo(); err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/todos", 302)
	}
}