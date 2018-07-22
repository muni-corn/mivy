package task

import "fmt"

type Task struct {
	Name    string
	Span    uint
	DueTime uint64
}

func (t Task) Display() {
	fmt.Println(t.Name, t.Span, t.DueTime);
}
