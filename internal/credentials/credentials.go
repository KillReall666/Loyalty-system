package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string
	PasswordHash string
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

/*
Сет юзер
user := User{
    Username: "example_user",
}
err := user.SetPassword("example_password")
if err != nil {
    panic(err)
}


Проверка пароля
if user.ComparePassword("example_password") {
    fmt.Println("Пароль верен")
} else {
    fmt.Println("Пароль неверен")
}
*/
