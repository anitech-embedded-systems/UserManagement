package dataimpl

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	data "main/Data"
	"main/model"
	"strings"
)

type UserData struct {
	data.UserRepo
	client *sql.DB
}

func New(client *sql.DB) (*UserData, error) {
	return &UserData{client: client}, nil
}

func (d UserData) FindByUsername(username string) (model.UserDetail, error) {
	// db query
	// select *
	var user model.UserDetail
	query := fmt.Sprintf("SELECT id, passwd, first_name, last_name FROM UserBase WHERE username=\"%s\"", username)
	res := d.client.QueryRow(query)
	errScan := res.Scan(&user.ID, &user.Passwd, &user.FirstName, &user.LastName)
	switch errScan {
	case sql.ErrNoRows:
		log.Println("No Row scanned")
		return model.UserDetail{}, errors.New("SQLTableScanErr")
	case nil:
		return user, nil
	default:
		log.Println("out of scope")
		return model.UserDetail{}, errors.New("SQLTableScanErr")
	}
}

func (d UserData) FindByUsername_anycase(username string) bool {
	query := fmt.Sprintf("SELECT username FROM UserBase")
	result, err := d.client.Query(query)
	if err != nil {
		log.Println("Error while query the DB")
		return false
	}
	for result.Next() {
		var userDup string
		err := result.Scan(&userDup)
		if err != nil {
			log.Println("Error in scanning username from DB table")
			log.Fatal(err)
		}
		if strings.ToLower(userDup) == strings.ToLower(username) {
			return true
		}
	}
	return false
}

func (d UserData) ExtractNewUserID(username string) int {
	idMax := 0
	id := 0
	res, err := d.client.Query("SELECT id FROM UserBase")
	if err != nil {
		log.Println("error while interacting with table")
		return model.IDNone
	}
	for res.Next() {
		_ = res.Scan(&id)
		if id > idMax {
			idMax = id
		}
	}
	idMax = idMax + 1
	return idMax
}

func (d UserData) CreateUser(user *model.UserDetail) bool {
	if user.ID == model.IDNone {
		return false
	}
	query := fmt.Sprintf("INSERT INTO UserBase(id, username, passwd, first_name, last_name) VALUES (\"%d\", \"%s\", \"%s\", \"%s\", \"%s\")", user.ID, user.UserName, user.Passwd, user.FirstName, user.LastName)
	_, err := d.client.Query(query)
	if err != nil {
		log.Println("error while creating user in the table", err)
		return false
	}
	log.Printf("New Assign ID: %d", user.ID)
	return true
}
