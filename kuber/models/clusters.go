package models

import (
	"errors"

	"example.com/kuber/db"
)

// Cluster represents the cluster with basic details.
type Cluster struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name" binding:"required"`
	Provisioner        string `json:"provisioner" binding:"required" validate:"oneof=aws azure gcp nutanix vmware"`
	Region             string `json:"region" binding:"required"`
	Workspace          string `json:"workspace" `
	Status             string `json:"status"`
	UserID             int64  `json:"user_id"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	KubeconfigFilePath string `json:"kubeconfig_file_path" binding:"required"`
}

// GetAllClusters retrieves all clusters and fetches control planes, worker nodes, and storage containers from kubeconfig.
func GetAllClusters() ([]map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace, kubeconfig, user_id, created_at, updated_at, status FROM clusters`
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
			&c.Workspace,
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
			"id":                   c.ID,
			"name":                 c.Name,
			"provisioner":          c.Provisioner,
			"region":               c.Region,
			"workspace":            c.Workspace,
			"status":               c.Status,
			"user_id":              c.UserID,
			"created_at":           c.CreatedAt,
			"updated_at":           c.UpdatedAt,
			"kubeconfig_file_path": c.KubeconfigFilePath,
			"control_planes":       controlPlanes,
			"worker_nodes":         workerNodes,
			"storage_containers":   storageContainers,
		}
		clusters = append(clusters, clusterMap)
	}
	return clusters, nil
}

func (c *Cluster) SaveCluster() error {
	query := `INSERT INTO clusters (name, provisioner, region, workspace, kubeconfig, user_id, created_at, updated_at, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err := db.DB.QueryRow(query,
		c.Name,
		c.Provisioner,
		c.Region,
		c.Workspace,
		c.KubeconfigFilePath,
		c.UserID,
		c.CreatedAt,
		c.UpdatedAt,
		c.Status).Scan(&c.ID)
	if err != nil {
		return errors.New("cannot create cluster: " + err.Error())
	}
	return nil
}

func GetClusterByID(id int64) (map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace, kubeconfig, user_id, created_at, updated_at, status 
			  FROM clusters WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var c Cluster
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Provisioner,
		&c.Region,
		&c.Workspace,
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
		"id":                   c.ID,
		"name":                 c.Name,
		"provisioner":          c.Provisioner,
		"region":               c.Region,
		"workspace":            c.Workspace,
		"status":               c.Status,
		"user_id":              c.UserID,
		"created_at":           c.CreatedAt,
		"updated_at":           c.UpdatedAt,
		"kubeconfig_file_path": c.KubeconfigFilePath,
		"control_planes":       controlPlanes,
		"worker_nodes":         workerNodes,
		"storage_containers":   storageContainers,
	}
	return clusterMap, nil
}

func GetClustersForUser(userID int64) ([]map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace, kubeconfig, user_id, created_at, updated_at, status 
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
			&c.Workspace,
			&c.KubeconfigFilePath,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Status,
		)
		if err != nil {
			return nil, err
		}

		clusterMap := map[string]interface{}{
			"id":                   c.ID,
			"name":                 c.Name,
			"provisioner":          c.Provisioner,
			"region":               c.Region,
			"workspace":            c.Workspace,
			"status":               c.Status,
			"user_id":              c.UserID,
			"created_at":           c.CreatedAt,
			"updated_at":           c.UpdatedAt,
			"kubeconfig_file_path": c.KubeconfigFilePath,
		}
		clusters = append(clusters, clusterMap)
	}
	return clusters, nil
}

func GetClustersByWorkspaceName(workspace string) ([]map[string]interface{}, error) {
	query := `SELECT id, name, provisioner, region, workspace, kubeconfig, user_id, created_at, updated_at, status 
			  FROM clusters WHERE workspace = $1`
	rows, err := db.DB.Query(query, workspace)
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
			&c.Workspace,
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
			"id":                   c.ID,
			"name":                 c.Name,
			"provisioner":          c.Provisioner,
			"region":               c.Region,
			"workspace":            c.Workspace,
			"status":               c.Status,
			"user_id":              c.UserID,
			"created_at":           c.CreatedAt,
			"updated_at":           c.UpdatedAt,
			"kubeconfig_file_path": c.KubeconfigFilePath,
			"control_planes":       controlPlanes,
			"worker_nodes":         workerNodes,
			"storage_containers":   storageContainers,
		}
		clusters = append(clusters, clusterMap)
	}
	return clusters, nil
}

func GetKubeconfigFilePathByID(id int64) (string, error) {
	query := `SELECT kubeconfig FROM clusters WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var kubeconfigFilePath string
	err := row.Scan(&kubeconfigFilePath)
	if err != nil {
		return "", errors.New("cannot fetch kubeconfig file path: " + err.Error())
	}
	return kubeconfigFilePath, nil
}

func GetProvisionerByID(id int64) (string, error) {
	query := `SELECT provisioner FROM clusters WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var provisioner string
	err := row.Scan(&provisioner)
	if err != nil {
		return "", errors.New("cannot fetch provisioner: " + err.Error())
	}
	return provisioner, nil
}
