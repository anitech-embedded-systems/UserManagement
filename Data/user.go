package data

import "main/model"

type UserRepo interface {
	FindByUsername(username string) (model.UserDetail, error)
	FindByUsername_anycase(username string) bool
	CreateUser(*model.UserDetail) bool
	ExtractNewUserID(username string) int
}
