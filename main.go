package main

import (
	"fmt"
	"mivy/data"
	t "mivy/tasks"
	"sort"
)

func main() {
	tasks := readData()
	fmt.Println("Welcome to Mivy!")

	displayHelp()
	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)
		switch input[0] {
		case 'a':
			addTask(&tasks)
			saveData(tasks)
		case 'e':
			fmt.Println("TODO: edit a task")
			saveData(tasks)
		case 'd':
			fmt.Println("TODO: delete a task")
			saveData(tasks)
		case 'v':
			fmt.Println("TODO: view a todo list")
		case 's':
			fmt.Println("TODO: settings")
			saveData(tasks)
		case 'q':
			quit(tasks)
			return
		case 'h':
			displayHelp()
		}
	}
}

// sorts Tasks into an entire todo list, limiting the amount of tasks
// per day as much as possible
func sortTasks(tasks *[]t.Task) {
	// sorts the tasks by due date, or span, or name
	sort.Slice(tasks, func(i, j int) bool {
		lhs, rhs := (*tasks)[i], (*tasks)[j]
		switch {
		case lhs.GetUserDueTimeDay() != rhs.GetUserDueTimeDay():
			return lhs.GetUserDueTimeDay() < rhs.GetUserDueTimeDay()
		case lhs.UserSpan != rhs.UserSpan:
			return lhs.UserSpan < rhs.UserSpan
		default:
			return lhs.Name < rhs.Name
		}
	})

	// the number of dates that match
	var firstMatchIndex = len(*tasks) - 1
	var taskToMatchForDate = (*tasks)[firstMatchIndex]
	for i := len(*tasks) - 2; i >= 0; i-- {
		// step backwards and find how many tasks exist on the
		// same due date until we come across a task that is not
		// part of the same due date
		//
		// the non-matching due date will be used to space out the
		// previous due dates over time

		task := (*tasks)[i]
		if task.GetUserDueTimeDay() == taskToMatchForDate.GetUserDueTimeDay() {
			continue
		} else {
			var numMatchingTasks = uint(firstMatchIndex - i)
			var daysBetweenDates = task.GetUserDueTimeDay() - taskToMatchForDate.GetUserDueTimeDay()
			var x uint = 1 // used for multiplying daysBetweenDates / numMatchingTasks
			for j := i + 1; j <= firstMatchIndex; j++ {
				jTask := (*tasks)[j]

				newSpan := daysBetweenDates / numMatchingTasks
				jTask.ModifiedDueTime = task.UserDueTime + daysBetweenDates*x/numMatchingTasks
				x++
				if newSpan > jTask.UserSpan {
					jTask.ModifiedSpan = newSpan
				} else {
					jTask.ModifiedSpan = jTask.UserSpan
				}
			}
		}
	}
}

func addTask(tasks *[]t.Task) {
	var foo t.Task
	foo.Prompt()
	*tasks = append(*tasks, foo)
	fmt.Println("Task added! Here are all the tasks you have now:")
	viewTasks(*tasks)
	fmt.Println()
}

func viewTasks(tasks []t.Task) {
	for _, task := range tasks {
		task.Display()
	}
}

func readData() (tasks []t.Task) {
	var d data.Data
	data.ReadData(&d)
	tasks = d.Tasks

	return
}

// saveData creates a Data object with the given
// information and uses the mivy/data package
// to write the data to a file.
func saveData(tasks []t.Task) {
	fmt.Println("Saving data...")

	// create the Data object and save it
	d := data.Data{Tasks: tasks}
	data.WriteData(d)

	fmt.Println("Saved!")
	fmt.Println()
}

func quit(tasks []t.Task) {
	saveData(tasks)
	fmt.Println("Goodbye :)")
}

func displayHelp() {
	fmt.Println("Here's what you can do:")
	fmt.Println("\ta\tadd a task")
	fmt.Println("\te\tedit a task")
	fmt.Println("\td\tdelete a task")
	fmt.Println("\tv\tview a todo list")
	fmt.Println("\ts\tsettings")
	fmt.Println("\tq\tquit")
	fmt.Println("\th\thelp")
}
