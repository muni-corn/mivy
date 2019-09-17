package data

import (
    "encoding/json"
	"log"
	"os"

    "github.com/muni-corn/mivy"
)

// Data is a collection of Tasks and settings
type Data struct {
    Tasks   []mivy.Task `json:"tasks"`
}

// Read returns a Data struct containing the groups and tasks stored in the
// user's .mivy file
func Read() Data {
    var d Data
	// get home directory file
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("error while getting home directory: " + err.Error())
		return d
	}

	// open the file
	file, err := os.Open(homePath + string(os.PathSeparator) + ".mivy")
    if os.IsNotExist(err) {
        return d
    }
	defer file.Close()

    // decompress (if needed) and decode
	// ungzip, err := gzip.NewReader(file)
    // if err != nil {
	// 	log.Println(err.Error())
    // }
    // defer ungzip.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&d)
    if err != nil {
		log.Println(err.Error())
    }
    return d
}

const dataVersion int = 1

// Write writes the Data to the user's .mivy file
func Write(d Data) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// create a file
	file, err := os.Create(homePath + string(os.PathSeparator) + ".mivy")
	defer file.Close()

    // encode into json first, then compress
    // gzipper := gzip.NewWriter(file)
    // defer gzipper.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    encoder.Encode(d)
}
