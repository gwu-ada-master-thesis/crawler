package worker


import (
	// builtin
	"os"
	"log"
	"fmt"
	"sync"
	"time"
	"net/http"
	"encoding/json"
	
	// custom
	"masterthesis/utils"
	"masterthesis/master"

	// 3rd party
	"github.com/gocolly/colly/v2"
)


func requestJob(args *master.JobRequest, reply *master.JobResponse) bool {

	for {

		// request task
		master.Call("Master.RequestJob", args, reply)

		// check job status
		switch reply.Status {
			case "assigned":
				return true
			case "wait":
				continue
			default:
				return false

		}

	}
	
}


func reportJob(args *master.JobRequest, reply *master.JobResponse) {
	for {

		// update args
		args.Job = reply.Job

		// request task
		master.Call("Master.ReportJob", args, reply)

		// check job status
		switch reply.Status {
			case "success":
				return
			default:
				continue

		}

	}

}


func requestDomainMetadata(domain string) *DomainMetadata{
	url 	:= "http://ip-api.com/json/" + domain + "?fields=status,message,continent,continentCode,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query"
	
	client 	:= &http.Client{Timeout: 40 * time.Second}

	resp, err := client.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	metadata := DomainMetadata{}

	json.NewDecoder(resp.Body).Decode(&metadata)

	return &metadata
}


func (worker *Worker) crawl(job *master.Job) []Url {

	// extract task from job
	task := job.Task

	// initiate slice to store urls
	urls := make([]Url, 0)

	// initiate empty url type, that will be overwritten
	var url *Url

	// initiate url to use
	// instantiate collector
	c := colly.NewCollector(
		// run multiple threads
		// colly.Async(),

		// set maximum depth
		colly.MaxDepth(2),
		
		// cache responses to prevent multiple download of pages even if the collector is restarted
		colly.CacheDir(fmt.Sprintf("%s/%s/cache", os.Getenv("OUTPUT"), task)),
	)

	c.OnRequest(func(r *colly.Request) {

		// parse link into parts and if is a url store
		if url = parseUrl(*r.URL); url != nil {

			// print worker status
			log.Printf("Visiting %s\n", r.URL)
		
		}
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		
		// get href link
		link := e.Attr("href")

		// scrape all links
		e.Request.Visit(link)
	})

	c.OnHTML("html p", func(e *colly.HTMLElement) {

		if url != nil {
			
			// lock write access
			worker.Access.Lock()

			url.Metadata.Context = fmt.Sprintf("%s\n\r\t\n%s", url.Metadata.Context, e.Text)

			// release write access
			worker.Access.Unlock()

		}
	})

	c.OnHTML("html[lang]", func(e *colly.HTMLElement) {
		
		if url != nil {

			// lock write access
			worker.Access.Lock()

			url.Metadata.Language = e.Attr("lang")

			// store url 
			urls = append(urls, *url)

			// lock write access
			worker.Access.Unlock()

		}
	})

	// start scraping
	c.Visit("http://" + task)

	// wait till all threads are done
	// c.Wait()

	fmt.Println(len(urls))

	return urls
}


func (worker *Worker) Work() {

	for {

		// initiate request/response
		args 	:= master.JobRequest{}
		reply 	:= master.JobResponse{}

		// print worker status
		log.Printf("requesting a new job\n");time.Sleep(1 * time.Second)

		// if master is dead
		if !requestJob(&args, &reply) {
			// kill the process
			return
		}

		// extract job from reply
		job := reply.Job

		// print for visual
		fmt.Printf(`


		                                                                            
		                                TASK: %s
															                
		                                                                             
		                                                     
		`, job.Task); time.Sleep(1 * time.Second)


		// print empty line for visualization
		fmt.Println()

		// generate final destination directory path
		finalDirPath := generateFinalDestinationStructures(job.Task)

		// print worker status
		log.Printf("processing\n");time.Sleep(1 * time.Second)
		
		// set metadata file path
		metadataFilePath := fmt.Sprintf("%s/metadata.json", finalDirPath)

		// check if metadata has already been created 
		var metadata *DomainMetadata 
		if _, err := os.Stat(metadataFilePath); err != nil {

			// print worker status
			log.Printf("requesting domain metadata\n"); time.Sleep(1 * time.Second)

			// request domain metadata
			metadata = requestDomainMetadata(job.Task)

		} else {

			// read existed metadata
			metadata = readDomainMetadataFromJson(metadataFilePath)

		}

		var urls []Url
		if metadata.Status == "success" {

			// print worker status
			log.Printf("crawling\n"); time.Sleep(1 * time.Second)

			// start crawling
			urls = worker.crawl(job)

		} else {

			// print worker status
			log.Printf("terminating crawling: domain is not reachable\n"); time.Sleep(1 * time.Second)

		}

		// if crawled
		if len(urls) != 0 {

			// update status to success
			metadata.Scrapable = true
			
		}

		// print worker status
		log.Printf("storing domain matadata\n"); time.Sleep(1 * time.Second)

		// store metadata as json
		dumpDomainMetadataToJson(*metadata, metadataFilePath)

		// print worker status
		log.Printf("storing crawled data\n"); time.Sleep(1 * time.Second)

		// create and write to .json
		dumpUrlsToJson(urls, fmt.Sprintf("%s/urls.json", finalDirPath))

		// print worker status
		log.Printf("reporting\n"); time.Sleep(1 * time.Second)

		// update status of the job as done
		reportJob(&args, &reply); time.Sleep(1*time.Second)
	}

}


func initialize() {

	// load environment
	utils.LoadEnvironment()
}


func MakeWorker() *Worker{

	//
	// ----------------------------------------------------------------* variables
	//

	// initialize necessary dependencies
	initialize()

	// create worker
	worker := Worker{
		Access:	new(sync.RWMutex),
	}

	return &worker
}