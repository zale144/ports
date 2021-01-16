package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zale144/ports/clientapi/internal/api"
	"github.com/zale144/ports/clientapi/internal/api/handler"
)

func ServeHTTP(port int64, service api.PortService) {
	http.HandleFunc("/port/process", handler.PortJSON(service))
	http.HandleFunc("/port/get", handler.GetPorts(service))
	log.Println("HTTP Listening on localhost:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
