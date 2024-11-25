package models

import (
	"Yakudza/pkg/database"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User - модель пользователей
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Login     string    `gorm:"index:,unique" json:"login,omitempty"`
	Password  string    `gorm:"not null" json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func AllUsers() ([]*User, error) {
	db := database.GetDB()
	var users []*User

	err := db.Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Create() error {
	db := database.GetDB()
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return db.Create(&u).Error
}

func (u *User) FindUserLogin() error {
	db := database.GetDB()

	return db.Where(&User{Login: u.Login}).Find(&u).Error
}

func ComparePassword(login string, password string) (bool, error) {
	u := &User{Login: login}
	if err := u.FindUserLogin(); err != nil {
		return false, err
	}

	if u.ID == 0 {
		return false, errors.New(fmt.Sprintf("Пользователь не найден"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil

}

func (u *User) Delete() error {
	db := database.GetDB()

	return db.Delete(&u).Error
}

func (u *User) DeleteByID() error {
	dbase := database.GetDB()
	result := dbase.Delete(&User{ID: u.ID})

	return result.Error
}

func (u *User) Update() error {
	dbase := database.GetDB()
	result := dbase.Updates(&u)
	return result.Error
}
