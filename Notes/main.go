package main

import (
	"fmt"

	"example.com/notes/users"
)

func getText(text string) string {
	var input string
	for {
		fmt.Print(text)
		n, err := fmt.Scanln(&input)
		if err != nil || n == 0 || input == "" {
			fmt.Println("Input cannot be empty. Please try again.")
			continue
		}
		break
	}
	return input
}

func getData() *users.User {
	firstName := getText("Enter first name: ")
	lastName := getText("Entre last name: ")
	dob := getText("Enter DOB in the format (DD/MM/YYYY): ")

	user, err := users.New(firstName, lastName, dob)

	if err != nil {
		fmt.Print(err)
		return nil
	}
	return user

}
func main() {
	user := getData()
	if user != nil {
		user.ShowOutput()
	}
}
