package tasks

import (
	"fmt"
	"time"
)

// Task is a task. this comment is just here to make gofmt happy
type Task struct {
	Name            string    `json:"name"`
	Span            uint      `json:"span"`
	UserDueTime     time.Time `json:"userDueTime"`
	ModifiedDueTime time.Time `json:"modifiedDueTime"`
	Done            bool      `json:"done"`
}

// Prompt prompts the user for the information of Task t
func (t *Task) Prompt() {
	fmt.Print("\nWhat's the name of this task? ")
	fmt.Scanln(&t.Name)
	fmt.Print("How many days will it take to complete? (default is 7) ")

	// n is the number of elements read in. if nothing is read (n == 0),
	// default t.Span to 7
	if n, _ := fmt.Scanln(&t.Span); n == 0 {
		t.Span = 7
	}

	var dueInput string
	fmt.Print("When is this task due? (mm/dd/yyyy) ")
	fmt.Scanln(&dueInput)

	var err error
	t.UserDueTime, err = time.ParseInLocation("01/02/2006", dueInput, time.Local)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
}

// Display - displays the Task t
func (t Task) Display() {
	fmt.Println(t.Name, t.Span, t.UserDueTime)
}
