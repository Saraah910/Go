package routes

import (
	"example.com/kuber/models"
)

func isAuthorized(userID int64) (bool, error) {
	role, permission, err := models.GetPermission(userID)
	if err != nil {
		return false, err
	}
	if role != "admin" {
		if permission != "write" && permission != "full" {
			return false, nil
		}
	}
	return true, nil
}
