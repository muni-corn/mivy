package main

import (
	"log"
    "os"

    "github.com/muni-corn/mivy/actions"
    "github.com/muni-corn/mivy/data"
)

func main() {
	// logging
    home, err := os.UserHomeDir()
    if err == nil {
        logFile, _ := os.Create(home + "/.mivy.log")
        log.SetOutput(logFile)
    } else {
        println("couldn't get home dir:", err.Error())
    }

    d := data.Read()

    d.Tasks = actions.ViewTasks(d.Tasks)

    data.Write(d)
}
