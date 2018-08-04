package main

import (
	"fmt"
)

func main() {
	var tasks []Task
	fmt.Println("Welcome to Mivy!")
	addTask(&tasks)
	tasks[0].Display()
}

func addTask(dest *[]Task) {
	var foo Task
	foo.Prompt()
	*dest = append(*dest, foo)
}

func displayHelp() {
	fmt.Println("Here's what you can do:");
	fmt.Println("\ta\tadd a task");
	fmt.Println("\te\tedit a task");
	fmt.Println("\td\tdelete a task");
	fmt.Println("\tv\tview a todo list");
	fmt.Println("\ts\tsettings");
	fmt.Println("\th\thelp");
}
