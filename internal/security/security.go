package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Username string
	Password string
}

var (
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
)

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(body)
	w.Write(payload)
}

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
	key, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)

	if err != nil {
		panic(err)
	}

	// dot + parenthesis = type assertion :)
	PrivateKey = key.(*rsa.PrivateKey)

	pubKey, err := ioutil.ReadFile(publicKeyPath)

	if err != nil {
		panic("Erro lendo arquivo da chave pública")
	}

	pubBlock, _ := pem.Decode(pubKey)
	pkey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)

	if err != nil {
		panic(err)
	}

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
	tokenCookie, err := r.Cookie("Token")

	if err != nil {
		log.Println("Erro ao ler o cookie")
		return false
	}

	token, errToken := jwt.Parse(tokenCookie.Value, func(jwtToken *jwt.Token) (interface{}, error) {
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

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var user User

	if err := decoder.Decode(&user); err != nil {
		WriteResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}

	// Faz de conta que tem uma chamada ao Keycloak ou uma autenticação via BD aqui...
	tokenString := ""

	if user.Username == "kânia" && user.Password == "búco" {
		tokenString = CreateToken(user.Username)
		cookie := http.Cookie{
			Name:     "Token",
			Value:    tokenString,
			MaxAge:   1800,
			HttpOnly: true, // Importantíssimo para evitar a leitura por JavaScript :)
		}
		http.SetCookie(w, &cookie)
		WriteResponse(http.StatusOK, map[string]string{"status": "autenticado"}, w)
	} else {
		WriteResponse(http.StatusUnauthorized, map[string]string{"error": "Credenciais inválidas"}, w)
	}
}
