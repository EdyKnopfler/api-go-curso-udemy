package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"com.blocopad/blocopad/internal/backend"
	"com.blocopad/blocopad/internal/db"
	"com.blocopad/blocopad/internal/security"

	"github.com/gorilla/mux"
)

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(body)
	w.Write(payload)
}

func ReadNote(w http.ResponseWriter, r *http.Request) {
	if !security.ValidateToken(r) {
		WriteResponse(http.StatusForbidden, map[string]string{"status": "Not authorized"}, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if data, err := backend.GetKey(id); err == nil {
		WriteResponse(http.StatusOK, data, w)
	} else {
		if err.Error() == "Not found" {
			WriteResponse(http.StatusNotFound, "Note not found", w)
		} else {
			log.Println(err.Error())
			WriteResponse(http.StatusInternalServerError, "Error", w)
		}
	}
}

func WriteNote(w http.ResponseWriter, r *http.Request) {
	if !security.ValidateToken(r) {
		WriteResponse(http.StatusForbidden, map[string]string{"status": "Not authorized"}, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var note db.Note

	if err := decoder.Decode(&note); err != nil {
		WriteResponse(
			http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}

	uuidString, err := backend.SaveKey(note.Text, note.OneTime)

	if err != nil {
		WriteResponse(
			http.StatusBadRequest, map[string]string{"error": "invalid request"}, w)
	} else {
		WriteResponse(http.StatusOK, map[string]string{"code": uuidString}, w)
	}
}

func main() {
	serverPort := "8080"
	if port, hasValue := os.LookupEnv("API_PORT"); hasValue {
		serverPort = port
	}

	databaseUrl := "localhost:6379"
	if dbUrl, hasValue := os.LookupEnv("API_DB_URL"); hasValue {
		databaseUrl = dbUrl
	}

	databasePassword := ""
	if dbPassword, hasValue := os.LookupEnv("API_DB_PASSWORD"); hasValue {
		databasePassword = dbPassword
	}

	db.DatabaseUrl = databaseUrl
	db.DatabasePassword = databasePassword

	security.GetKeys()

	router := mux.NewRouter()
	router.HandleFunc("/api/login", security.Login).Methods("POST")
	router.HandleFunc("/api/note/{id}", ReadNote).Methods("GET")
	router.HandleFunc("/api/note", WriteNote).Methods("POST")
	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), router)
	fmt.Println(err)
}
