package teams

import (
	"encoding/json"
	"net/http"
)

func EditTeams(w http.ResponseWriter, r *http.Request) {
	test := "{'edit': 'ksajdlks'}"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}
