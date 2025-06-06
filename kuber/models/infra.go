package models

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/kuber/db"
)

type AWSInfra struct {
	Name      string   `json:"infra_name" binding:"required"`
	AccessKey string   `json:"access_key" binding:"required"`
	SecretKey string   `json:"secret_key" binding:"required"`
	Region    string   `json:"region" binding:"required"`
	VPC_ID    string   `json:"vpc_id"`
	SubnetIDs []string `json:"subnet_ids"`
	UserID    int64    `json:"user_id"`
}

type InfraList struct {
	InfraID   int64
	Name      string
	Provider  string
	IsDefault bool
	CreatedAt time.Time
	UserID    int64
}

type NutanixInfra struct {
	Name        string `json:"infra_name" binding:"required"`
	Endpoint    string `json:"endpoint" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	Insecure    bool   `json:"insecure"`
	UserID      int64  `json:"user_id"`
}

func (n *NutanixInfra) SaveNutanixInfra() (int64, error) {
	config := map[string]interface{}{
		"infra_name":   n.Name,
		"endpoint":     n.Endpoint,
		"port":         n.Port,
		"cluster_name": n.ClusterName,
		"insecure":     n.Insecure,
		"user_id":      n.UserID,
	}
	configJson, err := json.Marshal(config)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal config: %w", err)
	}
	query := `INSERT INTO infra (name,provider,config,user_id) VALUES($1,$2,$3,$4) RETURNING id`
	var nutanix_infra_id int64
	err = db.DB.QueryRow(query, &n.Name, "nutanix", configJson, &n.UserID).Scan(&nutanix_infra_id)
	if err != nil {
		return 0, fmt.Errorf("failed to add into DB: %w", err)
	}
	return nutanix_infra_id, nil
}

func (i *AWSInfra) SaveAWSInfra() (int64, error) {
	config := map[string]interface{}{
		"access_key": i.AccessKey,
		"secret_key": i.SecretKey,
		"region":     i.Region,
		"vpc_id":     i.VPC_ID,
		"subnet_ids": i.SubnetIDs,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal config: %w", err)
	}
	query := `INSERT INTO infra (name,provider,config,is_default,user_id) VALUES($1,$2,$3,$4,$5) RETURNING id`
	var aws_cluster_id int64
	err = db.DB.QueryRow(query, i.Name, "aws", configJSON, "No", i.UserID).Scan(&aws_cluster_id)
	if err != nil {
		return 0, fmt.Errorf("failed to add into DB: %w", err)
	}
	return aws_cluster_id, nil
}

func GetInfrastructures() ([]InfraList, error) {
	query := `SELECT id, name, provider, is_default, created_at, user_id FROM infra`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	var infraList []InfraList
	for rows.Next() {
		var infra InfraList
		err := rows.Scan(&infra.InfraID, &infra.Name, &infra.Provider, &infra.IsDefault, &infra.CreatedAt, &infra.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		infraList = append(infraList, infra)
	}
	return infraList, nil
}

func GetInfraByName(name string) (*NutanixInfra, error) {
	query := `SELECT name, config, user_id FROM infra WHERE name = $1`

	var (
		infraName  string
		configJSON []byte
		userID     int64
	)

	row := db.DB.QueryRow(query, name)
	err := row.Scan(&infraName, &configJSON, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query infra: %w", err)
	}

	var nutanixInfra NutanixInfra
	err = json.Unmarshal(configJSON, &nutanixInfra)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	nutanixInfra.Name = infraName
	nutanixInfra.UserID = userID

	return &nutanixInfra, nil
}

func GetInfraByUserID(userID int64) ([]InfraList, error) {
	query := `SELECT id, name, provider, is_default, created_at, user_id, config FROM infra WHERE user_id = $1`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	var userSpecificInfraList []InfraList
	for rows.Next() {
		var infra InfraList
		err := rows.Scan(&infra.InfraID, &infra.Name, &infra.Provider, &infra.IsDefault, &infra.CreatedAt, &infra.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		userSpecificInfraList = append(userSpecificInfraList, infra)
	}
	return userSpecificInfraList, nil
}

func (i InfraList) GetConfig() (interface{}, error) {
	query := `SELECT config FROM infra WHERE provider = $1`
	row := db.DB.QueryRow(query, i.Provider)

	var configJSON []byte
	if err := row.Scan(&configJSON); err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	switch i.Provider {
	case "aws":
		var infra AWSInfra
		if err := json.Unmarshal(configJSON, &infra); err != nil {
			return infra, err
		}
		return infra, nil

	case "nutanix":
		var infra NutanixInfra
		if err := json.Unmarshal(configJSON, &infra); err != nil {
			return infra, err
		}
		return infra, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", i.Provider)
	}
}
