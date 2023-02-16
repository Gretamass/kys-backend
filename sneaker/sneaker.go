package sneaker

type Sneaker struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	Model               string `json:"model"`
	Brand               string `json:"brand"`
	Description         string `json:"description"`
	ProviderInformation map[string]struct {
		ProviderInformation
	} `json:"provider-information"`
	AvailabilityScrappers map[string]struct {
		AvailabilityScrappers
	} `json:"availability_scrappers"`
}

type ProviderInformation struct {
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

type AvailabilityScrappers struct {
	SearchFor string `json:"search_for"`
}
