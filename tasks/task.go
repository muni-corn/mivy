package tasks

import (
	"fmt"
	"time"
)

// Task is a task. this comment is just here to make gofmt happy
type Task struct {
	Name            string `json:"name"`
	UserSpan        uint   `json:"userSpan"`
	ModifiedSpan    uint   `json:"modifiedSpan"`
	UserDueTime     uint   `json:"userDueTime"`
	ModifiedDueTime uint   `json:"modifiedDueTime"`
	Done            bool   `json:"done"`
}

// Prompt prompts the user for the information of Task t
func (t *Task) Prompt() {
	fmt.Print("\nWhat's the name of this task? ")
	fmt.Scanln(&t.Name)
	fmt.Print("How many days minimum will it take to complete? (default is 3) ")

	// n is the number of elements read in. if nothing is read (n == 0),
	// default t.UserSpan to 3
	if n, _ := fmt.Scanln(&t.UserSpan); n == 0 {
		t.UserSpan = 3
	}

	var dueInput string
	fmt.Print("When is this task due? (mm/dd/yyyy) ")
	fmt.Scanln(&dueInput)

	var err error
	dueTime, err := time.ParseInLocation("01/02/2006", dueInput, time.UTC)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	t.UserDueTime = uint(dueTime.Unix())

	// default modifieds to the user's values
	t.ModifiedDueTime = t.UserDueTime
	t.ModifiedSpan = t.UserSpan
}

// GetUserDueTimeDay returns the day of this task in uint form
func (t Task) GetUserDueTimeDay() uint {
	// divide seconds by the amount of seconds per
	// day
	return uint(t.UserDueTime / (3600 * 24))
}

// Display - displays the Task t
func (t Task) Display() {
	fmt.Println(t.Name, t.UserSpan, t.UserDueTime)
}
