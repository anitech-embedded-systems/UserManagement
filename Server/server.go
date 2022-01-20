package server

import (
	"encoding/json"
	"log"
	dbm "main/Data"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(http.StatusOK)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response dbm.Response
	response.Message = ""
	//read body of the request
	var user dbm.UserDetail
	ret := 0
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error while reading body of login page", err)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Error in reading HTTP Body"
		json.NewEncoder(w).Encode(response)
		return
	}

	//check the existance of user
	//return the user id
	user.ID, ret = dbm.UsernameExistanceCheck(user.UserName)
	if user.ID == dbm.IDNone {
		w.WriteHeader(http.StatusUnauthorized)
		response.Message = "User Not Found, Please Signup First, http://localhost:8090/signup"
		json.NewEncoder(w).Encode(response)
		return
	}
	if ret == dbm.SQLTableConnErr || ret == dbm.SQLTableScanErr {
		w.WriteHeader(http.StatusUnauthorized)
		response.Message = "Error in DataBase"
		json.NewEncoder(w).Encode(response)
		return
	}
	userlogin, ret := dbm.UserLogin(user.UserName, user.Passwd)
	if ret == dbm.UserLoginSuccess {
		w.WriteHeader(http.StatusOK)
		response.Message = "Login Successful"
		userlogin.Passwd = ""
		response.UserInfo = userlogin
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		response.Message = "User Name or Password is wrong, please try gain!"
		json.NewEncoder(w).Encode(response)
	}
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response dbm.Response
	response.Message = ""

	var user dbm.UserDetail
	ret := 0
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
		userCheck := dbm.UsernameDupCheck(user.UserName)
		if userCheck == dbm.UserFound {
			w.WriteHeader(http.StatusUnauthorized)
			response.Message = "Username already exists"
			json.NewEncoder(w).Encode(response)
			return
		} else if userCheck == dbm.UserNotFound {
			user.ID, ret = dbm.CreateUser(user.UserName, user.Passwd, user.FirstName, user.LastName)
			if ret != dbm.AllOK {
				w.WriteHeader(http.StatusBadRequest)
				response.Message = "Something wrong with the database"
				json.NewEncoder(w).Encode(response)
				return
			}
			if user.ID == dbm.IDNone {
				w.WriteHeader(http.StatusBadRequest)
				response.Message = "Something wrong with the database"
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response.Message = "Something wrong with the database"
			json.NewEncoder(w).Encode(response)
		}
		w.WriteHeader(http.StatusOK)
		user.Passwd = strings.Repeat("*", len(user.Passwd))
		response.Message = "Signup Successful"
		response.UserInfo = user
		json.NewEncoder(w).Encode(response)
		return
	}

}
func HandleRequest() {
	r := mux.NewRouter()
	r.HandleFunc("/homepage", HomePage)
	r.HandleFunc("/login", LoginPage).Methods("post")
	r.HandleFunc("/signup", SignupPage).Methods("post")
	log.Fatal(http.ListenAndServe(":8090", r))
}
