package main

import (
	"fmt"
	t "mivy/tasks"
)

func main() {
	var tasks []t.Task
	fmt.Println("Welcome to Mivy!")

	displayHelp()
	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)
		switch input[0] {
		case 'a':
			addTask(&tasks)
		case 'e':
			fmt.Println("TODO: edit a task")
		case 'd':
			fmt.Println("TODO: delete a task")
		case 'v':
			fmt.Println("TODO: view a todo list")
		case 's':
			fmt.Println("TODO: settings")
		case 'q':
			quit()
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

func quit() {
	fmt.Println("TODO: save work and quit")
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
