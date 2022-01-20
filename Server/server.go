package server

import (
	"encoding/json"
	"fmt"
	"log"
	"main/model"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	port        int
	userservice service.UserService
}

func New(port int, userservice service.UserService) (*Server, error) {
	return &Server{port: port, userservice: userservice}, nil
}

func (s Server) HomePage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(http.StatusOK)
	return
}

func (s Server) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user model.UserDetail
	var response model.Response
	w.Header().Set("Content-Type", "application/json")
	//read body of the request
	err := json.NewDecoder(r.Body).Decode(&user)
	response.Message = ""
	if err != nil {
		log.Println("Error while reading body of login page", err)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Error in reading HTTP Body"
		json.NewEncoder(w).Encode(response)
		return
	}
	if s.userservice.Login(user.UserName, user.Passwd) {
		response.Message = "Login Successful"
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		response.Message = "Login Unsuccessful"
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (s Server) SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response model.Response
	response.Message = ""
	var user model.UserDetail
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error while reading body of signup page", err)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Error in reading HTTP Body"
		json.NewEncoder(w).Encode(response)
		return
	}
	if len(user.UserName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "User Name Missing"
		json.NewEncoder(w).Encode(response)
		return
	} else if len(user.Passwd) == 0 || len(user.FirstName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "First Name Missing"
		json.NewEncoder(w).Encode(response)
		return
	} else if len(user.Passwd) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Password Length is insufficient"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		if s.userservice.Signup(&user) {
			w.WriteHeader(http.StatusOK)
			response.Message = "SignUp Succesfull"
			json.NewEncoder(w).Encode(response)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			response.Message = "SignUp Unsuccesfull"
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (s Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/homepage", s.HomePage)
	r.HandleFunc("/login", s.LoginHandler).Methods("post")
	r.HandleFunc("/signup", s.SignupHandler).Methods("post")
	port := fmt.Sprintf(":%d", s.port)
	log.Fatal(http.ListenAndServe(port, r))
}

func (s Server) Stop() {
	//stop server port
}
