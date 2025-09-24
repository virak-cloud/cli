package responses

type ObjectStorageBucket struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Status    string `json:"status"`
	Policy    string `json:"policy"`
	Size      int    `json:"size"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	Tier      string `json:"tier"`
	IsFailed  bool   `json:"is_failed"`
	Message   string `json:"message"`
}

type ObjectStorageBucketsResponse struct {
	Data []ObjectStorageBucket `json:"data"`
}

type ObjectStorageBucketResponse struct {
	Data ObjectStorageBucket `json:"data"`
}

type ObjectStorageBucketCreationResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type ObjectStorageEvent struct {
	ProductModel  string `json:"product_model"`
	ProductID     string `json:"product_id"`
	ProductSource string `json:"product_source"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	CreatedAt     int    `json:"created_at"`
}

type ObjectStorageEventsResponse struct {
	Data []ObjectStorageEvent `json:"data"`
	Meta Meta                 `json:"meta"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	From        int `json:"from"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	To          int `json:"to"`
	Total       int `json:"total"`
}
