package teams

import (
	"encoding/json"
	"net/http"
)

func CreateTeams(w http.ResponseWriter, r *http.Request) {
	test := "{'create': 'ksajdlks'}"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}
