package worker


import (
	// builtin
	"sync"

	// custom
	"masterthesis/master"
)


type Worker struct {
	Access *sync.RWMutex
}


type Url struct {
	Protocol      	string
	Subdomain     	string
	DomainName    	string
	Extension 		string
	Path			string
	Metadata		*UrlMetadata

}

type UrlMetadata struct {
	Context		string
	Language	string
}

type Domain struct {
	Name		string
	Extension	string
	Metadata	*DomainMetadata
}


type DomainMetadata struct {
	IpAddress   	string    	`json:"query"`
	Status 			string    	`json:"status"`
	Message 		string    	`json:"message"`
	Continent 		string    	`json:"continent"`
	ContinentCode 	string 		`json:"continentCode"`
	Country 		string    	`json:"country"`
	Region 			string		`json:"region"`
	RegionName 		string		`json:"regionName"`
	City 			string		`json:"city"`
	Zip 			string		`json:"zip"`
	Lat 			string		`json:"lat"`
	Lon 			string		`json:"lon"`
	Timezone 		string		`json:"timezone"`
	Isp 			string		`json:"isp"`
	As 				string		`json:"as"`
	AsName 			string		`json:"asname"`
	Scrapable		bool		`json:"scrapable"`
}

type JobRequest struct {
}

type JobResponse struct {
	Job     *master.Job
	Status 	string
}