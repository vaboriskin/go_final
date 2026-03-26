package server

import (
	"fmt"
	"go_final/pkg/api"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	port := getPort()
	api.Init()
	webDir := "./web"

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	fmt.Println("сервер запущен на порту", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ошибка запуска сервера", err)
	}

}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
