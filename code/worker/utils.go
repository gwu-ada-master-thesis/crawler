package worker


import (
	// builtin
	"net/url"
	"strings"
	"log"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"

	//custom
	"masterthesis/utils"
)


func parseUrl(u url.URL) *Url {

	// if erro is not null or scheme is empty or host is empty, url is not valid
    if u.Scheme == "" || u.Host == "" {
        return nil
    }

	// set protocol
	protocol := u.Scheme

	// split host into: subdomain, domain name and extension
	hostSplitted := strings.Split(u.Host, ".")
	
	// if less than 2, means url is not valid
	if len(hostSplitted) < 2 {
		return nil
	}

	// set subdomaind
	subdomain := strings.Join(hostSplitted[: len(hostSplitted)-2], ".")

	// set domain name
	domainName 	:= hostSplitted[len(hostSplitted) - 2]

	// set extension
	extension 	:= hostSplitted[len(hostSplitted)  -1]

	// set path
	path := u.Path

	return &Url {
		Protocol: 	protocol,
		Subdomain: 	subdomain,
		DomainName: domainName,
		Extension: 	extension,
		Path: 		path,
		
		Metadata:	&UrlMetadata{},
	}

}


func dumpUrlsToJson(lines []Url, path string) {
    file, err := os.Create(path)
    if err != nil {
        log.Fatalf("Cannot create file %q: %s\n", path, err)
        return
    }
    defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// dump json to the standard output
	enc.Encode(lines)
}


func dumpDomainMetadataToJson(metadata DomainMetadata, path string) {
    file, err := os.Create(path)
    if err != nil {
        log.Fatalf("Cannot create file %q: %s\n", path, err)
        return
    }
    defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// dump json to the standard output
	enc.Encode(metadata)
}

func readDomainMetadataFromJson(path string) *DomainMetadata{
	
	// open file and handle error
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// close opened file when done
	defer file.Close()

	// read file as byte array
	lines, _ := ioutil.ReadAll(file)

	var metadata *DomainMetadata
	json.Unmarshal([]byte(lines), &metadata)

	return metadata
}


func generateFinalDestinationStructures(task string) string {

	// get output directory
	outPath := os.Getenv("OUTPUT")

	// set domain folder path
	taskDirPath := fmt.Sprintf("%s/%s", outPath, task)

	// create folder related if not exist
	utils.CreateDirectoryIfNotExist(taskDirPath)

	return taskDirPath
}