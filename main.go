package main

import (
	"fmt"
)

func main() {
	var tasks []Task
	fmt.Println("Welcome to Mivy!")
	var foo Task
	foo.Prompt()
	tasks = append(tasks, foo)
	tasks[0].Display();
}
