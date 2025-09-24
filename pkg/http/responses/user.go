package responses

// UserProfileResponse represents the response for the user profile endpoint.
type UserProfileResponse struct {
	Data struct {
		ID           string      `json:"id"`
		Name         string      `json:"name"`
		Language     string      `json:"language"`
		NationalCode string      `json:"national_code"`
		Email        string      `json:"email"`
		Phone        string      `json:"phone"`
		Country      interface{} `json:"country"`
		State        interface{} `json:"state"`
		City         interface{} `json:"city"`
		Address      interface{} `json:"address"`
		Zip          interface{} `json:"zip"`
		Website      interface{} `json:"website"`
		Extra        struct {
			ReferralCode interface{} `json:"referral_code"`
		} `json:"extra"`
		Status             string      `json:"status"`
		Type               string      `json:"type"`
		CreatedAt          string      `json:"created_at"`
		UpdatedAt          string      `json:"updated_at"`
		CustomerZonesCount interface{} `json:"customer_zones_count"`
		InstancesCount     interface{} `json:"instances_count"`
		PaymentsCount      interface{} `json:"payments_count"`
		WalletsCount       interface{} `json:"wallets_count"`
		InviteCode         string      `json:"invite_code"`
		InvitedByMe        int         `json:"invited_by_me"`
		Picture            string      `json:"picture"`
	} `json:"data"`
}

// UserTokenAbilitiesResponse represents the response for the token abilities endpoint.
type UserTokenAbilitiesResponse struct {
	Abilities []string `json:"abilities"`
}

// UserSSHKey represents a single SSH key in the SSH key list response.
type UserSSHKey struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	DataKey     string `json:"datakey"`
	DataValue   string `json:"datavalue"`
	CreatedAt   string `json:"created_at"`
}

// UserSSHKeyListResponse represents the response for the SSH key list endpoint.
type UserSSHKeyListResponse struct {
	UserData []UserSSHKey `json:"userData"`
}

// AddUserSSHKeyResponse represents the response for adding a new SSH key.
type AddUserSSHKeyResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

// DeleteUserSSHKeyResponse represents the response for deleting an SSH key.
type DeleteUserSSHKeyResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

// WalletsBalanceResponse represents the response for the wallets balance endpoint.
type WalletsBalanceResponse struct {
	Data struct {
		Name           string  `json:"name"`
		Track          string  `json:"track"`
		Type           string  `json:"type"`
		Balance        float64 `json:"balance"`
		BalanceLimit   float64 `json:"balance_limit"`
		IsBlocked      bool    `json:"is_blocked"`
		MaxCost        float64 `json:"max_cost"`
		RemainingHours float64 `json:"remaining_hours"`
		UpdatedAt      string  `json:"updated_at"`
	} `json:"data"`
}

// CostDocument represents a single cost document in the yearly cost documents response.
type CostDocument struct {
	DateFrom                            string  `json:"dateFrom"`
	DateTo                              string  `json:"dateTo"`
	Instance                            float64 `json:"Instance"`
	NetworkNetflow                      float64 `json:"NetworkNetflow"`
	InstanceSnapshot                    float64 `json:"InstanceSnapshot"`
	InstanceDataVolumes                 float64 `json:"InstanceDataVolumes"`
	SupportOfferings                    float64 `json:"SupportOfferings"`
	NetworkInternetPublicAddressV4      float64 `json:"NetworkInternetPublicAddressV4"`
	NetworkDevice                       float64 `json:"NetworkDevice"`
	InstanceNetworkSecondaryIpAddressV4 float64 `json:"InstanceNetworkSecondaryIpAddressV4"`
	BucketSize                          float64 `json:"BucketSize"`
	BucketDownloadTraffic               float64 `json:"BucketDownloadTraffic"`
	BucketUploadTraffic                 float64 `json:"BucketUploadTraffic"`
	KubernetesNode                      float64 `json:"KubernetesNode"`
}

// CostDocumentsYearlyResponse represents the response for the yearly cost documents endpoint.
type CostDocumentsYearlyResponse struct {
	Data []CostDocument `json:"data"`
}

// PaymentListResponse represents the response for the payment list endpoint.
type PaymentListResponse struct {
	Data []interface{} `json:"data"`
}
