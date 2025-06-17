package geodbclient

type CountryDTO struct {
	CountryCode   string   `json:"code"`
	CurrencyCodes []string `json:"currencyCodes,omitempty"`
	CountryName   string   `json:"name"`
	WikiDataId    string   `json:"wikiDataId"`
}

// I wish there was a better way of doing this...
type GeoDBAPIResponse struct {
	Metadata struct {
		Count  int `json:"totalCount"`
		Offset int `json:"offset"`
	} `json:"metadata"`
	Countries []CountryDTO `json:"data"`
}
