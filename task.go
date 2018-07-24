package main

import "fmt"

// Task is a task. this comment is just here to make gofmt happy
type Task struct {
	Name    string
	Span    uint
	DueTime string
}

// Prompt prompts the user for the information of Task t
func (t *Task) Prompt() {
	fmt.Print("\nWhat's the name of this task? ")
	fmt.Scanln(&t.Name)
	fmt.Print("How many days will it take to complete? (default is 7) ")

	if n, _ := fmt.Scanln(&t.Span); n == 0 {
		t.Span = 7
	}

	fmt.Print("When is this task due? (mm/dd/yyyy) ")
	fmt.Scanln(&t.DueTime)
}

// Display displays the Task t
func (t Task) Display() {
	fmt.Println(t.Name, t.Span, t.DueTime)
}
