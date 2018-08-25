package tasks

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task is a task. this comment is just here to make gofmt happy
type Task struct {
	Name           string `json:"name"`
	UserSpan       int    `json:"userSpan"`
	ModifiedSpan   int    `json:"modifiedSpan"`
	UserDueDay     int    `json:"userDueTime"`
	ModifiedDueDay int    `json:"modifiedDueTime"`
	Done           bool   `json:"done"`
}

// Prompt prompts the user for the information of Task t
func (t *Task) Prompt() {

	// PROMPT FOR NAME
	in := bufio.NewReader(os.Stdin)
	fmt.Print("\nWhat's the name of this task? ")
	var err error
	t.Name, err = in.ReadString('\n')
	if err != nil {
		panic(err) // we shouldn't do this but oh well lol
	}
	// trim whitespace and stuff
	t.Name = strings.TrimSpace(t.Name)

	//////////////////////////

	// PROMPT FOR DUE DAY
	var dueInput string
	fmt.Print("When is this task due? (mm/dd/yyyy) ")
	fmt.Scanln(&dueInput)

	dueTime, err := time.ParseInLocation("01/02/2006", dueInput, time.Local)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}
	t.UserDueDay = int(dueTime.Unix()/(3600*24)) + 1 // for some reason we just have to add one here
	fmt.Println("\t===>added task for day", t.UserDueDay)
	///////////////////////////

	// PROMPT FOR SPAN
	fmt.Print("How many days minimum will it take to complete? (default is 3) ")

	// n is the number of elements read in. if nothing is read (n == 0),
	// default t.UserSpan to 3
	if n, _ := fmt.Scanln(&t.UserSpan); n == 0 {
		t.UserSpan = 3
	}
	///////////////////////////

	// default modifieds to the user's values
	t.ModifiedDueDay = t.UserDueDay
	t.ModifiedSpan = t.UserSpan
}

// IsDateInRange returns true if the given timestamp
// is within range of the task's span
func (t Task) IsDateInRange(currentDay int) bool {
	return currentDay >= t.ModifiedDueDay-t.ModifiedSpan
}

// GetDaysUntilDue returns the amount of days until
// this task is due.
func (t Task) GetDaysUntilDue(currentDay int) int {
	return int(t.ModifiedDueDay) - int(currentDay)
}

// Display displays the Task t, telling the user
// how much of it to do. ideally, this function should
// not be called if the current date is out of its
// range.
func (t Task) Display(currentDay int) {
	daysUntilDue := t.GetDaysUntilDue(currentDay)
	if daysUntilDue == 1 {
		fmt.Println("Finish", t.Name)
	} else if daysUntilDue < 1 {
		fmt.Println("Finish (OVERDUE!)", t.Name)
	} else if t.IsDateInRange(currentDay) {
		fmt.Println("Do 1/"+strconv.Itoa(daysUntilDue), "of", t.Name)
	} else {
		fmt.Println(t.Name, "isn't due for another", daysUntilDue, "days")
	}
}
