package main

import (
	"log"
	"net/http"
)

const ffmpegPath = "D://ffmpeg/bin/ffmpeg.exe"

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/result.gif", gif)
	if err := http.ListenAndServe(":1789", nil); err != nil {
		log.Fatal("服务启动失败")
	}
}
