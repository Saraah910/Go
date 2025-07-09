package routes

import (
	"example.com/kuber/Middleware"
	"github.com/gin-gonic/gin"
)

func Routes(server *gin.Engine) {
	server.GET("/users/list", getAllUsers)
	server.POST("/users/signup", signUpUser)
	server.POST("/users/login", loginUser)
	server.POST("/users/logout", logoutUser)
	server.GET("/kube/clusters/list", getAllClusters)
	server.GET("/kube/cluster/:id", getClusterByID)
	server.GET("/kube/infrastructure/list", getAllInfrastructures)
	// server.GET("/kube/infrastructure/user/:id", getInfraForUser)
	server.GET("/kube/workspaces/list", getWorkspaces)

	authenticationGroups := server.Group("/")
	authenticationGroups.Use(Middleware.Authentication)
	authenticationGroups.POST("/kube/clusters/list", addCluster)
	authenticationGroups.PUT("/users/:id/update", updateUser)
	authenticationGroups.GET("/kube/clusters/list/user", getClustersForUser)
	authenticationGroups.POST("/kube/workspaces/list", createWorkspace)
	authenticationGroups.POST("/kube/cluster/create", createCluster)
	authenticationGroups.GET("/kube/cluster/workspaces/:workspace_uuid", getClustersByWorkspace)
	authenticationGroups.GET("/kube/users/list/permission", getPermissions)
	authenticationGroups.DELETE("/users/:id/delete", deleteUser)
	// authenticationGroups.DELETE("/kube/clusters/:id/delete", deleteCluster)
	authenticationGroups.POST("/kube/infrastructure/aws", AWSInfra)
	authenticationGroups.POST("/kube/infrastructure/nutanix", NutanixInfra)
	// authenticationGroups.POST("/kube/infrastructure/vmware", VMwareInfra)

	authenticationGroups.POST("/kube/dr/apply", clusterDR)
	clusterActionGroup := server.Group("/kube/cluster/actions")
	clusterActionGroup.Use(Middleware.Authentication)
	clusterActionGroup.GET(":id/services", getServices)
	// clusterActionGroup.GET("/pods", getPods)
	// clusterActionGroup.GET("/deployments", getDeployments)
	// clusterActionGroup.GET("/pvcs", getPVCs)
	// clusterActionGroup.GET("/pv", getPV)
	// clusterActionGroup.GET("/statefulsets", getStorageContainers)
	// clusterActionGroup.GET("/secrets", getSecrets)
	// clusterActionGroup.GET("/configmaps", getConfigMaps)
	// clusterActionGroup.GET("/ingresses", getIngresses)
	// clusterActionGroup.GET("/roles", getRoles)
	// clusterActionGroup.GET("/rolebindings", getRoleBindings)
	// clusterActionGroup.GET("/clusterroles", getClusterRoles)
	// clusterActionGroup.GET("/clusterrolebindings", getClusterRoleBindings)
	clusterActionGroup.GET(":id/namespaces", getNamespaces)
	// clusterActionGroup.GET("/events", getEvents)

}
