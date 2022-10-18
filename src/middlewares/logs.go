package middlewares

import (
	"log"
	"net/http"
)

func CallLog(req *http.Request) {
	log.Printf("successful %s on %s from %s\n", req.Method, req.RequestURI, req.UserAgent())
}