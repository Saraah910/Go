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
	InfraID   string          `json:"infra_id"`
	Name      string          `json:"name"`
	Provider  string          `json:"provider"`
	IsDefault bool            `json:"is_default"`
	CreatedAt time.Time       `json:"created_at"`
	UserID    int64           `json:"user_id"`
	Config    json.RawMessage `json:"config"`
}

type NutanixInfra struct {
	Name        string `json:"infra_name" binding:"required"`
	Endpoint    string `json:"endpoint" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	Insecure    bool   `json:"insecure"`
	UserID      int64  `json:"user_id"`
}

func (n *NutanixInfra) SaveNutanixInfra() (string, error) {
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
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	infraID, err := db.GetUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	query := `INSERT INTO infra (id, name, provider, config, user_id) VALUES($1, $2, $3, $4, $5)`
	_, err = db.DB.Exec(query, infraID, n.Name, "nutanix", configJson, n.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to insert into DB: %w", err)
	}

	return infraID, nil
}

func (i *AWSInfra) SaveAWSInfra() (string, error) {
	config := map[string]interface{}{
		"infra_name": i.Name,
		"access_key": i.AccessKey,
		"secret_key": i.SecretKey,
		"region":     i.Region,
		"vpc_id":     i.VPC_ID,
		"subnet_ids": i.SubnetIDs,
		"user_id":    i.UserID,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}
	infraID, err := db.GetUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	query := `INSERT INTO infra (id, name, provider, config, user_id) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.DB.Exec(query, infraID, i.Name, "aws", configJSON, i.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to add into DB: %w", err)
	}
	return infraID, nil
}

func GetInfrastructures() ([]InfraList, error) {
	query := `SELECT id, name, provider, is_default, created_at, config, user_id FROM infra`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	var infraList []InfraList
	for rows.Next() {
		var infra InfraList
		err := rows.Scan(&infra.InfraID, &infra.Name, &infra.Provider, &infra.IsDefault, &infra.CreatedAt, &infra.Config, &infra.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		infraList = append(infraList, infra)
	}
	return infraList, nil
}

func GetInfraByName(name string) (*InfraList, error) {
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

	var infra InfraList
	err = json.Unmarshal(configJSON, &infra)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	infra.Name = infraName
	infra.UserID = userID

	return &infra, nil
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
		err := rows.Scan(&infra.InfraID, &infra.Name, &infra.Provider, &infra.IsDefault, &infra.CreatedAt, &infra.Config, &infra.UserID)
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
