package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example.com/kuber/db"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type ResourceMetric struct {
	Used     string `json:"used"`
	Capacity string `json:"capacity"`
	Unit     string `json:"unit"`
}

type ResourceUsage struct {
	ClusterID string         `json:"cluster_uuid"`
	CPU       ResourceMetric `json:"cpu"`
	Memory    ResourceMetric `json:"memory"`
	Storage   ResourceMetric `json:"storage"`
}

type Workspace struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	OwnerID           int64             `json:"owner_id"`
	CreatedAt         time.Time         `json:"created_at"`
	Members           []string          `json:"members"`
	Roles             []string          `json:"roles"`
	ClusterCount      int               `json:"cluster_count"`
	CloudProviders    []string          `json:"cloud_providers"`
	AppsCount         int               `json:"apps_count"`
	ResourceUsage     []ResourceUsage   `json:"resource_usage"`
	MonitoringEnabled bool              `json:"monitoring_enabled"`
	LoggingEnabled    bool              `json:"logging_enabled"`
	Tags              map[string]string `json:"tags"`
	Clusters          []Cluster         `json:"clusters"`
}

func (w *Workspace) Save() error {
	workspaceID, err := db.GetUUID()
	if err != nil {
		return err
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
	}
	w.ID = workspaceID
	w.ClusterCount = 0
	w.AppsCount = 0
	w.MonitoringEnabled = false
	w.LoggingEnabled = false
	// if w.Tags == nil {
	// 	w.Tags = pq.StringArray{}
	// }
	// if w.CloudProviders == nil {
	// 	w.CloudProviders = pq.StringArray{}
	// }

	query := `INSERT INTO workspaces (
		id, name, description, owner_id, created_at, monitoring_enabled, logging_enabled
	) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = db.DB.QueryRow(query,
		w.ID,
		w.Name,
		w.Description,
		w.OwnerID,
		w.CreatedAt,
		w.MonitoringEnabled,
		w.LoggingEnabled,
	).Scan(&w.ID)

	if err != nil {
		return err
	}
	return nil
}

func GetWorkspaceIDByName(name string) (string, error) {
	query := `SELECT id FROM workspaces WHERE name = $1`
	var id string
	err := db.DB.QueryRow(query, name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func GetWorkspaceNameByID(id string) (string, error) {
	query := `SELECT name FROM workspaces WHERE id = $1`
	var name string
	err := db.DB.QueryRow(query, id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetWorkspaces() ([]*Workspace, error) {
	query := `SELECT id, name, description, owner_id, created_at,
       members, roles, cluster_count, cloud_providers,
       apps_count, monitoring_enabled,
       logging_enabled, tags
FROM workspaces;`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*Workspace
	for rows.Next() {
		var ws Workspace
		var membersJSON, rolesJSON, cloudProvidersJSON, tagsJSON []byte
		err := rows.Scan(
			&ws.ID,
			&ws.Name,
			&ws.Description,
			&ws.OwnerID,
			&ws.CreatedAt,
			&membersJSON,
			&rolesJSON,
			&ws.ClusterCount,
			&cloudProvidersJSON,
			&ws.AppsCount,
			&ws.MonitoringEnabled,
			&ws.LoggingEnabled,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(membersJSON, &ws.Members); err != nil {
			return nil, fmt.Errorf("error unmarshaling members: %w", err)
		}
		if err := json.Unmarshal(rolesJSON, &ws.Roles); err != nil {
			return nil, fmt.Errorf("error unmarshaling roles: %w", err)
		}
		if err := json.Unmarshal(cloudProvidersJSON, &ws.CloudProviders); err != nil {
			return nil, fmt.Errorf("error unmarshaling cloud providers: %w", err)
		}

		if err := json.Unmarshal(tagsJSON, &ws.Tags); err != nil {
			return nil, fmt.Errorf("error unmarshaling tags: %w", err)
		}
		clusters, err := GetClustersByWorkspaceID(ws.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting clusters: %w", err)
		}
		var ResourceUsageList []ResourceUsage
		for _, cluster := range clusters {
			clientSet, metricSet, err := CalculateClusterResourceUsage(cluster.ID)
			if err != nil {
				return nil, err
			}
			resourceUsage, err := GetClusterResourceUsage(cluster.ID, clientSet, metricSet)
			if err != nil {
				return nil, err
			}
			ResourceUsageList = append(ResourceUsageList, *resourceUsage)
		}
		ws.ResourceUsage = ResourceUsageList
		workspaces = append(workspaces, &ws)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func GetWorkspaceByID(id string) (*Workspace, error) {
	query := `SELECT id, name, description, owner_id, created_at, members, roles, cluster_count, cloud_providers, apps_count, resource_usage, monitoring_enabled, logging_enabled, tags FROM workspaces WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var w Workspace
	var membersJSON, rolesJSON, cloudProvidersJSON, tagsJSON []byte
	err := row.Scan(
		&w.ID,
		&w.Name,
		&w.Description,
		&w.OwnerID,
		&w.CreatedAt,
		&membersJSON,
		&rolesJSON,
		&w.ClusterCount,
		&cloudProvidersJSON,
		&w.AppsCount,
		&w.ResourceUsage,
		&w.MonitoringEnabled,
		&w.LoggingEnabled,
		&tagsJSON,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(membersJSON, &w.Members); err != nil {
		return nil, fmt.Errorf("error unmarshaling members: %w", err)
	}
	if err := json.Unmarshal(rolesJSON, &w.Roles); err != nil {
		return nil, fmt.Errorf("error unmarshaling roles: %w", err)
	}
	if err := json.Unmarshal(cloudProvidersJSON, &w.CloudProviders); err != nil {
		return nil, fmt.Errorf("error unmarshaling cloud providers: %w", err)
	}

	if err := json.Unmarshal(tagsJSON, &w.Tags); err != nil {
		return nil, fmt.Errorf("error unmarshaling tags: %w", err)
	}
	return &w, nil
}

func GetWorkspaceByUUID(uuid string) (*Workspace, error) {
	query := `SELECT id, name, description, owner_id, created_at, members, roles, cluster_count, cloud_providers, apps_count, monitoring_enabled, logging_enabled, tags FROM workspaces WHERE id = $1`
	row := db.DB.QueryRow(query, uuid)

	var w Workspace
	var membersJSON, rolesJSON, cloudProvidersJSON, tagsJSON []byte
	err := row.Scan(
		&w.ID,
		&w.Name,
		&w.Description,
		&w.OwnerID,
		&w.CreatedAt,
		&membersJSON,
		&rolesJSON,
		&w.ClusterCount,
		&cloudProvidersJSON,
		&w.AppsCount,
		&w.MonitoringEnabled,
		&w.LoggingEnabled,
		&tagsJSON,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(membersJSON, &w.Members); err != nil {
		return nil, fmt.Errorf("error unmarshaling members: %w", err)
	}
	if err := json.Unmarshal(rolesJSON, &w.Roles); err != nil {
		return nil, fmt.Errorf("error unmarshaling roles: %w", err)
	}
	if err := json.Unmarshal(cloudProvidersJSON, &w.CloudProviders); err != nil {
		return nil, fmt.Errorf("error unmarshaling cloud providers: %w", err)
	}

	if err := json.Unmarshal(tagsJSON, &w.Tags); err != nil {
		return nil, fmt.Errorf("error unmarshaling tags: %w", err)
	}
	clusters, err := GetClustersByWorkspaceID(w.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting clusters: %w", err)
	}
	var clustersList []Cluster
	for _, cluster := range clusters {
		clustersList = append(clustersList, *cluster)
	}
	w.Clusters = clustersList
	return &w, nil
}
func GetWorkspacesForUser(userID string) ([]*Workspace, error) {
	query := `SELECT id, name, description, owner_id, created_at, members, roles, cluster_count, cloud_providers, apps_count, resource_usage, monitoring_enabled, logging_enabled, tags FROM workspaces WHERE owner_id = $1`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*Workspace
	for rows.Next() {
		var w Workspace
		if err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.Description,
			&w.OwnerID,
			&w.CreatedAt,
			&w.Members,
			&w.Roles,
			&w.ClusterCount,
			&w.CloudProviders,
			&w.AppsCount,
			&w.ResourceUsage,
			&w.MonitoringEnabled,
			&w.LoggingEnabled,
			&w.Tags,
		); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, &w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workspaces, nil
}
func (w *Workspace) Update() error {
	query := `UPDATE workspaces SET 
		name = $1, 
		description = $2, 
		owner_id = $3, 
		members = $4, 
		roles = $5, 
		cluster_count = $6, 
		cloud_providers = $7, 
		apps_count = $8, 
		resource_usage = $9, 
		monitoring_enabled = $10, 
		logging_enabled = $11, 
		tags = $12 
	WHERE id = $13`
	_, err := db.DB.Exec(query,
		w.Name,
		w.Description,
		w.OwnerID,
		w.Members,
		w.Roles,
		w.ClusterCount,
		w.CloudProviders,
		w.AppsCount,
		w.ResourceUsage,
		w.MonitoringEnabled,
		w.LoggingEnabled,
		w.Tags,
		w.ID,
	)
	return err
}
func (w *Workspace) Delete() error {
	query := `DELETE FROM workspaces WHERE id = $1`
	_, err := db.DB.Exec(query, w.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetClustersByWorkspaceID(workspaceID string) ([]*Cluster, error) {
	query := `SELECT id,name,provisioner,region,kubeconfig,created_at,user_id,status FROM clusters WHERE workspace_id = $1`
	rows, err := db.DB.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []*Cluster
	for rows.Next() {
		var c Cluster
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Provisioner,
			&c.Region,
			&c.KubeconfigFilePath,
			&c.CreatedAt,
			&c.UserID,
			&c.Status,
		)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, &c)
	}
	return clusters, nil
}

func CalculateClusterResourceUsage(clusterID string) (*kubernetes.Clientset, *versioned.Clientset, error) {
	path, err := GetKubeconfigFilePathByID(clusterID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot fetch kubeconfig clientset: %w", err)
	}
	clientSet, err := GetKubeClient(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create clientset: %w", err)
	}
	metricSet, err := GetMetricClient(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create metricset: %w", err)
	}
	return clientSet, metricSet, nil
}

func GetClusterResourceUsage(clusterID string, clientset *kubernetes.Clientset, metricsClient *versioned.Clientset) (*ResourceUsage, error) {
	ctx := context.Background()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}

	nodeMetricsList, err := metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %v", err)
	}

	totalCPUUsed := resource.NewQuantity(0, resource.DecimalSI)
	totalMemUsed := resource.NewQuantity(0, resource.BinarySI)
	for _, m := range nodeMetricsList.Items {
		totalCPUUsed.Add(*m.Usage.Cpu())
		totalMemUsed.Add(*m.Usage.Memory())
	}

	totalCPUCap := resource.NewQuantity(0, resource.DecimalSI)
	totalMemCap := resource.NewQuantity(0, resource.BinarySI)
	for _, node := range nodes.Items {
		totalCPUCap.Add(*node.Status.Capacity.Cpu())
		totalMemCap.Add(*node.Status.Capacity.Memory())
	}
	pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PVs: %v", err)
	}
	totalStorageCap := resource.NewQuantity(0, resource.BinarySI)
	for _, pv := range pvs.Items {
		if pv.Spec.Capacity != nil {
			if qty, ok := pv.Spec.Capacity[corev1.ResourceStorage]; ok {
				totalStorageCap.Add(qty)
			}
		}
	}

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PVCs: %v", err)
	}
	totalStorageUsed := resource.NewQuantity(0, resource.BinarySI)
	for _, pvc := range pvcs.Items {
		if pvc.Status.Capacity != nil {
			if qty, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
				totalStorageUsed.Add(qty)
			}
		}
	}

	usage := &ResourceUsage{
		ClusterID: clusterID,
		CPU: ResourceMetric{
			Used:     totalCPUUsed.String(),
			Capacity: totalCPUCap.String(),
			Unit:     "cores",
		},
		Memory: ResourceMetric{
			Used:     totalMemUsed.String(),
			Capacity: totalMemCap.String(),
			Unit:     "bytes",
		},
		Storage: ResourceMetric{
			Used:     totalStorageUsed.String(),
			Capacity: totalStorageCap.String(),
			Unit:     "bytes",
		},
	}

	return usage, nil
}
