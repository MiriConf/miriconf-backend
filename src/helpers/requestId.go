package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InvalidID(resource string) []byte {
	error, err := json.Marshal(
		Error{
			Error: "invalid " + resource + " id",
		})
	if err != nil {
		panic(err)
	}
	return error
}

func GetRequestID(resource string, request *http.Request, response http.ResponseWriter) (int, primitive.ObjectID) {
	status := 0
	reqParams := mux.Vars(request)
	teamID, err := primitive.ObjectIDFromHex(reqParams["id"])
	if err != nil {
		EndpointError("invalid "+resource+" id", request)
		error := InvalidID(resource)
		response.Write(error)
		status = 1
		return status, teamID
	}
	return status, teamID
}
