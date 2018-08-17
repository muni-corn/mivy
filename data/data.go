package data

import (
	"bufio"
	j "encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	t "mivy/tasks"
	"os"
	"strconv"
)

// Data is a collection of Tasks and settings
type Data struct {
	Version uint     `json:"dataVersion"`
	Tasks   []t.Task `json:"tasks"`
}

// ReadData returns a Data struct through pass-by-reference
func ReadData(d *Data) {
	// get home directory file
	homePath, err := homedir.Dir()
	if err != nil {
		fmt.Println("error while getting home directory: " + err.Error())
		return
	}

	// open the file
	f, err := os.Open(homePath + string(os.PathSeparator) + ".mivy")
	defer f.Close()

	// create a reader
	reader := bufio.NewReader(f)
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("error while opening reader: " + err.Error())
		return
	}

	// get the version number from the first line of the file
	version, err := strconv.Atoi(string(line))
	fmt.Println("version: " + strconv.Itoa(version))
	if err != nil {
		fmt.Println("error while converting version to string: ", err.Error())
		return
	}

	// read the rest of the file into a buffer
	buffer := make([]byte, reader.Buffered())
	reader.Read(buffer)

	// decode the json
	if err := j.Unmarshal(buffer, &d); err != nil {
		fmt.Println("there was an error in parsing the json: ", err)
		return
	}

	switch version {
	case 1:
		fmt.Println(d.Tasks)
	}
}

const dataVersion int = 1

// WriteData writes the user's data to a file
func WriteData(d Data) {
	homePath, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return
	}

	// ENCODE TO JSON
	json, err := j.Marshal(d)
	if err != nil {
		fmt.Println(err)
		return
	}

	// WRITE TO FILE
	// create a file
	f, err := os.Create(homePath + string(os.PathSeparator) + ".mivy")
	defer f.Close()

	// open a writing stream
	writer := bufio.NewWriter(f)
	writer.WriteString(strconv.Itoa(dataVersion) + "\n")
	_, err2 := writer.WriteString(string(json))
	if err2 != nil {
		fmt.Println(err2)
	}
	writer.Flush()
}
