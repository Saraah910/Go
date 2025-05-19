package models

import (
	"errors"
	"strings"

	"example.com/kuber/Utils"
	"example.com/kuber/db"
)

type Users struct {
	Id            int64
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Role          string `json:"role"`
	OrgName       string `json:"org_name"`
	OrgDepartment string `json:"org_department"`
	CityLocation  string `json:"city_location"`
	Permission    string `json:"permission"`
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
		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.OrgName, &user.OrgDepartment, &user.CityLocation, &user.Permission)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *Users) CreateUser() error {
	query := `
	INSERT INTO users(email, password, role, org_name, org_department, city_location, permission) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id
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
	err = db.DB.QueryRow(query, u.Email, hashedPassword, u.Role, u.OrgName, u.OrgDepartment, u.CityLocation, u.Permission).Scan(&u.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("email already exists")
		}
		return err
	}
	return nil
}

func (u *Users) ValidateCreds() error {
	query := `SELECT id, password FROM users WHERE email = $1`
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

func GetPermission(userID int64) (string, string, error) {
	query := `SELECT role, permission FROM users WHERE id = $1`
	row := db.DB.QueryRow(query, userID)
	var role string
	var permission string
	err := row.Scan(&role, &permission)
	return role, permission, err
}

func DeleteUser(userID int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no user found with the given ID")
	}

	return nil
}
