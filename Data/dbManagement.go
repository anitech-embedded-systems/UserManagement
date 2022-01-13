package dbm

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type UserDetail struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Passwd    string `json:"passwd"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

const (
	SQLTableConnErr = iota
	SQLTableScanErr
	UserFound
	UserNotFound
	UserDetailWrong
	UserDetailCorrect
	UserDetailMissing
	PasswdWrong
	PasswdSyntaxWrong
	UserLoginSuccess
	AllOK
)
const (
	IDNone = 0
)

type USERDB struct {
	MyDB *sql.DB
}

var Mydb USERDB

func ConnectDB() (db *sql.DB) {
	var err error
	Mydb.MyDB, err = sql.Open("mysql", "root:Aaaa@(127.0.0.1:3306)/user_login")
	if err != nil {
		panic(err.Error())
	}
	Mydb.MyDB.Ping()
	return Mydb.MyDB
}

//Checking the duplicate username
func UsernameDupCheck(username string) int {
	result, err := Mydb.MyDB.Query("SELECT username FROM UserBase")
	if err != nil {
		log.Println("Couldn't connect to DB", err)
		return SQLTableConnErr
	}
	for result.Next() {
		var userDup string
		err := result.Scan(&userDup)
		if err != nil {
			log.Println("Error in scanning username from DB table")
			log.Fatal(err)
		}
		if strings.ToLower(userDup) == strings.ToLower(username) {
			return UserFound
		}
	}
	return UserNotFound
}

//User SignUp
func CreateUser(username string, passwd string, firstname string, lastname string) (idRet int, ret int) {
	//extract assignable user id from table
	//do not reassign deleted user's id
	//extract biggest number out of all user id
	id := ExtractNewUserID()
	if id <= 0 {
		return IDNone, SQLTableConnErr
	}
	query := fmt.Sprintf("INSERT INTO UserBase(id, username, passwd, first_name, last_name) VALUES (\"%d\", \"%s\", \"%s\", \"%s\", \"%s\")", id, username, passwd, firstname, lastname)
	_, err := Mydb.MyDB.Query(query)
	if err != nil {
		log.Println("error while creating user in the table", err)
		return id, SQLTableConnErr
	}
	log.Printf("New Assign ID: %d", id)
	return id, AllOK
}

//extract user id
func ExtractNewUserID() int {
	idMax := 0
	id := 0
	res, err := Mydb.MyDB.Query("SELECT id FROM UserBase")
	if err != nil {
		log.Println("error while interacting with table")
		return SQLTableConnErr
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

//login user
//return: userdetail, err
//{all feilds},{-1, 0, 200}
func UserLogin(username string, passwd string) (userdetail UserDetail, ret int) {
	//check if user exists
	var user UserDetail
	err := 0
	user.ID, err = UsernameExistanceCheck(username)
	if err == 0 {
		log.Printf("user doesn't exist, please signup first")
		return user, UserNotFound
	}

	query := fmt.Sprintf("SELECT id, passwd, first_name, last_name FROM UserBase WHERE username=\"%s\"", username)
	res := Mydb.MyDB.QueryRow(query)
	errScan := res.Scan(&user.ID, &user.Passwd, &user.FirstName, &user.LastName)

	switch errScan {
	case sql.ErrNoRows:
		log.Println("No Row scanned")
		return user, SQLTableScanErr
	case nil:
		if user.Passwd == passwd {
			log.Println("LOGIN SUCCESSFUL!")
			return user, UserLoginSuccess
		}
	default:
		log.Println("out of scope")
		return user, SQLTableScanErr
	}
	return user, SQLTableScanErr
}

//user existance check function
//Case sensitive input and comparison
//return: userid, ret
//{userid, 0}, {-1, 0, 1}
//{}, {error in sql, not in the table, found in the table}
func UsernameExistanceCheck(username string) (userid int, ret int) {
	var user UserDetail
	query := fmt.Sprintf("SELECT id FROM UserBase WHERE username=\"%s\"", username)
	res := Mydb.MyDB.QueryRow(query)
	user.ID = 0
	errScan := res.Scan(&user.ID)
	switch errScan {
	case sql.ErrNoRows:
		log.Println("User not found")
		return IDNone, UserNotFound
	case nil:
		if user.ID != IDNone {
			log.Println("username found in the table")
			return user.ID, UserFound
		}
		return IDNone, UserNotFound
	default:
		log.Println("out of scope")
		return IDNone, SQLTableScanErr
	}
}
