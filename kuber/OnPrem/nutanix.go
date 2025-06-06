package OnPrem

// func CreateVMs(providerConfig, machineConfig json.RawMessage) {
// 	ClusterConfig := json.Unmarshal(cluster.ClusterConfig, &models.ClusterConfig)
// 	ProviderConfig := json.Unmarshal(cluster.ProviderConfig, &models.ProviderConfig)
// 	MachineConfig := json.Unmarshal(cluster.MachineConfig, &models.MachineConfig)

// 	infra, err := models.GetInfraByName(ClusterConfig.Provider)
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
// 		return
// 	}
// 	Endpoint, Port, Insecure, PE := infra.Endpoint, infra.Port, infra.Insecure, infra.ClusterName

// }
