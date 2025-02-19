package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ErrorResponse struct {
	Error string
}

// Функия ответа на запрос с ошибкой.
func respondWithError(text string, w http.ResponseWriter) {
	res := ErrorResponse{text}
	toWrite, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Println("Error marshalling data")
		return
	}
	_, err = w.Write(toWrite)
	if err != nil {
		log.Println("Error writing the data")
		return
	}
}

/*func authorized(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(403)
		respondWithError("Bad method", w)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Message  string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" /*|| user.Message == "" * {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

}
*/
// r.Method - медот HTTP(GET, POST, PUT) и др.
// r.URL.Path - часть пути URL-адреса
// r.RemoteAddr - IP-адрес клиента, сделавшего запрос
func authHandler(w http.ResponseWriter, r *http.Request) {
	// /login
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method != "POST" {
		w.WriteHeader(403)
		respondWithError("Bad method", w)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Message  string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" /*|| user.Message == "" */ {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//Здесь должна быть проверка пользователя в БД
	//
	if user.Username != ExpectedUsername && user.Password != ExpectedPassword {
		http.Error(w, "Unauthorization", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.Username)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(token))

	/*body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error read query", http.StatusBadRequest)
		return
	}
	message := string(body)
	err = addEchotoDB(message)
	if err != nil {
		http.Error(w, "Error write in DB", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message saved in DB: %s\n", message)
	return

	*/

	//stuff, _ := io.ReadAll(r.Body)
	//w.Write(stuff)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// /echo
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method == http.MethodGet {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		log.Println(token)
		if !validateToken(token) {
			http.Error(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}

		message, err := getEchomessageList()
		if err != nil {
			http.Error(w, "Error read in DB", http.StatusInternalServerError)
			return
		}

		var response string
		for _, msg := range message {
			response += fmt.Sprintf("ID: %d, Message: %s \n", msg.ID, msg.Message)
		}

		w.Write([]byte(response))
		return
	}

}
func addEchoHandler(w http.ResponseWriter, r *http.Request) {
	// /add
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method != "POST" {
		w.WriteHeader(403)
		respondWithError("Bad method", w)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Message  string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" || user.Message == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	log.Println(token)
	if !validateToken(token) {
		http.Error(w, "Unauthorized!", http.StatusUnauthorized)
		return
	}

	message := user.Message
	err = addEchotoDB(message)
	if err != nil {
		http.Error(w, "Error write in DB", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message saved in DB: %s\n", message)
}
