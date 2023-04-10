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

type AccessToken struct {
	AccessToken string
}

func TestSaveOK(t *testing.T) {
	postUrl := "http://localhost:8080/api/note"
	loginUrl := "http://localhost:8080/api/login"
	jsonData := []byte(`{
		"data": "eu deveria salvar isto?",
		"onetime": false
	}`)
	credentialsData := []byte(`{
		"username": "kânia",
		"password": "búco"
	}`)
	client := &http.Client{}

	// POST da nota sem autenticação
	postRequest, _ := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	postRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postResponse, err := client.Do(postRequest)

	if err != nil {
		panic(err)
	}

	defer postResponse.Body.Close()

	if postResponse.StatusCode != http.StatusForbidden {
		t.Fatal("TestSaveOK deveria devolver 403 Forbidden")
	}

	// GET da nota sem autenticação
	getRequest, _ := http.NewRequest("GET", postUrl+"/abc", nil)
	getResponse, err := client.Do(getRequest)

	if err != nil {
		panic(err)
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusForbidden {
		t.Fatal("TestSaveOK deveria devolver 403 Forbidden")
	}

	// Autenticação :)
	loginRequest, _ := http.NewRequest("POST", loginUrl, bytes.NewBuffer(credentialsData))
	loginRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")
	loginResponse, err := client.Do(loginRequest)

	if err != nil {
		panic(err)
	}

	defer loginResponse.Body.Close()
	loginDecoder := json.NewDecoder(loginResponse.Body)
	var token AccessToken

	if err := loginDecoder.Decode(&token); err != nil {
		panic(err)
	}

	// POST autenticado
	postRequest.Header.Set("Authorization", "Bearer "+token.AccessToken)
	postResponse, err = client.Do(postRequest)

	if err != nil {
		panic(err)
	}

	if postResponse.StatusCode != http.StatusOK {
		t.Fatal("TestSaveOK não deveria devolver erro")
	}

	var savedNote SavedNote
	err = json.NewDecoder(postResponse.Body).Decode(&savedNote)

	if err != nil {
		t.Fatal("TestSaveOK deveria devolver um JSON válido")
	}

	if len(savedNote.Code) == 0 {
		t.Fatal("TestSaveOK deveria devolver um UUID válido")
	}

	// GET autenticado
	getRequest, _ = http.NewRequest("GET", postUrl+"/"+savedNote.Code, nil)
	getRequest.Header.Set("Authorization", "Bearer "+token.AccessToken)
	getResponse, err = client.Do(getRequest)

	if err != nil {
		panic(err)
	}

	defer getResponse.Body.Close()
	binary, readError := io.ReadAll(getResponse.Body)

	if readError != nil {
		panic(readError)
	}

	if string(binary) != "\"eu deveria salvar isto?\"" {
		t.Fatal("TestSaveOK não recuperou a nota" + string(binary))
	}

	// Expiração da nota
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
