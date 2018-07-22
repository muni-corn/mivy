package main

import "mivy/task"

func main() {
	t := task.Task{Name: "Name", Span: 7, DueTime: 1000000}
	t.Display()
}
