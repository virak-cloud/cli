package responses

type KubernetesCluster struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ZoneID            string `json:"zone_id"`
	Status            string `json:"status"`
	FailedReason      string `json:"failure_message"`
	KubernetesVersion struct {
		ID      string `json:"id"`
		Version string `json:"version"`
		Enabled bool   `json:"enabled"`
	} `json:"kubernetes_version"`
	ServiceOffering struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"service_offering"`
	SSHKey      string `json:"ssh_key"`
	HAEnabled   bool   `json:"ha_enabled"`
	ClusterSize int    `json:"cluster_size"`
	CreatedAt   int    `json:"created_at"`
	UpdatedAt   int    `json:"updated_at"`
}

type KubernetesClusterResponse struct {
	Data KubernetesCluster `json:"data"`
}

type KubernetesClusterListResponse struct {
	Data []KubernetesCluster `json:"data"`
}

type KubernetesMessage struct {
	Message string `json:"message"`
}

type KubernetesVersion struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
}

type KubernetesVersionsListResponse struct {
	Data []KubernetesVersion `json:"data"`
}

type KubernetesEvent struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

type KubernetesEventsListResponse struct {
	Data []KubernetesEvent `json:"data"`
}

type KubernetesServiceOffering struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsPublic    bool   `json:"is_public"`
	IsAvailable bool   `json:"is_available"`
	HourlyPrice struct {
		Up   int `json:"up"`
		Down int `json:"down"`
	} `json:"hourly_price"`
	HourlyPriceNoDiscount struct {
		Up   int `json:"up"`
		Down int `json:"down"`
	} `json:"hourly_price_no_discount"`
	Description string `json:"description"`
	Hardware    struct {
		CPUCore        int `json:"cpu_core"`
		MemoryMB       int `json:"memory_mb"`
		CPUSpeedMHz    int `json:"cpu_speed_MHz"`
		RootDiskSizeGB int `json:"root_disk_size_gB"`
		NetworkRate    int `json:"network_rate"`
		DiskIOPS       int `json:"disk_iops"`
	} `json:"hardware"`
}

type KubernetesServiceOfferingsListResponse struct {
	Data []KubernetesServiceOffering `json:"data"`
}
