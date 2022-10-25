package helpers

import "encoding/json"

type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Success string `json:"success"`
}

func SuccessMsg(msg string) []byte {
	success, err := json.Marshal(
		Success{
			Success: msg,
		})
	if err != nil {
		panic(err)
	}
	return success
}

func ErrorMsg(msg string) []byte {
	error, err := json.Marshal(
		Error{
			Error: msg,
		})
	if err != nil {
		panic(err)
	}
	return error
}
