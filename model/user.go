package model

import (
	"fmt"
	"fruit_shop/pkg/auth"
	"fruit_shop/pkg/constattr"

	validator "gopkg.in/go-playground/validator.v9"
)

// UserModel represents a registered user
type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
}

func (u *UserModel) TableName() string {
	return "users"
}

// Create creates a new user account
func (u *UserModel) Create() error {
	return DB.Self.Create(&u).Error
}

// DeleteUser delete the user by user id
func DeleteUser(id uint64) error {
	user := UserModel{}
	user.BaseModel.Id = id
	return DB.Self.Delete(&user).Error
}

// Update updates the user information
func (u *UserModel) Update() error {
	return DB.Self.Save(u).Error
}

// GetUser gets an user by the username
func GetUser(username string) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("username = ?", username).First(&u)
	return u, d.Error
}

// ListUser list all users
func ListUser(username string, offset, limit int) ([]*UserModel, uint64, error) {
	if limit == 0 {
		limit = constattr.DefaultLimit
	}

	users := make([]*UserModel, 0)
	var count uint64

	query := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Self.Model(&UserModel{}).Where(query).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.Self.Where(query).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

// Compare compare with the plain text password.
func (u *UserModel) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// Encrypt user password
func (u *UserModel) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

// Validate the fields
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
