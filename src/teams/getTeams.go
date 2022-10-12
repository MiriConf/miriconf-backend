package teams

import (
	"encoding/json"
	"net/http"
)

func GetTeams(w http.ResponseWriter, r *http.Request) {
	test := "{'get': 'ksajdlks'}"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}
