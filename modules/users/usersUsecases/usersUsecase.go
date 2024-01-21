package usersUsecases

import (
	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/Vodka479/go-shop-tutorial/modules/users"
	"github.com/Vodka479/go-shop-tutorial/modules/users/usersRepositories"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg               config.IConfig
	usersRepositories usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepositories usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:               cfg,
		usersRepositories: usersRepositories,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//Insert user  //ใช้งาน userRepo
	result, err := u.usersRepositories.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}
