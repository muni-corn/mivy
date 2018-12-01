package tasks

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Group holds the name of a group as a string
type Group string

// Task is a task. this comment is just here to make gofmt happy
type Task struct {
	Name            string
	Span            int
	GroupIndex      int
	UserDueDate     time.Time
	ModifiedDueDate time.Time
}

// Prompt prompts the user for the information of Task t
func (t *Task) Prompt(groups *[]Group) {
	t.PromptName()
	t.PromptDueDate()
	t.PromptSpan()
	t.PromptGroup(groups)
}

// PromptName prompts the user for Task t's name
func (t *Task) PromptName() {
	in := bufio.NewReader(os.Stdin)

prompt:
	fmt.Print("\nWhat's the name of this task? ")
	var err error
	t.Name, err = in.ReadString('\n')
	if err != nil {
		fmt.Println("Oops. There was an issue. Try again.")
		goto prompt
	}
	// trim whitespace and stuff
	t.Name = strings.TrimSpace(t.Name)
}

// PromptSpan prompts the user for Task t's span
func (t *Task) PromptSpan() {
	fmt.Print("How many days minimum will it take to complete? (default is 5): ")

	// n is the number of elements read in. if nothing is read (n == 0),
	// default the span to 5

prompt:
	if n, err := fmt.Scanln(&t.Span); err != nil {
		if n == 0 {
			t.Span = 5
		}
	} else {
		fmt.Println("There was a problem. Try that again: ")
		goto prompt
	}
}

// PromptGroup prompts for the Task's group
func (t *Task) PromptGroup(groups *[]Group) {
	// create a reader. a bufio.Reader lets us
	// get a line of input, even with spaces
	// included
	in := bufio.NewReader(os.Stdin)

	fmt.Println("Which group do you want to assign to this task?")
	for i := 0; i < len(*groups); i++ {
		fmt.Println("\t" + strconv.Itoa(i) + "\t" + string((*groups)[i]))
	}
	fmt.Println("\n\tn\tnew group")
	fmt.Println("\tu\tleave unassigned")
	fmt.Print("\n: ")

	// get the user input
promptSelection:
	// get the input for the number
	var input string
	_, err := fmt.Scan(&input)
	switch {
	// make sure the input is either a number or 'n'
	// if it's not (or if there was an error), go back to the prompt
	case err != nil && input[0] != 'n':
		fmt.Print("I didn't understand that. Try again: ")
		goto promptSelection

	case input[0] == 'u':
		fmt.Println("All right, we won't assign this task to a group.")
		t.GroupIndex = -1

	case input[0] == 'n':
		// if the input is valid, then check to see if it's 'n'
		// if it is, get the user's group name
		// if there was an issue, prompt again

		fmt.Print("Enter the group name: ")

	promptGroupName:
		groupName, err := in.ReadString('\n')

		if err != nil {
			fmt.Print("There was a problem. Try again: ")
			goto promptGroupName
		}

		*groups = append(*groups, Group(groupName))

		t.GroupIndex = len(*groups) - 1

	default:
		// otherwise, we're guessing this is a number
		// we'll make sure it's in range
		// if it's not, we'll prompt again
		// if it is, we'll pass it back up to the Task
		var groupIndex int
		_, err := fmt.Scan(&groupIndex)

		if err != nil {
			fmt.Print("There was an issue... Try again: ")
		} else if groupIndex < 0 || groupIndex >= len(*groups) {
			fmt.Print("That's not a number you can choose :) Try again: ")
			goto promptSelection
		} else {
			t.GroupIndex = groupIndex
		}
	}
}

// PromptDueDate prompts the user for Task t's due day
func (t *Task) PromptDueDate() {
	var dueInput string

prompt:
	fmt.Print("When is this task due? (m/d/yyyy) ")
	fmt.Scanln(&dueInput)

	var err error
	t.UserDueDate, err = time.ParseInLocation("1/2/2006", dueInput, time.Local)
	if err != nil {
		fmt.Println("Oops. There was an issue. Try again.")
		goto prompt
	}
	// log.Println("\t===>added task for day", t.UserDueDate)
	t.ModifiedDueDate = t.UserDueDate
}

// MarkDoneForToday changes the Span to mark a task
// as done
func (t *Task) MarkDoneForToday() {
	tomorrow := floorDate(time.Now().Add(time.Hour * 24))
	t.Span = int(floorDate(t.ModifiedDueDate).Sub(tomorrow) / (time.Hour * 24))
}

// IsDone returns true if the user span of Task t is less than or equal to zero.
func (t Task) IsDone() bool {
	return t.Span <= 0
}

// IsDateInRange returns true if the given timestamp
// is within range of the task's span
func (t Task) IsDateInRange() bool {
	now := time.Now()
	withinSpan := now.After(t.ModifiedDueDate.Add(time.Duration(-t.Span) * time.Hour * 24))
	return withinSpan
}

// GetDaysUntil returns the amount of days until
// the given date
func GetDaysUntil(date time.Time) int {
	untilDueDate := math.Ceil((float64(time.Until(date)) / float64(time.Hour*24)))

	return int(untilDueDate)
}

// Display displays the Task t, telling the user
// how much of it to do. ideally, this function should
// not be called if the current date is out of its
// range.
func (t Task) Display(groups []Group) {
	daysUntilModified := GetDaysUntil(t.ModifiedDueDate)
	daysUntilUser := GetDaysUntil(t.UserDueDate)

	var groupString string
	if t.GroupIndex >= 0 {
		groupString = string(groups[t.GroupIndex])
	}

	dueDateString := t.UserDueDate.Format("Monday, January 2")
	if t.IsDone() {
		return
	} else if t.IsDateInRange() {
		if daysUntilModified == 1 {
			fmt.Println("\tFinish", t.Name)
		} else if daysUntilModified < 1 {
			// only display as overdue if both the modified date and user
			// date are past
			if daysUntilUser < 1 {
				fmt.Println("\t(OVERDUE!) Finish", t.Name)
			} else {
				fmt.Println("\tFinish", t.Name)
			}
		} else if t.IsDateInRange() {
			fmt.Println("\tDo 1/"+strconv.Itoa(int(daysUntilModified)), "of", t.Name)
		} else {
			fmt.Println(t.Name, "isn't due for another", daysUntilModified, "days")
		}

		if daysUntilUser == 1 {
			fmt.Print("\t\t", groupString, ". Due tomorrow, ", dueDateString)
		} else if daysUntilUser > 0 {
			fmt.Print("\t\t", groupString, ". Due in ", daysUntilUser, " days on ", dueDateString)
		} else if daysUntilUser == 0 {
			fmt.Print("\t\t", groupString, ". Due today, ", dueDateString)
		} else if daysUntilUser == -1 {
			fmt.Print("\t\t", groupString, ". Due yesterday, on ", dueDateString)
		} else if daysUntilUser < 0 {
			fmt.Print("\t\t", groupString, ". Due ", -daysUntilUser, " days ago on ", dueDateString)
		}


		fmt.Print("\n\n")
	}
}

// Unmarshal parses a Task from its string form, usually from a line
// in a file.
//
// Data version 1 formats Task strings like so:
// 		T yyyymmdd span task name with many spaces ...
//		0 1		   2    3    4    5    6    7      ...
func Unmarshal(dataVersion int, task string) (t Task, err error) {
	switch dataVersion {
	case 1:
		// If the string we're given doesn't start with
		// T (the task tag), forget about it
		if !strings.HasPrefix(task, "T") {
			return
		}

		// split the string into separate pieces of data
		split := strings.Split(task, " ")

		t.UserDueDate, err = time.ParseInLocation("20060102", split[1], time.Local)

		if err != nil {
			panic(err)
		}

		t.Span, err = strconv.Atoi(split[2])

		if err != nil {
			panic(err)
		}

		for i := 3; i < len(split); i++ {
			t.Name += split[i] + " "
		}

		t.Name = strings.TrimSpace(t.Name)
	}
	return
}

// Marshal converts Task t into its text representation.
func (t Task) Marshal() string {
	dueDate := t.UserDueDate.Format("20060102")

	return "T " + dueDate + " " + strconv.Itoa(t.Span) + " " + t.Name
}

func floorDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
