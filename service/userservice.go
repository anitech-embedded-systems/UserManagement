package service

import (
	"log"
	data "main/Data"
	"main/model"
)

type UserService struct {
	userrepo data.UserRepo
}

func New(userrepo data.UserRepo) (*UserService, error) {
	return &UserService{userrepo: userrepo}, nil
}

func (s *UserService) Login(username, password string) bool {
	user, err := s.userrepo.FindByUsername(username)
	if err != nil {
		return false
	}
	if user.Passwd == password {
		return true
	}
	return false
}

func (s *UserService) Signup(user *model.UserDetail) bool {
	// user, err := s.userrepo.FindByUsername(user.UserName)
	if !(s.userrepo.FindByUsername_anycase(user.UserName)) {
		id := s.userrepo.ExtractNewUserID(user.UserName)
		if id == model.IDNone {
			return false
		}
		user.ID = id
		if s.userrepo.CreateUser(user) {
			return true
		}
		return false
	}
	log.Println("username already exist")
	return false
}
