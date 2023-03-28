package backend

import (
	"errors"

	"com.blocopad/blocopad/internal/db"
)

const tam_32k = 32*1024

var GetKey = func(key string) (string, error) {
	if len(key) == 0 || len(key) > 36 {
		return "", errors.New("Tam. máx. da chave: 36")
	}

	oneTime, data, err := db.GetNote(key)

	if err != nil {
		return "", errors.New(err.Error())
	}

	if oneTime {
		if err := db.DeleteNote(key); err != nil {
			panic("Não foi possível deletar a nota de leitura única!")
		}
	}

	return data, nil
}

var SaveKey = func(data string, onetime bool) (string, error) {
	byteSize := len([]rune(data))

	if byteSize == 0 || byteSize > tam_32k {
		return "", errors.New("Tamanho de nota inválido")
	}

	uuidCode, err := db.SaveNote(data, onetime)

	if err != nil {
		return "", errors.New(err.Error())
	}

	return uuidCode, nil
}
