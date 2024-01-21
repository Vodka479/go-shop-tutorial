package usersUsecases

import (
	"fmt"

	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/Vodka479/go-shop-tutorial/modules/users"
	"github.com/Vodka479/go-shop-tutorial/modules/users/usersRepositories"
	"github.com/Vodka479/go-shop-tutorial/pkg/goauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetUserProfile(userId string) (*users.User, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//Insert user  //ใช้งาน userRepo
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//Insert user  //ใช้งาน userRepo
	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	//Sign Token
	accessToken, err := goauth.NewGoAuth(goauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	//Refresh Token
	refreshToken, err := goauth.NewGoAuth(goauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	// Set Passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil

}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	claims, err := goauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := goauth.NewGoAuth(
		goauth.Access,
		u.cfg.Jwt(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := goauth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
