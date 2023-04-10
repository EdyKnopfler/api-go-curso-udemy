package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
)

func GetKeys() {
	var privateKeyPath string
	var publicKeyPath string

	if privKeyFile, hasValue := os.LookupEnv("API_PRIVATE_KEY"); hasValue {
		privateKeyPath = privKeyFile
	} else {
		panic("Favor setar a variável API_PRIVATE_KEY")
	}

	if pubKeyFile, hasValue := os.LookupEnv("API_PUBLIC_KEY"); hasValue {
		publicKeyPath = pubKeyFile
	} else {
		panic("Favor setar a variável API_PUBLIC_KEY")
	}

	privKey, err := ioutil.ReadFile(privateKeyPath)

	if err != nil {
		panic("Erro lendo arquivo da chave privada")
	}

	privBlock, _ := pem.Decode(privKey)
	key, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)

	if err != nil {
		panic(err)
	}

	PrivateKey = key

	pubKey, err := ioutil.ReadFile(publicKeyPath)

	if err != nil {
		panic("Erro lendo arquivo da chave pública")
	}

	pubBlock, _ := pem.Decode(pubKey)
	pkey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)

	if err != nil {
		panic(err)
	}

	// dot + parenthesis = type assertion :)
	rsaKey, ok := pkey.(*rsa.PublicKey)

	if !ok {
		panic("Tipo da chave pública incorreto")
	}

	PublicKey = rsaKey
}

func CreateToken(username string) string {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	claims["authorized"] = true
	claims["user"] = username
	tokenString, err := token.SignedString(PrivateKey)

	if err != nil {
		fmt.Println(err)
		panic(err.Error)
	}

	return tokenString
}

func ValidateToken(r *http.Request) bool {
	if r.Header["Authorization"] == nil {
		return false
	}

	tokenString := strings.Replace(r.Header["Authorization"][0], "Bearer ", "", 1)
	token, errToken := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		return PublicKey, nil
	})

	if errToken != nil {
		log.Println("Token error: ", errToken)
		return false
	}

	if token == nil {
		log.Println("Token inválido")
		return false
	}

	if !token.Valid {
		log.Println("Token marcado como inválido")
		return false
	}

	return true
}
