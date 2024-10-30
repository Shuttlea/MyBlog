package main

import (
	"log"
	"myblog/db"
	"myblog/handlers"
	"net/http"
	"time"

	"sync/atomic"
)

func worker(limit *int32){
	time.Sleep(time.Second)
	atomic.AddInt32(limit,1)
	// FOR DEMONSTRATE COMMENT PREV AND UNCOMMENT NEXT LINE
	// log.Println(atomic.AddInt32(limit,1))
}

func limiter(limit *int32,next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		go worker(limit)
		if atomic.AddInt32(limit,-1) < 0 {
			http.Error(w,http.StatusText(http.StatusTooManyRequests),http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w,r)
	})
}

func main() {
	var limit int32 = 100
	dataBase := db.ConnectDB()
	defer db.CloseDB(dataBase)
	h := handlers.DBHandler{DB: dataBase}
	mux := http.NewServeMux()
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))))
	mux.Handle("/Exercise00/", http.StripPrefix("/Exercise00/", http.FileServer(http.Dir("../../Exercise00"))))
	mux.HandleFunc("/styles/style.css", handlers.HandlerFavicon)
	mux.HandleFunc("/", h.HandlerMain)
	mux.HandleFunc("/article", h.HandlerArticle)
	mux.HandleFunc("/create", h.HandlerCreate)
	mux.HandleFunc("/authpage", handlers.HandlerAuth)
	log.Println("server is listening")
	log.Fatal(http.ListenAndServe(":8888", limiter(&limit,mux)))
}
