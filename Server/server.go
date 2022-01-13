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
	//read body of the request
	var user dbm.UserDetail
	ret := 0
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error while reading body of login page", err)
		http.Error(w, "Error in reading body", http.StatusBadRequest)
		return
	}

	//check the existance of user
	//return the user id
	user.ID, ret = dbm.UsernameExistanceCheck(user.UserName)
	if user.ID == dbm.IDNone {
		json.NewEncoder(w).Encode(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("http://localhost:8090/signup")
		return
	}
	if ret == dbm.SQLTableConnErr || ret == dbm.SQLTableScanErr {
		json.NewEncoder(w).Encode(http.StatusUnauthorized)
		return
	}
	userlogin, ret := dbm.UserLogin(user.UserName, user.Passwd)
	if ret == dbm.UserLoginSuccess {
		json.NewEncoder(w).Encode(http.StatusOK)
		json.NewEncoder(w).Encode(userlogin.ID)
	} else {
		json.NewEncoder(w).Encode(http.StatusUnauthorized)
	}
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	var user dbm.UserDetail
	ret := 0
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error while reading body of signup page", err)
		http.Error(w, "Error in reading body", http.StatusBadRequest)
		return
	}

	// db = connectDB()
	// defer db.Close()

	if len(user.UserName) == 0 {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	} else if len(user.Passwd) == 0 || len(user.FirstName) == 0 {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	} else if len(user.Passwd) < 8 {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	} else {
		userCheck := dbm.UsernameDupCheck(user.UserName)
		if userCheck == dbm.UserFound {
			json.NewEncoder(w).Encode(http.StatusUnauthorized)
			return
		} else if userCheck == dbm.UserNotFound {
			user.ID, ret = dbm.CreateUser(user.UserName, user.Passwd, user.FirstName, user.LastName)
			if ret != dbm.AllOK {
				json.NewEncoder(w).Encode(http.StatusBadRequest)
				return
			}
			if user.ID == dbm.IDNone {
				json.NewEncoder(w).Encode(http.StatusBadRequest)
				return
			}
		} else {
			json.NewEncoder(w).Encode(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(http.StatusOK)
		user.Passwd = strings.Repeat("*", len(user.Passwd))
		json.NewEncoder(w).Encode(user)
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
