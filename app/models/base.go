package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"sampleapp/config"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB
var err error

//テーブル名を宣言
const (
	tableNameUser = "users"
	tableNameTodo = "todos"
	tableNameSession = "sessions"
)

func init() {
	//SQLの起動
	Db, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		log.Fatal(err)
	}

	//usersテーブル作成
	sqlU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		name STRING,
		email STRING,
		password STRING,
		created_at DATETIME
		)`, tableNameUser)

	Db.Exec(sqlU)

	//todosテーブル作成
	sqlT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		created_at DATETIME
	)`, tableNameTodo)

	Db.Exec(sqlT)

	sqlS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING STRING NOT NULL UNIQUE,
		email STRING,
		user_id INTEGER,
		created_at DATETIME
	)`, tableNameSession)

	Db.Exec(sqlS)
}

//uuidの作成・UUIDは唯一無二のIDのこと
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

//passwordをハッシュ値に変換
func Encrypt(plaintext string) (cryptext string) {
	//ハッシュ値に変換
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}