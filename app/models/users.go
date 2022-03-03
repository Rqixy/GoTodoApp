package models

import (
	"log"
	"time"

)

type User struct {
	ID 		  int
	UUID 	  string
	Name 	  string
	Email 	  string
	Password  string
	CreatedAt time.Time
	Todos     []Todo
}

type Session struct {
	ID 		  int
	UUID 	  string
	Email 	  string
	UserID 	  int
	CreatedAt time.Time
}

//ユーザーの作成
func (u *User) CreateUser() (err error) {
	sql := `INSERT INTO users(uuid, name, email, password, created_at) values (?, ?, ?, ?, ?)`
	_, err = Db.Exec(sql, 
		createUUID(), 
		u.Name, 
		u.Email, 
		Encrypt(u.Password), 
		time.Now(),
	)
		
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//Userの取得
func GetUser(id int) (user User, err error) {
	//Userの構造体を取得
	user = User{}
	sql := `SELECT id, uuid, name, email, password, created_at FROM users WHERE id = ?`
	//レコードを1つ取得する
	err = Db.QueryRow(sql, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	return user, err
}

//Userのデータ更新
func (u *User) UpdateUser() (err error) {
	sql := `UPDATE users SET  name = ?, email = ? WHERE id = ?`
	_, err = Db.Exec(sql, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//Userの削除
func (u *User) DeleteUser() (err error) {
	sql := `DELETE FROM users WHERE id = ?`
	_, err = Db.Exec(sql, u.ID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//メールアドレスからUserを取得
func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	sql := `SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?`
	err = Db.QueryRow(sql, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	return user, err
}

//ログイン情報を保持するためにセッションを保存する
func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	sql1 := `INSERT INTO sessions (uuid, email, user_id, created_at) VALUES (?, ?, ?, ?)`
	_, err = Db.Exec(sql1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	sql2 := `SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = ? AND email = ?`
	err = Db.QueryRow(sql2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
	}

	return session, err
}

//セッションがDBに存在するかチェック
func (sess *Session) CheckSession() (valid bool, err error) {
	sql := `SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = ?`
	err = Db.QueryRow(sql, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt,
	)
	if err != nil {
		valid = false
		return
	}

	valid = true
	
	return valid, err
}

//DBからセッションを削除する
func (sess *Session) DeleteSessionByUUID() (err error) {
	sql := `DELETE FROM sessions WHERE uuid = ?`
	_, err = Db.Exec(sql, sess.UUID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//セッションからUser情報を取得する
func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	sql := `SELECT id, uuid, name, email, created_at FROM users WHERE id = ?`
	err = Db.QueryRow(sql, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	return user, err
}