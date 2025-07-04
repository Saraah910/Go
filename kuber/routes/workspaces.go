package routes

import (
	"net/http"

	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
)

func getWorkspaces(c *gin.Context) {
	workspaces, err := models.GetWorkspaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspaces", "Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workspaces": workspaces})
}
func getWorkspaceByID(c *gin.Context) {
	workspaceID := c.Param("id")
	workspace, err := models.GetWorkspaceByID(workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspace"})
		return
	}
	if workspace == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workspace": workspace})
}
func createWorkspace(c *gin.Context) {
	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	workspace.OwnerID = c.GetInt64("userID")
	if err := workspace.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace", "Error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Workspace created successfully", "workspace_id": workspace.ID})
}

func updateWorkspace(c *gin.Context) {
	workspaceID := c.Param("id")
	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	workspace.ID = workspaceID // Ensure the ID is set for the update
	if err := workspace.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workspace updated successfully"})
}
