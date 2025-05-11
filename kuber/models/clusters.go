package models

import (
	"errors"

	"example.com/kuber/db"
)

type Kube struct {
	ClusterID          int64
	ClusterName        string `json:"cluster_name" binding:"required"`
	Provisioner        string `json:"provisioner" binding:"required"`
	KubeconfigFilePath string `json:"kubeconfig_path" binding:"required"`
	Status             string `json:"status"`
	UserID             int64
}

func GetAllClusters() ([]Kube, error) {
	query := `SELECT * FROM kubeclusters`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, errors.New("cannot fetch the results")
	}
	defer rows.Close()
	var Clusters []Kube

	for rows.Next() {
		var Cluster Kube

		err := rows.Scan(&Cluster.ClusterID, &Cluster.ClusterName, &Cluster.Provisioner, &Cluster.KubeconfigFilePath, &Cluster.Status, &Cluster.UserID)
		if err != nil {
			return nil, err
		}
		Clusters = append(Clusters, Cluster)
	}
	return Clusters, nil
}

func (k *Kube) SaveCluster() error {
	query := `
	INSERT INTO kubeclusters(cluster_name,provisioner,kubeconfig_path,status,user_id) VALUES(?,?,?,?,?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(k.ClusterName, k.Provisioner, k.KubeconfigFilePath, k.Status, k.UserID)
	if err != nil {
		return err
	}
	clusterID, err := result.LastInsertId()
	k.ClusterID = clusterID

	return err
}

func GetClusterById(clusterID int64) (*Kube, error) {
	query := `SELECT * FROM kubeclusters WHERE cluster_id = ?`
	row := db.DB.QueryRow(query, clusterID)
	var cluster Kube

	err := row.Scan(&cluster.ClusterID, &cluster.ClusterName, &cluster.Provisioner, &cluster.KubeconfigFilePath, &cluster.Status, &cluster.UserID)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}
