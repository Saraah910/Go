package Cloud

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"

// 	"example.com/kuber/models"
// )

// type AWSConfig struct {
// 	AccessKey string   `json:"access_key"`
// 	SecretKey string   `json:"secret_key"`
// 	VPCID     string   `json:"vpc_id"`
// 	Subnets   []string `json:"subnet_ids"`
// }

// func AWSProvision(cluster models.Clusters) error {
// 	var cfg AWSConfig
// 	if err := json.Unmarshal(cluster.ProviderConfig, &cfg); err != nil {
// 		return fmt.Errorf("invalid AWS config: %w", err)
// 	}

// 	_, err := os.MkdirTemp("", "aws-tf-*")
// 	// if err != nil {
// 	// 	return fmt.Errorf("cannot create temp dir: %w", err)
// 	// }

// 	// if err != nil {
// 	// 	return fmt.Errorf("cannot generate AWS Terraform files: %w", err)
// 	// }
// 	return err
// 	// return RunFiles(tmpDir)
// }
