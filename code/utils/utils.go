package utils


import (
	// built-in
	"bufio"
    "fmt"
    "os"
    "log"
    "strconv"

    // custom

	// 3rd party
	"github.com/joho/godotenv"
)


func LoadEnvironment() {

  // load .env file
  err := godotenv.Load(".env")

  // check error
  if err != nil {
    log.Fatalf("Error loading .env file")
  }
}


// readLines reads a whole file into memory
func ReadLines(path string) ([]string, error) {
	
	// open file
    file, err := os.Open(path)

	// check for error
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}


// writeLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}


func CreateFile(path string) {
    file, err := os.Create(path)
    if err != nil {
        log.Fatalf("Cannot create file %q: %s\n", path, err)
        return
    }
    defer file.Close()
}


func CreateDirectoryIfNotExist(path string) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			log.Fatalf("Cannot folder file %q: %s\n", path, err)
		}
	}
}


func GenerateSocketName(customIdentifier int, serverName string) string {
	return fmt.Sprintf("/var/tmp/%s-%s-%s", strconv.Itoa(customIdentifier), serverName, strconv.Itoa(os.Getuid()))
}
