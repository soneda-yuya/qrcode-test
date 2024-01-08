package main

import (
	"awesomeProject/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	// リクエストログを出力するミドルウェアを追加
	r.Use(loggingMiddleware)

	assetsHandler := http.FileServer(http.Dir("./assets"))
	// /assets/ パスにハンドラを割り当て
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))

	handler := handlers.NewHandler()
	r.HandleFunc("/tickets", handler.Ticket.Create).Methods("POST")
	r.HandleFunc("/tickets/{ticketToken}", handler.Ticket.Get).Methods("GET")
	r.HandleFunc("/tickets/{ticketToken}/confirm", handler.Ticket.Confirm).Methods("GET")
	r.HandleFunc("/tickets/{ticketToken}/approve", handler.Ticket.Confirm).Methods("GET")

	// ルーターを使用してサーバーを作成
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("start server using port 8080")
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエスト情報をログに出力
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
