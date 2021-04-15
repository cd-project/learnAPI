package repository

import (
	"time"
	"todo/infrastructure"
	"todo/model"
)

type userRepository struct {
}

func (r *userRepository) GetAll() ([]model.UserResponse, error) {
	db := infrastructure.GetDB()

	var responseList []model.UserResponse
	if err := db.Table("users").Select("id, username, role").Order("id").Scan(&responseList).Error; err != nil {
		return nil, err
	}

	return responseList, nil

}

func (r *userRepository) GetByID(id int) (*model.User, error) {
	db := infrastructure.GetDB()

	var user model.User
	// first record ordered by PK
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	db := infrastructure.GetDB()

	var user model.User
	// get first record ordered by PK
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Insert(newUser *model.User) (*model.User, error) {
	db := infrastructure.GetDB()

	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *userRepository) ChangePassword(id int, newPwd string) error {
	db := infrastructure.GetDB()

	if err := db.Model(&model.User{ID: id}).Update("password", newPwd).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) ChangeRole(id int, newRole string) (*model.User, error) {
	db := infrastructure.GetDB()

	if err := db.Model(&model.User{ID: id}).Update("role", newRole).Error; err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *userRepository) DeleteUser(id int) error {
	db := infrastructure.GetDB()

	if err := db.Table("users").Delete("id", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) LoginTokenRequest(user *model.User) (bool, error) {
	db := infrastructure.GetDB()

	var userInfo model.User
	if err := db.Where(&model.User{
		Username: user.Username,
		Password: user.Password,
	}).First(&userInfo).Error; err != nil {
		infrastructure.ErrLog.Println(err)
		return false, nil
	}

	user.ExpiresAt = time.Now().Local().Add(time.Hour*time.Duration(infrastructure.Extend_Hour)).UnixNano() / infrastructure.NANO_TO_SECOND
	return true, nil
}
func NewUserRepository() model.UserRepository {
	return &userRepository{}
}
