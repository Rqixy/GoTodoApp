package config

import (
	"log"
	"sampleapp/utils"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port 	  string	//ポート番号
	SQLDriver string	//SQLの名前
	DbName    string	//データベースの名前
	LogFile   string	//ログを残すファイル
	Static 	  string	//静的ファイルのある階層
}


//外部から呼べるようにする
var Config ConfigList

//mainが呼ばれる前にConfigListを作成する
func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

func LoadConfig() {
	//iniファイルを読み込む
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	// グローバルに宣言されたConfigに値を代入する
	Config = ConfigList{
		Port: 	   cfg.Section("web").Key("port").MustString("8080"),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		DbName:    cfg.Section("db").Key("name").String(),
		LogFile:   cfg.Section("web").Key("logfile").String(),
		Static:    cfg.Section("web").Key("static").String(),
	}

}