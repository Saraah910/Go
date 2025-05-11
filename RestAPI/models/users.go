package models

import (
	"errors"

	"example.com/APIs/DB"
	"example.com/APIs/Utils"
)

type User struct {
	Id       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users(email, password) VALUES(?,?)`
	stmt, err := DB.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	hashedPassword, err := Utils.ConvertToHash(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.Id = userId
	return err
}

func (u *User) ValidateCreds() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := DB.DB.QueryRow(query, u.Email)

	var retrivedPassword string
	err := row.Scan(&u.Id, &retrivedPassword)
	if err != nil {
		return errors.New("no result found")
	}
	isEqual := Utils.CheckPassword(u.Password, retrivedPassword)
	if !isEqual {
		return errors.New("invalid password")
	}
	return nil
}
