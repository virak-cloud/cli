package responses

// DataCenter represents the response for the zone list endpoint.
type DataCenter struct {
	Data []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Active   bool   `json:"active"`
	} `json:"data"`
}

// ZoneActiveServicesResponse represents the response for the zone active services endpoint.
type ZoneActiveServicesResponse struct {
	Instance      bool `json:"Instance"`
	DataVolume    bool `json:"DataVolume"`
	Network       bool `json:"Network"`
	ObjectStorage bool `json:"ObjectStorage"`
	K8s           bool `json:"K8s"`
}

// CustomerResourceResponse represents the response for the customer resource endpoint.
type CustomerResourceResponse struct {
	InstanceResourceCollected struct {
		Memory struct {
			Collected int `json:"collected"`
			Total     int `json:"total"`
		} `json:"memory"`
		CPUNumber struct {
			Collected int `json:"collected"`
			Total     int `json:"total"`
		} `json:"cpunumber"`
		DataVolume struct {
			Collected int `json:"collected"`
			Total     int `json:"total"`
		} `json:"datavolume"`
		VMLimit struct {
			Collected int `json:"collected"`
			Total     int `json:"total"`
		} `json:"vmlimit"`
	} `json:"instance_resource_collected"`
}

// ZoneNetworksResponse represents the response for the zone networks endpoint.
type ZoneNetworksResponse struct {
	Data []struct {
		ID              string        `json:"id"`
		Name            string        `json:"name"`
		Status          string        `json:"status"`
		IPConfig        []interface{} `json:"ip_config"`
		IPConfigV6      []interface{} `json:"ip_config_v6"`
		NetworkOffering struct {
			ID                     string  `json:"id"`
			Name                   string  `json:"name"`
			DisplayName            string  `json:"displayname"`
			DisplayNameFa          string  `json:"displayname_fa"`
			HourlyStartedPrice     float64 `json:"hourly_started_price"`
			TrafficPricePerGig     float64 `json:"traffic_price_per_gig"`
			TrafficTransferFreeGig float64 `json:"traffic_transfer_free_gig"`
			NetworkRate            float64 `json:"networkrate"`
			Type                   string  `json:"type"`
			Description            string  `json:"description"`
			InternetProtocol       string  `json:"internet_protocol"`
		} `json:"network_offering"`
		InstanceNetwork []interface{} `json:"instance_network"`
	} `json:"data"`
}
