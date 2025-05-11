package models

import (
	"errors"

	"example.com/kuber/Utils"
	"example.com/kuber/db"
)

type Users struct {
	Id       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func FetchAllUsers() ([]Users, error) {
	query := `SELECT * FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, errors.New("cannot fetch results from database")
	}
	defer rows.Close()
	var users []Users

	for rows.Next() {
		var user Users
		err := rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *Users) CreateUser() error {
	query := `
	INSERT INTO users(email, password) VALUES(?,?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	hashedPassword, err := Utils.ConvertToHashString(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}
	userID, err := result.LastInsertId()
	u.Id = userID
	return err
}

func (u *Users) ValidateCreds() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var retrivedPassword string
	row := stmt.QueryRow(u.Email)
	err = row.Scan(&u.Id, &retrivedPassword)
	if err != nil {
		return err
	}
	if !Utils.ComparePasswords(u.Password, retrivedPassword) {
		return errors.New("password does not match")
	}
	return nil
}
