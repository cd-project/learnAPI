package service

import (
	"errors"
	"log"
	"time"
	"todo/infrastructure"
	"todo/model"
	"todo/repository"

	middleware "todo/middlewares"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAll() ([]model.UserResponse, error)
	GetByID(id int) (*model.UserResponse, error)
	CreateUser(newUser *model.User) (*model.UserResponse, error)
	ChangePassword(id int, newPwd, oldPwd string) (*model.UserResponse, error)
	ChangeRole(id int, newRole string) (*model.UserResponse, error)
	DeleteUser(id int) error
	ResetPassword(id int) (*model.UserResponse, error)

	CheckCredential(id int, password string) (*model.User, error)
	LoginRequest(username, password string) (*model.User, string, string, error)
	LoginWithToken(token string) (*model.User, string, bool, error)
}

type userService struct {
	userRepository model.UserRepository
}

func (s *userService) GetAll() ([]model.UserResponse, error) {
	return s.userRepository.GetAll()
}

func (s *userService) GetByID(id int) (*model.UserResponse, error) {
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	userResponse := model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	return &userResponse, nil
}

func (s *userService) CreateUser(newUser *model.User) (*model.UserResponse, error) {
	newUser.Password = hashAndSalt(newUser.Password)
	user, err := s.userRepository.Insert(newUser)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil

}

func (s *userService) ChangePassword(id int, newPwd, oldPwd string) (*model.UserResponse, error) {
	// validate username/password
	user, err := s.CheckCredential(id, oldPwd)
	if err != nil {
		return nil, err
	}
	// validated, update new password
	hashedPwd := hashAndSalt(newPwd)
	err = s.userRepository.ChangePassword(id, hashedPwd)
	if err != nil {
		return nil, err
	}

	// User reponse
	userResponse := model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	return &userResponse, nil
}

func (s *userService) ChangeRole(id int, newRole string) (*model.UserResponse, error) {
	user, err := s.userRepository.ChangeRole(id, newRole)
	if err != nil {
		return nil, err
	}
	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepository.DeleteUser(id)
}

func (s *userService) ResetPassword(id int) (*model.UserResponse, error) {
	newPwd := model.DefaultPassword
	hashedPwd := hashAndSalt(newPwd)
	if err := s.userRepository.ChangePassword(id, hashedPwd); err != nil {
		return nil, err
	}
	return s.GetByID(id)

}
func (s *userService) LoginRequest(username, password string) (*model.User, string, string, error) {
	// validate username and password
	user, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return nil, "", "", err
	}
	err = s.checkPassword(user, password)
	if err != nil {
		return nil, "", "", err
	}
	// err getting JWT
	// Get JWT
	tokenString, refreshToken, err := middleware.GetTokenString(user)
	if err != nil {
		infrastructure.ErrLog.Printf("Problem with LoginRequest by Authentication: %v/n", err)
		return nil, "", "", err
	}
	//
	return user, tokenString, refreshToken, nil

}

func (s *userService) LoginWithToken(token string) (*model.User, string, bool, error) {
	user, err := middleware.GetClaimsData(token)
	if err != nil {
		infrastructure.ErrLog.Printf("Problem with LoginWithToken at Getting claims data: %v/n", err)
		return nil, "Token is invalid", false, nil
	}

	timeLeft := user.VerifyExpiresAt(time.Now().UnixNano()/infrastructure.NANO_TO_SECOND, true)
	if !timeLeft {
		infrastructure.ErrLog.Println("Expire exceeded; session expired!")
		return nil, "Token has exceeded expire time!", false, nil
	}

	if ok, err := s.userRepository.LoginTokenRequest(user); err != nil {
		infrastructure.ErrLog.Printf("Problem with Login with Token, Login Request at repo level~ : %v/n", err)
	} else {
		if ok == false {
			return nil, "Token invalid", false, nil
		}
	}

	newTokenString, _, err := middleware.GetTokenString(user)
	if err != nil {
		infrastructure.ErrLog.Printf("Problem with Login Request, by authentication: %v,n", err)
		return nil, "", false, err
	}

	return user, newTokenString, true, nil
}

func (s *userService) CheckCredential(id int, password string) (*model.User, error) {
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !comparePassword(user.Password, password) {
		return nil, errors.New("incorrect password from service/check credential")
	}
	return user, nil
}

func (s *userService) checkPassword(user *model.User, password string) error {
	// check password: hashed one.
	// compare: hashed, plain
	if !comparePassword(user.Password, password) {
		return errors.New("incorrect password from service/checkPassword")
	}

	return nil
}

func hashAndSalt(password string) string {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println(err.Error() + "/service/hashAndSalt")
	}
	log.Println("hashedPwd is:", string(hashedPwd))
	return string(hashedPwd)
}

func comparePassword(hashedPwd string, plainPwd string) bool {
	// cmphashandpass : true ~ nil.
	log.Println("Login session:")
	log.Println("hashpwd: ", hashedPwd)
	log.Println("plainpwd: ", plainPwd)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}
	return true
}

func NewUserService() UserService {
	userRepo := repository.NewUserRepository()
	return &userService{
		userRepository: userRepo,
	}
}
