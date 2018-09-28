package main

import (
	"fmt"
	"mivy/data"
	"sort"
	"log"
	"time"
	t "mivy/tasks"
	"io/ioutil"
)

func main() {
	// Disable logging
	log.SetOutput(ioutil.Discard)

	tasks, groups := readData()
	fmt.Println("Welcome to Mivy!")
	displayHelp()
	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)
		switch input[0] {
		case 'a':
			addTask(&tasks, &groups)
			saveData(tasks, groups)
		case 'e':
			fmt.Println("TODO: edit a task")
			saveData(tasks, groups)
		case 'd':
			markATaskAsDone(&tasks)
			saveData(tasks, groups)
		case 'v':
			viewTasks(&tasks, &groups)
		case 's':
			fmt.Println("TODO: settings")
			saveData(tasks, groups)
		case 'q':
			quit(tasks, groups)
			return
		case 'h':
			displayHelp()
		}
	}
}

// sorts Tasks into an entire todo list, limiting the amount of tasks
// per day as much as possible
func sortTasks(tasks *[]t.Task) {
	log.Println("sorting tasks")

	// sorts the tasks by due date, or span, or name
	sort.Slice(*tasks, func(i, j int) bool {
		lhs, rhs := (*tasks)[i], (*tasks)[j]
		switch {
		case lhs.UserDueDate != rhs.UserDueDate:
			return time.Until(lhs.UserDueDate) < time.Until(rhs.UserDueDate)
		case lhs.Span != rhs.Span:
			return lhs.Span < rhs.Span
		default:
			return lhs.Name < rhs.Name
		}
	})

	// first, let's just change all of the modified due dates
	// to the user-specified due date
	for _, t := range *tasks {
		t.ModifiedDueDate = t.UserDueDate
	}

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
		// log.Println("\t===> Starting iteration on", task.Name)
		// log.Println("\t\t-> firstMatchIndex points to", (*tasks)[firstMatchIndex].Name)
		if t.GetDaysUntil(task.UserDueDate) == t.GetDaysUntil(taskToMatchForDate.UserDueDate) {
			log.Println(task.Name, "and", taskToMatchForDate.Name, "have matching due dates. continuing.")
			continue
		} else {
			log.Println(task.Name, "and", taskToMatchForDate.Name, "have different dates. sorting all preceding dates")

			numMatchingTasks := (firstMatchIndex - i)
			// log.Println("\t\t-> This amounts to", numMatchingTasks, "tasks, right?")

			daysBetweenDates := t.GetDaysUntil(taskToMatchForDate.UserDueDate) - t.GetDaysUntil(task.UserDueDate)

			log.Println("\tThere are", daysBetweenDates, "days between the two aforementioned tasks")

			// if numMatchingTasks > 1 {
			var x = 1 // used for multiplying daysBetweenDates / numMatchingTasks
			for j := i + 1; j <= firstMatchIndex; j++ {
				jTask := (*tasks)[j]
				// log.Println("\t\t-> looking at", jTask.Name, "due on day", jTask.UserDueDate)

				daysToAdd := (daysBetweenDates) * x / numMatchingTasks
				// log.Println("\t\t-> also, we're going to add", daysToAdd, "days out from the due date of", task.Name, "which is due on day", task.UserDueDate)

				log.Println("Looking at task", jTask.Name)

				durationToAdd := time.Duration(daysToAdd) * time.Hour * 24
				log.Println("\tAdding", durationToAdd)


				jTask.ModifiedDueDate = task.UserDueDate.Add(durationToAdd)
				// log.Println("\t\t-> so, the new due day of", jTask.Name, "is", jTask.ModifiedDueDate)
				x++

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

func addTask(tasks *[]t.Task, groups *[]t.Group) {
	var foo t.Task
	foo.Prompt(groups)
	*tasks = append(*tasks, foo)
	fmt.Println("Task added!")
	fmt.Println()
}

func editTask(tasks *[]t.Task) {
	task := (*tasks)[getATaskIndex(*tasks)]

	fmt.Println("What do you want to do with this task?\n", task.Name)
	fmt.Println()
}

func viewTasks(tasks *[]t.Task, groups *[]t.Group) {
	sortTasks(tasks)

	fmt.Print("\nHere's your todo list today:\n\n")

	for _, task := range *tasks {
		task.Display(*groups)
	}
	fmt.Println()
}

func markATaskAsDone(tasks *[]t.Task) {
	index := getATaskIndex(*tasks)
	(*tasks)[index].MarkDoneForToday()
	t := (*tasks)[index]
	fmt.Println("\"" + t.Name + "\" has been marked as done. Good job! ;)")
}

func getATaskIndex(tasks []t.Task) int {
	fmt.Print("\n")
	for i, task := range tasks {
		fmt.Println("\t", i, "\t"+task.Name)
	}
	fmt.Print("\nEnter the number of the task you want to change: ")
	var index int
	fmt.Scan(&index)

	return index
}

func readData() (tasks []t.Task, groups []t.Group) {
	var d data.Data
	data.ReadData(&d)
	tasks = d.Tasks
	groups = d.Groups

	return
}

// saveData creates a Data object with the given
// information and uses the mivy/data package
// to write the data to a file.
func saveData(tasks []t.Task, groups []t.Group) {
	// create the Data object and save it
	d := data.Data{Tasks: tasks, Groups: groups}
	data.WriteData(d)
}

func quit(tasks []t.Task, groups []t.Group) {
	saveData(tasks,groups)
	fmt.Println("Goodbye :)")
	fmt.Println()
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
