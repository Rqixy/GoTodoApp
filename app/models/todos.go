package models

import (
	"log"
	"time"
)

type Todo struct {
	ID 		  int
	Content   string
	User_ID   int
	CreatedAt time.Time
}

//Todoリストの作成
func (u *User) CreateTodo(content string) (err error) {
	sql := `INSERT INTO todos(content, user_id, created_at) VALUES (?, ?, ?)`
	_, err = Db.Exec(sql, content, u.ID, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//1つのTodoを取得
func GetTodo(id int) (todo Todo, err error) {
	sql := `SELECT id, content, user_id, created_at FROM todos WHERE id = ?`
	todo = Todo{}

	err = Db.QueryRow(sql, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.User_ID,
		&todo.CreatedAt,
	)

	if err != nil {
		log.Fatal(err)
	}

	return todo, err
}

//全てのTodoリストを取得
func GetTodos() (todos []Todo, err error) {
	sql := `SELECT id, content, user_id, created_at FROM todos`
	//レコードを全て取得
	rows, err := Db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	//取得したレコードを一つ一つ配列に格納する
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.User_ID,
			&todo.CreatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}

		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

//特定のUserのTodoリストを取得する
func (u *User) GetTodosByUser() (todos []Todo, err error) {
	sql := `SELECT id, content, user_id, created_at FROM todos WHERE user_id = ?`
	rows, err := Db.Query(sql, u.ID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.User_ID,
			&todo.CreatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

//Todoリストの更新
func (t *Todo) UpdateTodo() (err error) {
	sql := `UPDATE todos SET content = ?, user_id = ? WHERE id = ?`
	_, err = Db.Exec(sql, t.Content, t.User_ID, t.ID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//Todoリストの削除
func (t *Todo) DeleteTodo() (err error) {
	sql := `DELETE FROM todos WHERE id = ?`
	_, err = Db.Exec(sql, t.ID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}