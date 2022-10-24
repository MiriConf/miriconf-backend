package helpers

import (
	"log"
	"net/http"
	"strings"
)

func SuccessLog(req *http.Request) {
	sanitizedURL := strings.Replace(req.URL.String(), "\n", "", -1)
	sanitizedURL = strings.Replace(sanitizedURL, "\r", "", -1)
	sanitizedAgent := strings.Replace(string(req.UserAgent()), "\n", "", -1)
	sanitizedAgent = strings.Replace(sanitizedAgent, "\r", "", -1)

	log.Printf("successful %s on %s from %s\n", req.Method, sanitizedURL, sanitizedAgent)
}

func EndpointError(msg string, req *http.Request) {
	sanitizedURL := strings.Replace(req.URL.String(), "\n", "", -1)
	sanitizedURL = strings.Replace(sanitizedURL, "\r", "", -1)
	sanitizedAgent := strings.Replace(string(req.UserAgent()), "\n", "", -1)
	sanitizedAgent = strings.Replace(sanitizedAgent, "\r", "", -1)

	log.Printf("error: %s for %s request on %s from %s", msg, req.Method, sanitizedURL, sanitizedAgent)
}
