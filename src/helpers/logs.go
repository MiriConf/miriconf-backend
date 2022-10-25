package helpers

import (
	"log"
	"net/http"
)

func SuccessLog(req *http.Request) {
	log.Printf("successful %s on %s from %s\n", req.Method, req.RequestURI, req.UserAgent())
}

func EndpointError(msg string, req *http.Request) {
	log.Printf("error: %s for %s request on %s from %s", msg, req.Method, req.RequestURI, req.UserAgent())
}
