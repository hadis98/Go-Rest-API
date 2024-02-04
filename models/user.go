package models

import (
	"errors"

	"github.com/hadis98/rest-api/db"
	"github.com/hadis98/rest-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email,password) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//* we don't save the plaintext password
	// * instead, we use the hashed version of it
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (user *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	// in our database creation, we said that the emial column must be unique
	// so an email cannot exists more than once
	row := db.DB.QueryRow(query, user.Email)
	var retrievedPassword string
	err := row.Scan(&user.ID, &retrievedPassword) //the password that we found in our database for the entered email
	// Scan will return an error if it can't find any matched query
	if err != nil {
		// !we did not find a user for that email
		return errors.New("credentials invalid")
	}
	// compare the password to the password that user has submitted with request
	passwordIsValid := utils.CheckPasswordHash(user.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("credentials invalid")
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
