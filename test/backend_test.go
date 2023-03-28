package tests

import (
	"errors"
	"strings"
	"testing"

	"com.blocopad/blocopad/internal/backend"
	"com.blocopad/blocopad/internal/db"
)

func TestGetKeyOk(t *testing.T) {
    // Dado
    db.GetNote = func(key string) (bool, string, error) {
        return false, "OK", nil
    }
    
    // Quando
    data, err := backend.GetKey("Kânia")
    
    // Então
    if err != nil {
        t.Fatal("TestGetKeyOk Não deve devolver erro")
    }
    
    if data != "OK" {
        t.Fatal("TestGetKeyOk Resposta inválida")
    }
}

func TestGetKeyOkMaxSize(t *testing.T) {
    key := strings.Repeat("k", 36)
    
    db.GetNote = func(key string) (bool, string, error) {
        return false, key, nil
    }
    
    data, err := backend.GetKey(key)
    
    if err != nil {
        t.Fatal("TestGetKeyOkMaxSize Não deve devolver erro")
    }
    
    if data != key {
        t.Fatal("TestGetKeyOkMaxSize Resposta inválida")
    }
}

func TestGetErrorSizeZero(t *testing.T) {
    keyString := ""
    
    _, err := backend.GetKey(keyString)
    
    if err == nil {
        t.Fatal("TestGetErrorSizeZero Deve devolver erro com chave de tamanho zero")
    }
}

func TestGetErrorSizeBiggerThan36(t *testing.T) {
    keyString := strings.Repeat("a", 37)
    
    _, err := backend.GetKey(keyString)
    
    if err == nil {
        t.Fatal("TestGetErrorSizeBiggerThan36")
    }
}

func TestGetKeyDbError(t *testing.T) {
    db.GetNote = func(key string) (bool, string, error) {
        return false, "OK", errors.New("Erro qualquer")
    }
    
    _, err := backend.GetKey("mimimi")
    
    if err == nil {
        t.Fatal("TestGetKeyDbError Deveria devolver erro")
    }
}

func TestSaveKeyOK(t *testing.T) {
    db.SaveNote = func(data string, oneTime bool) (string, error) {
        return "123456", nil
    }
    
    uuid, err := backend.SaveKey("bla bla bla", false)
    
    if err != nil {
        t.Fatal("TestSaveKeyOK Não deveria devolver erro")
    }
    
    if uuid != "123456" {
        t.Fatal("TestSaveKeyOK Resposta inválida")
    }
}

func TestSaveKeyDbError(t *testing.T) {
    db.SaveNote = func(data string, oneTime bool) (string, error) {
        return "123456", errors.New("Erro qualquer")
    }
    
    _, err := backend.SaveKey("bla bla bla", false)
    
    if err == nil {
        t.Fatal("TestSaveKeyDbError Deveria devolver erro")
    }
}

func TestSaveInvalidSizeZero(t *testing.T) {
    dataZeroLength := ""
    
    _, err := backend.SaveKey(dataZeroLength, false)
    
    if err == nil {
        t.Fatal("TestSaveInvalidSizeZero Deveria devolver erro")
    }
}

func TestSaveInvalidSizeBiggerThan36(t *testing.T) {
    dataTooBig := strings.Repeat("a", 37)
    
    _, err := backend.SaveKey(dataTooBig, false)
    
    if err == nil {
        t.Fatal("TestSaveInvalidSizeBiggerThan36 Deveria devolver erro")
    }
}

func TestGetKeyDeleteOk(t *testing.T) {
    deleteInvoked := false
    deletedKey := ""
    
    db.GetNote = func(key string) (bool, string, error) {
        return true, "OK", nil
    }
    
    db.DeleteNote = func(key string) error {
        deleteInvoked = true
        deletedKey = key
        return nil
    }
    
    data, err := backend.GetKey("mimimi")
    
    if err != nil {
        t.Fatal("TestGetKeyDeleteOk Não deveria devolver erro")
    }
    
    if data != "OK" {
		t.Fatal("TestGetKeyDeleteOk Dado inválido")
	}
	
	if !deleteInvoked {
		t.Fatal("TestGetKeyDeleteOk Deveria apagar a chave \"onetime\" após consulta")
	}
	
	if deletedKey != "mimimi" {
		t.Fatal("TestGetKeyDeleteOk Deveria ter apagado a chave correta")
	}
}

func TestGetKeyDeleteDbError(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("TestGetKeyDeleteDbError O código não \"priou cânico\"")
        }
    }()
    
    deleteInvoked := false
    deletedKey := ""
    
    db.GetNote = func(key string) (bool, string, error) {
        return true, "OK", nil
    }
    
    db.DeleteNote = func(key string) error {
        deleteInvoked = true
        deletedKey = key
        return errors.New("Erro qualquer")
    }
    
    _, err := backend.GetKey("mimimi")
    
    if err == nil {
        t.Fatal("TestGetKeyDeleteDbError Deveria ter devolvido um erro")
    }
    
    if !deleteInvoked {
		t.Fatal("TestGetKeyDeleteDbError Deveria ter tentado apagar a chave")
	}
	
	if deletedKey != "mimimi" {
		t.Fatal("TestGetKeyDeleteDbError Deveria ter tentado apagar a chave correta")
	}
}

