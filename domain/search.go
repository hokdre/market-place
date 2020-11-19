package domain

type Search struct {
	Tags      []string   `json:"tags"`
	Products  []Product  `json:"products"`
	Merchants []Merchant `json:"merchants"`
}

type CategorySearch struct {
	Buckets []struct {
		Count          int    `json:"doc_count"`
		Key            string `json:"key"`
		SecondCategory struct {
			Buckets []struct {
				Count         int    `json:"doc_count"`
				Key           string `json:"key"`
				ThirdCategory struct {
					Buckets []struct {
						Count int    `json:"doc_count"`
						Key   string `json:"key"`
					} `json:"buckets"`
				} `json:"third_category"`
			} `json:"buckets"`
		} `json:"second_category"`
	} `json:"buckets"`
}

type CitySearch struct {
	Buckets []struct {
		Count int    `json:"doc_count"`
		Key   string `json:"key"`
	} `json:"buckets"`
}

type SearchProduct struct {
	Categories CategorySearch `json:"categories"`
	Cities     CitySearch     `json:"cities"`
	Products   []Product      `json:"products"`
}
