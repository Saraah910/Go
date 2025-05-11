package users

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	fname     string
	lname     string
	dob       string
	createdAt time.Time
}

func New(firstName, lastName, dob string) (*User, error) {
	if firstName == "" || lastName == "" || dob == "" {
		err := errors.New("invalid input")
		return nil, err
	}

	user := User{
		fname:     firstName,
		lname:     lastName,
		dob:       dob,
		createdAt: time.Now(),
	}
	return &user, nil
}

func (user *User) ShowOutput() {
	fmt.Printf("The user details are %v, %v, DOB: %v \n", user.fname, user.lname, user.dob)
	fmt.Printf("User created at: %v", user.createdAt)
}
