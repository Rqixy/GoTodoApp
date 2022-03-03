package utils

import (
	"io"
	"log"
	"os"
)

//logの設定
func LoggingSettings(logFile string) {
	//logfileを開く
	//logの読み書きと作成と追記
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if 	 err != nil {
		log.Fatal(err)
	}
	//logの書き込み先を標準出力とLogfileに指定
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	//logのフォーマットを指定
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
	//logの出力先を指定
	log.SetOutput(multiLogFile)
}