package models

import (
	"errors"
	"fmt"
	"time"

	"example.com/kuber/db"
)

// Cluster represents the cluster with basic details.
type Cluster struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name" binding:"required"`
	Provisioner        string    `json:"provisioner" binding:"required" validate:"oneof=aws azure gcp nutanix vmware"`
	Region             string    `json:"region" binding:"required"`
	WorkspaceID        string    `json:"workspace_id"`
	WorkspaceName      string    `json:"workspace_name"`
	Status             string    `json:"status"`
	UserID             int64     `json:"user_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	KubeconfigFilePath string    `json:"kubeconfig_file_path" binding:"required"`
}

func GetClusterIDByName(name string) (string, error) {
	query := `SELECT id FROM clusters WHERE name = $1`
	var id string
	err := db.DB.QueryRow(query, name).Scan(&id)
	if err != nil {
		return "", errors.New("workspace not found: " + err.Error())
	}
	return id, nil
}

func GetAllClusters() ([]map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace_id, kubeconfig, user_id, created_at, updated_at, status FROM clusters;`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, errors.New("cannot fetch the results. check the kubeconfig file")
	}
	defer rows.Close()

	var clusters []map[string]interface{}
	for rows.Next() {
		var c Cluster
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Provisioner,
			&c.Region,
			&c.WorkspaceID,
			&c.KubeconfigFilePath,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Status,
		)
		if err != nil {
			return nil, err
		}
		workspaceName, err := GetWorkspaceNameByID(c.WorkspaceID)
		if err != nil {
			return nil, errors.New("cannot fetch workspace name by ID: " + err.Error())
		}
		if workspaceName == "" {
			workspaceName = "default"
		}
		controlPlanes, workerNodes, storageContainers := fetchClusterDetailsFromKubeconfig(c.KubeconfigFilePath)
		clusterMap := map[string]interface{}{
			"cluster_info": map[string]interface{}{
				"id":                   c.ID,
				"name":                 c.Name,
				"provisioner":          c.Provisioner,
				"region":               c.Region,
				"workspace_name":       workspaceName,
				"status":               c.Status,
				"user_id":              c.UserID,
				"created_at":           c.CreatedAt,
				"updated_at":           c.UpdatedAt,
				"kubeconfig_file_path": c.KubeconfigFilePath,
			},
			"cluster_details": map[string]interface{}{
				"control_planes":     controlPlanes,
				"worker_nodes":       workerNodes,
				"storage_containers": storageContainers,
			},
		}
		clusters = append(clusters, clusterMap)
	}
	return clusters, nil
}

func (c *Cluster) SaveCluster() error {
	query := `INSERT INTO clusters (id, name, provisioner, region, workspace_id, kubeconfig, user_id, created_at, updated_at, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	err := db.DB.QueryRow(query,
		c.ID,
		c.Name,
		c.Provisioner,
		c.Region,
		c.WorkspaceID,
		c.KubeconfigFilePath,
		c.UserID,
		c.CreatedAt,
		c.UpdatedAt,
		c.Status).Scan(&c.ID)
	if err != nil {
		return errors.New("cannot create cluster: " + err.Error())
	}
	err = UpdateWorkspaceByClusterInfo(c.WorkspaceID)
	if err != nil {
		return errors.New("cannot update cluster count: " + err.Error())
	}
	return nil
}

func UpdateWorkspaceByClusterInfo(workspaceID string) error {
	updateQuery := `UPDATE workspaces SET cluster_count = cluster_count + 1 WHERE id = $1`
	_, err := db.DB.Exec(updateQuery, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to update workspace cluster count: %w", err)
	}
	return nil
}

func GetClusterByID(id string) (map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace_id, kubeconfig, user_id, created_at, updated_at, status 
			  FROM clusters WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var c Cluster
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Provisioner,
		&c.Region,
		&c.WorkspaceName,
		&c.KubeconfigFilePath,
		&c.UserID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Status,
	)
	if err != nil {
		return nil, errors.New("cannot fetch cluster by ID: " + err.Error())
	}

	controlPlanes, workerNodes, storageContainers := fetchClusterDetailsFromKubeconfig(c.KubeconfigFilePath)

	clusterMap := map[string]interface{}{
		"cluster_info": map[string]interface{}{
			"id":                   c.ID,
			"name":                 c.Name,
			"provisioner":          c.Provisioner,
			"region":               c.Region,
			"workspace_name":       c.WorkspaceName,
			"status":               c.Status,
			"user_id":              c.UserID,
			"created_at":           c.CreatedAt,
			"updated_at":           c.UpdatedAt,
			"kubeconfig_file_path": c.KubeconfigFilePath,
		},
		"cluster_details": map[string]interface{}{
			"control_planes":     controlPlanes,
			"worker_nodes":       workerNodes,
			"storage_containers": storageContainers,
		},
	}
	return clusterMap, nil
}

func GetClustersForUser(userID int64) ([]map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace_id, kubeconfig, user_id, created_at, updated_at, status 
			  FROM clusters WHERE user_id = $1`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, errors.New("cannot fetch clusters for user: " + err.Error())
	}
	defer rows.Close()

	var clusters []map[string]interface{}
	for rows.Next() {
		var c Cluster
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Provisioner,
			&c.Region,
			&c.WorkspaceName,
			&c.KubeconfigFilePath,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Status,
		)
		if err != nil {
			return nil, err
		}

		controlPlanes, workerNodes, storageContainers := fetchClusterDetailsFromKubeconfig(c.KubeconfigFilePath)

		clusterMap := map[string]interface{}{
			"cluster_info": map[string]interface{}{
				"id":                   c.ID,
				"name":                 c.Name,
				"provisioner":          c.Provisioner,
				"region":               c.Region,
				"workspace_name":       c.WorkspaceName,
				"status":               c.Status,
				"user_id":              c.UserID,
				"created_at":           c.CreatedAt,
				"updated_at":           c.UpdatedAt,
				"kubeconfig_file_path": c.KubeconfigFilePath,
			},
			"cluster_details": map[string]interface{}{
				"control_planes":     controlPlanes,
				"worker_nodes":       workerNodes,
				"storage_containers": storageContainers,
			},
		}
		clusters = append(clusters, clusterMap)
	}
	return clusters, nil
}

func GetClustersByWorkspaceName(workspace string) ([]map[string]interface{}, error) {
	workspaceID, err := GetWorkspaceIDByName(workspace)
	if err != nil {
		return nil, errors.New("cannot fetch workspace ID by name: " + err.Error())
	}
	if workspaceID == "" {
		return nil, errors.New("workspace not found")
	}
	query := `SELECT id, name, provisioner, region, workspace_id, kubeconfig, user_id, created_at, updated_at, status 
			  FROM clusters WHERE workspace_id = $1`
	rows, err := db.DB.Query(query, workspaceID)
	if err != nil {
		return nil, errors.New("cannot fetch clusters by workspace name: " + err.Error())
	}
	defer rows.Close()

	var clusters []map[string]interface{}
	for rows.Next() {
		var c Cluster
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Provisioner,
			&c.Region,
			&c.WorkspaceID,
			&c.KubeconfigFilePath,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Status,
		)
		if err != nil {
			return nil, err
		}
		controlPlanes, workerNodes, storageContainers := fetchClusterDetailsFromKubeconfig(c.KubeconfigFilePath)
		clusterMap := map[string]interface{}{
			"cluster_info": map[string]interface{}{
				"id":                   c.ID,
				"name":                 c.Name,
				"provisioner":          c.Provisioner,
				"region":               c.Region,
				"workspace_id":         c.WorkspaceID,
				"status":               c.Status,
				"user_id":              c.UserID,
				"created_at":           c.CreatedAt,
				"updated_at":           c.UpdatedAt,
				"kubeconfig_file_path": c.KubeconfigFilePath,
			},
			"cluster_details": map[string]interface{}{
				"control_planes":     controlPlanes,
				"worker_nodes":       workerNodes,
				"storage_containers": storageContainers,
			},
		}
		clusters = append(clusters, clusterMap)
	}
	return clusters, nil
}

func GetKubeconfigFilePathByID(id string) (string, error) {
	query := `SELECT kubeconfig FROM clusters WHERE id = $1`
	fmt.Printf("%s\n", id)
	row := db.DB.QueryRow(query, id)

	var kubeconfigFilePath string
	err := row.Scan(&kubeconfigFilePath)
	if err != nil {
		return "", errors.New("cannot fetch kubeconfig file path: " + err.Error())
	}
	return kubeconfigFilePath, nil
}

func GetProvisionerByID(id string) (string, error) {
	query := `SELECT provisioner FROM clusters WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var provisioner string
	err := row.Scan(&provisioner)
	if err != nil {
		return "", errors.New("cannot fetch provisioner: " + err.Error())
	}
	return provisioner, nil
}
