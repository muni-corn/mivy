package data

import (
	"bufio"
	"log"
	"github.com/mitchellh/go-homedir"
	"os"
	"strconv"
	"strings"
	t "mivy/tasks"
)

// Data is a collection of Tasks and settings
type Data struct {
	Version uint     
	Groups  []t.Group
	Tasks   []t.Task
}

// ReadData returns a Data struct through pass-by-reference
func ReadData(d *Data) {
	// get home directory file
	homePath, err := homedir.Dir()
	if err != nil {
		log.Println("error while getting home directory: " + err.Error())
		return
	}

	// open the file
	file, err := os.Open(homePath + string(os.PathSeparator) + ".mivy")
	defer file.Close()

	// create a reader
	scanner := bufio.NewScanner(file)

	// get the version number from the first line of the file
	var line string
	if scanner.Scan() {
		line = scanner.Text()
	} else {
		return
	}

	version, err := strconv.Atoi(string(line))
	if err != nil {
		log.Println("error while converting version to string: ", err.Error())
		return
	}

	groupIndex := -1

	switch version {
	case 1:
		// iterate through the file. lines that start with G are a group name. all
		// tasks under a G line are a part of that group. Task lines start with T.
		// anything else should be ignored.

		// scan each line in the file. i'm guessing scanner.Scan() returns false if
		// we've reached the end of the file.
		for scanner.Scan() {
			line := scanner.Text() // gets the actual line

			switch {
			case strings.HasPrefix(line, "G"):
				split := strings.Split(line, " ")

				groupName := "" // create an empty string from the group name

				// the group name takes up the rest of the line starting
				// at index 1
				for i := 1; i < len(split); i++ {
					groupName += split[i] + " "
				}
				groupName = strings.TrimSpace(groupName)

				d.Groups = append(d.Groups, t.Group(groupName))

				groupIndex++
			case strings.HasPrefix(line, "T"):
				task, err := t.Unmarshal(version, line)
				if err != nil {
					panic(err)
				}
				task.GroupIndex = groupIndex
				d.Tasks = append(d.Tasks, task)
			}
		}
	}
}

const dataVersion int = 1

// WriteData writes the user's data to a file
func WriteData(d Data) {
	homePath, err := homedir.Dir()
	if err != nil {
		log.Println(err)
		return
	}

	// create a file
	f, err := os.Create(homePath + string(os.PathSeparator) + ".mivy")
	defer f.Close()

	// open a writing stream
	writer := bufio.NewWriter(f)
	writer.WriteString(strconv.Itoa(dataVersion) + "\n")

	// first, write any tasks that are unassigned

	for i := 0; i < len(d.Groups); i++ {
		// marshal the group name
		_, err2 := writer.WriteString("G " + string(d.Groups[i]) + "\n")

		// write tasks, but only if they match the current index
		for _, t := range d.Tasks {
			if t.GroupIndex == i {
				writer.WriteString(t.Marshal() + "\n")
			}
		}

		if err2 != nil {
			log.Println(err2)
		}
	}

	writer.Flush()
}
