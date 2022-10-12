package teams

import (
	"encoding/json"
	"net/http"
)

func DeleteTeams(w http.ResponseWriter, r *http.Request) {
	test := "{'delete': 'ksajdlks'}"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}
