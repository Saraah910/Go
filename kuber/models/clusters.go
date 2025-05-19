package models

import (
	"errors"
	"strings"

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
	INSERT INTO kubeclusters(cluster_name,provisioner,kubeconfig_path,status,user_id) VALUES($1,$2,$3,$4,$5) RETURNING cluster_id
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = db.DB.QueryRow(query, k.ClusterName, k.Provisioner, k.KubeconfigFilePath, k.Status, k.UserID).Scan(&k.ClusterID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("cluster already exists")
		}
		return err
	}
	return nil
}

func GetClusterById(clusterID int64) (*Kube, error) {
	query := `SELECT * FROM kubeclusters WHERE cluster_id = $1`
	row := db.DB.QueryRow(query, clusterID)
	var cluster Kube

	err := row.Scan(&cluster.ClusterID, &cluster.ClusterName, &cluster.Provisioner, &cluster.KubeconfigFilePath, &cluster.Status, &cluster.UserID)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Kube) UpdateCluster(clusterID int64) error {
	query := `UPDATE kubeclusters SET cluster_name = $1, provisioner = $2, kubeconfig_path = $3 WHERE cluster_id = $4`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ClusterName, c.Provisioner, c.KubeconfigFilePath, c.ClusterID)
	return err
}
