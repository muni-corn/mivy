package main

import (
	"fmt"
	"mivy/data"
	t "mivy/tasks"
	"sort"
	"time"
)

func main() {
	tasks := readData()
	fmt.Println("Welcome to Mivy!")
	fmt.Println("Today is day", time.Now().Unix()/(3600*24))

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
			markATaskAsDone(&tasks);
			saveData(tasks)
		case 'v':
			viewTasks(&tasks)
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
	sort.Slice(*tasks, func(i, j int) bool {
		lhs, rhs := (*tasks)[i], (*tasks)[j]
		switch {
		case lhs.UserDueDay != rhs.UserDueDay:
			return lhs.UserDueDay < rhs.UserDueDay
		case lhs.UserSpan != rhs.UserSpan:
			return lhs.UserSpan < rhs.UserSpan
		default:
			return lhs.Name < rhs.Name
		}
	})

	var firstMatchIndex = len(*tasks) - 1
	for i := len(*tasks) - 2; i >= 0; i-- {
		// step backwards and find how many tasks exist on the
		// same due date until we come across a task that is not
		// part of the same due date
		//
		// the non-matching due date will be used to space out the
		// previous due dates over time

		task := (*tasks)[i]
		taskToMatchForDate := (*tasks)[firstMatchIndex]
		// fmt.Println("\t===> Starting iteration on", task.Name)
		// fmt.Println("\t\t-> firstMatchIndex points to", (*tasks)[firstMatchIndex].Name)
		if task.UserDueDay == taskToMatchForDate.UserDueDay {
			// fmt.Println("\t===>", task.Name, "and", taskToMatchForDate.Name, "have matching due dates. continuing.")
			continue
		} else {
			// fmt.Println("\t===>", task.Name, "and", taskToMatchForDate.Name, "have different dates. sorting all preceding dates")

			numMatchingTasks := (firstMatchIndex - i)
			// fmt.Println("\t\t-> This amounts to", numMatchingTasks, "tasks, right?")

			daysBetweenDates := taskToMatchForDate.UserDueDay - task.UserDueDay
			// fmt.Println("\t\t-> There are", daysBetweenDates, "days between the two aforementioned tasks")

			newSpan := daysBetweenDates / numMatchingTasks
			// fmt.Println("\t\t-> we'll try to set the new span of these tasks to", newSpan)

			// if numMatchingTasks > 1 {
				var x = 1 // used for multiplying daysBetweenDates / numMatchingTasks
				for j := i + 1; j <= firstMatchIndex; j++ {
					jTask := (*tasks)[j]
					// fmt.Println("\t\t-> looking at", jTask.Name, "due on day", jTask.UserDueDay)

					daysToAdd := (daysBetweenDates) * x / numMatchingTasks
					// fmt.Println("\t\t-> also, we're going to add", daysToAdd, "days out from the due date of", task.Name, "which is due on day", task.UserDueDay)

					jTask.ModifiedDueDay = task.UserDueDay + daysToAdd
					// fmt.Println("\t\t-> so, the new due day of", jTask.Name, "is", jTask.ModifiedDueDay)
					x++
					if newSpan > jTask.UserSpan {
						jTask.ModifiedSpan = newSpan
					} else {
						jTask.ModifiedSpan = jTask.UserSpan
					}

					// so i guess we have to send this back to the slice
					(*tasks)[j] = jTask
				}
			// } else {
			// 	if newSpan > (*tasks)[i].UserSpan {
			// 		(*tasks)[i+1].ModifiedSpan = newSpan
			// 	} else {
			// 		(*tasks)[i+1].ModifiedSpan = (*tasks)[i].UserSpan
			// 	}
			// }
			firstMatchIndex = i
		}
	}
}

func addTask(tasks *[]t.Task) {
	var foo t.Task
	foo.Prompt()
	*tasks = append(*tasks, foo)
	fmt.Println("Task added! Here are all the tasks you have now:")
	viewTasks(tasks)
	fmt.Println()
}

func viewTasks(tasks *[]t.Task) {
	sortTasks(tasks)

	fmt.Print("\nHere's your todo list today:\n\n")

	currentDay := int(time.Now().Unix() / (3600 * 24))
	for _, task := range *tasks {
		task.Display(currentDay)
	}
	fmt.Println()
}

func markATaskAsDone(tasks *[]t.Task) {
	index := getATaskIndex(*tasks)
	(*tasks)[index].Done = true;
	t := (*tasks)[index]
	fmt.Println("\"" + t.Name + "\" has been marked as done. Good job! ;)")
}

func getATaskIndex(tasks []t.Task) int {
	fmt.Print("\n");
	for i, task := range tasks {
		fmt.Println("\t", i, "\t" + task.Name)
	}
	fmt.Print("\nEnter the number of the task you want to change: ")
	var index int
	fmt.Scan(&index)

	return index
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
	fmt.Println("\td\tmark a task as done")
	fmt.Println("\tv\tview today's todo list")
	fmt.Println("\ts\tsettings")
	fmt.Println("\tq\tquit")
	fmt.Println("\th\thelp")
}
