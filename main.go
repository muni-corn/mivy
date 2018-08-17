package main

import (
	"fmt"
	t "mivy/tasks"
	"mivy/data"
)

func main() {
	tasks := readData();
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
			return;
		case 'h':
			displayHelp()
		}
	}
}

func addTask(tasks *[]t.Task) {
	var foo t.Task
	foo.Prompt()
	*tasks = append(*tasks, foo)
	fmt.Println("Task added! Here are all the tasks you have now:")
	viewTasks(*tasks)
	fmt.Println();
}

func viewTasks(tasks []t.Task) {
	for _, task := range tasks {
		task.Display()
	}
}


func readData() (tasks []t.Task) {
	var d data.Data;
	data.ReadData(&d);
	tasks = d.Tasks;
	
	return
}

// saveData creates a Data object with the given
// information and uses the mivy/data package
// to write the data to a file.
func saveData(tasks []t.Task) {
	fmt.Println("Saving data...")

	// create the Data object and save it
	d := data.Data{Tasks: tasks}
	data.WriteData(d);
	
	fmt.Println("Saved!")
	fmt.Println()
}

func quit(tasks []t.Task) {
	saveData(tasks);
	fmt.Println("Goodbye :)");
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
