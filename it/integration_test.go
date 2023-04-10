package it

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-redis/redis/v8"
)

type SavedNote struct {
	Code string `json:"code"`
}

func TestSaveOK(t *testing.T) {
	postUrl := "http://localhost:8080/api/note"
	jsonData := []byte(`{
		"data": "eu deveria salvar isto?",
		"onetime": false
	}`)

	request, _ := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatal("TestSaveOK não deveria devolver erro")
	}

	var savedNote SavedNote
	err = json.NewDecoder(response.Body).Decode(&savedNote)

	if err != nil {
		t.Fatal("TestSaveOK deveria devolver um JSON válido")
	}

	if len(savedNote.Code) == 0 {
		t.Fatal("TestSaveOK deveria devolver um UUID válido")
	}

	request2, _ := http.NewRequest("GET", postUrl+"/"+savedNote.Code, nil)
	response2, err2 := client.Do(request2)

	if err2 != nil {
		panic(err2)
	}

	defer response2.Body.Close()
	binary, readError := io.ReadAll(response2.Body)

	if readError != nil {
		panic(readError)
	}

	if string(binary) != "\"eu deveria salvar isto?\"" {
		t.Fatal("TestSaveOK não recuperou a nota" + string(binary))
	}

	rDB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	expiration := rDB.TTL(ctx, savedNote.Code)

	if expiration.Val().Hours() < 23 {
		t.Fatal("TestSaveOK deveria setar a expiração para 24 horas")
	}
}
