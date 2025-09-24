package responses

type DomainList struct {
	Data []Domain `json:"data"`
}

type Domain struct {
	Domain  string `json:"domain"`
	Status  string `json:"status"`
	DNSInfo struct {
		VirakDNS  []*string `json:"virak_dns"`
		DomainDNS []string  `json:"domain_dns"`
	} `json:"dns_info"`
}

type DomainShow struct {
	Data Domain `json:"data"`
}

type DnsMessage struct {
	Message string `json:"message"`
}

type Record struct {
	Name        string    `json:"name"`
	TTL         int       `json:"ttl"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	IsProtected bool      `json:"is_protected"`
	Content     []Content `json:"content"`
}

type Content struct {
	ID         string `json:"id"`
	ContentRaw string `json:"content_raw"`
}

type RecordList struct {
	Data []Record `json:"data"`
}

type DNSEvent struct {
	ProductModel  string `json:"product_model"`
	ProductID     string `json:"product_id"`
	ProductSource string `json:"product_source"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	CreatedAt     int    `json:"created_at"`
}

type DNSEventsResponse struct {
	Data []DNSEvent `json:"data"`
	Meta Meta       `json:"meta"`
}
