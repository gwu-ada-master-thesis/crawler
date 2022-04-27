package master


import (
	"sync"
)


type Master struct {
	Jobs   			map[string]*Job
	Access 			*sync.RWMutex
	ProcessDetails 	*ProcessDetails
	SocketName 		string
}


type ProcessDetails struct {
	latestJobId						int
	idleJobsCount 					int
	doneJobsCount 					int
	occupiedJobsCount 				int
	domainMetadataApiRequestCount 	int
}


type Job struct {
	Id  	int
	Task   	string
	Status 	string
}