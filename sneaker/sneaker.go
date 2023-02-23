package sneaker

type Sneaker struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
	Brand string `json:"brand"`
	//Description         string `json:"description"`
	//ProviderInformation map[string]struct {
	//	ProviderInformation
	//} `json:"provider-information"`
	//AvailabilityScrappers map[string]struct {
	//	AvailabilityScrappers
	//} `json:"availability_scrappers"`
}

type SneakerInformation struct {
	Id                 int         `json:"id"`
	Name               string      `json:"name"`
	Model              string      `json:"model"`
	Brand              string      `json:"brand"`
	SneakerInformation SneakerInfo `json:"sneakerInformation"`
	//AvailabilityScrappers map[string]struct {
	//	AvailabilityScrappers
	//} `json:"availability_scrappers"`
}

type SneakerInfo struct {
	SneakerId      int    `json:"sneakerId"`
	MainInfo       string `json:"mainInfo"`
	MainImageUrl   string `json:"mainImageUrl"`
	AdditionalInfo string `json:"additionalInfo"`
}

type SneakerAvailability struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Model        string         `json:"model"`
	Brand        string         `json:"brand"`
	Availability []Availability `json:"availability"`
}

type Availability struct {
	Id         int     `json:"id"`
	ProductId  int     `json:"productId"`
	ProviderId int     `json:"providerId"`
	Available  bool    `json:"available"`
	Price      float32 `json:"price"`
}

type AvailabilityScrappers struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Model    string     `json:"model"`
	Brand    string     `json:"brand"`
	Scrapper []Scrapper `json:"scrapper"`
}

type ProviderInformation struct {
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

type Scrapper struct {
	Id         int    `json:"id"`
	ProductId  int    `json:"productId"`
	ProviderId int    `json:"providerId"`
	SearchFor  string `json:"search_for"`
}
