package model

type UserDetail struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Passwd    string `json:"passwd"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type Response struct {
	Message  string     `json:"message"`
	UserInfo UserDetail `json: userinfo`
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

// type Errorstring struct {
// 	s string
// }

// func (e *Errorstring) Error() string {
// 	return e.s
// }
