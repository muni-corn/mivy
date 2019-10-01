// Package actions contains all functions needed to start actions in mivy
package actions

import (
    "bufio"
    "fmt"
	"log"
    "sort"
    "time"

    "github.com/muni-corn/mivy/util"
    "github.com/muni-corn/mivy"
)

func ViewTasks(tasks mivy.TaskSlice) mivy.TaskSlice {
    log.Println("opening tasks")

    mesg := "Here's your list of tasks. Select one to edit it, or type a new task to add."
    in, out, err := util.OpenRofiWithMesg(mesg, "Mivy")
    if err != nil {
        util.RofiShowError(err)
        return tasks
    }
    defer out.Close()

    optimizeTasks(tasks)

    for _, task := range tasks {
        in.Write([]byte(task.DisplayString() + "\n"))
    }
    in.Close()

    bufout := bufio.NewReader(out)
    log.Println("buffered reader created; waiting for response from user")
    o, _, _ := bufout.ReadLine()
    if len(o) <= 0 {
        return tasks
    }
    log.Printf("got response: '%s'", string(o))

    // add or edit?
    task, index := getTaskFromDisplayString(string(o), tasks)
    if task == nil {
        log.Println("no existing task found, so adding")
        tasks = append(tasks, mivy.NewTask(string(o), tasks.Groups()))
    } else {
        log.Println("existing task found, editing")
        deleted := EditTask(task, tasks)
        if deleted {
            println("deleting at index", index, "while slice has length of", len(tasks))
            tasks = append(tasks[:index], tasks[index+1:]...)
        } else {
            tasks[index] = *task
        }
    }

    return tasks
}

// returns true if the task was deleted instead of edited
func EditTask(task *mivy.Task, tasks mivy.TaskSlice) (deleted bool) {
    log.Printf("editing task '%s'", task.Name)

    mesg := fmt.Sprintf("What do you want to do with '%s'?\n", task.Name)

    // display information about the task
    if task.Group != "" {
        mesg += fmt.Sprintf("\tin group %s\n", task.Group)
    }
    if !task.UserDueDate.IsZero() {
        mesg += fmt.Sprintf("\t%s\n", util.GetDueDateString(task.UserDueDate))
    }
    if task.URL != "" {
        mesg += fmt.Sprintf("\tlinks to %s\n", util.GetDueDateString(task.UserDueDate))
    }
    if time.Now().Before(task.SnoozedUntil) {
        mesg += fmt.Sprintf("\tsnoozed until %s\n", task.SnoozedUntil.Format(util.DueDateFormat))
    }

    // trim off trailing newline
    mesg = mesg[:len(mesg)-1]

    in, out, err := util.OpenRofiWithMesg(mesg, "Action")
    if err != nil {
        util.RofiShowError(err)
    }
    defer out.Close()

    // actions
    const (
        actionMarkComplete = "Mark it complete"
        actionMarkDoneToday = "Mark it done for today"
        actionVisitURL = "Visit URL"
        actionChangeName = "Change its name"
        actionChangeGroup = "Change its group"
        actionChangeDueDate = "Change its due date"
        actionChangeURL = "Change its URL"
        actionDelete = "Delete this task"
    )

    in.Write([]byte(actionMarkComplete + "\n"))
    in.Write([]byte(actionMarkDoneToday + "\n"))
    if task.URL != "" {
        in.Write([]byte(actionVisitURL + "\n"))
    }
    in.Write([]byte(actionChangeName + "\n"))
    in.Write([]byte(actionChangeGroup + "\n"))
    in.Write([]byte(actionChangeDueDate + "\n"))
    in.Write([]byte(actionChangeURL + "\n"))
    in.Write([]byte(actionDelete + "\n"))
    in.Close()

    bufout := bufio.NewReader(out)
    if bufout.Size() <= 0 {
        return false
    }
    o, _, _ := bufout.ReadLine()

    switch string(o) {
    case actionMarkDoneToday:
        task.MarkDoneForToday()
    case actionMarkComplete:
        task.MarkComplete()
    case actionVisitURL:
        task.VisitURL()
    case actionChangeName:
        task.PromptName()
    case actionChangeGroup:
        task.PromptGroup(tasks.Groups())
    case actionChangeDueDate:
        task.PromptDueDate()
    case actionChangeURL:
        task.PromptURL()
    case actionDelete:
        return true
    default:
        if string(o) != "" {
            task.Set(string(o))
        }
    }

    return false
}

// returns the task and the index
func getTaskFromDisplayString(displayLine string, tasks mivy.TaskSlice) (*mivy.Task, int) {
    log.Printf("getting pointer for display line '%s'", displayLine)
    for i, t := range tasks {
        if t.DisplayString() == displayLine {
            return &t, i
        }
    }
    return nil, -1
}

// sorts tasks into an entire todo list, limiting the amount of tasks
// per day as much as possible
func optimizeTasks(tasks mivy.TaskSlice) { // {{{
    if tasks == nil || len(tasks) <= 0 {
        return
    }

    log.Println("optimizing tasks")

    // sorts the tasks by due date, or span, or name
    log.Println("sorting...")
    sort.Sort(tasks)

    // first, let's just change all of the modified due dates
    // to the user-specified due date
    for _, t := range tasks {
        t.OptimalDueDate = t.UserDueDate
    }

    var firstMatchIndex = len(tasks) - 1
    for i := len(tasks) - 2; i >= 0; i-- {
        // step backwards and find how many tasks exist on the
        // same due date until we come across a task that is not
        // part of the same due date

        // tasks of the same due date are spaced out between one due date and
        // the other

        task := tasks[i]
        taskToMatchForDate := tasks[firstMatchIndex]
        log.Println("\t===> Starting iteration on", task.Name)
        log.Println("\t\t-> firstMatchIndex points to", tasks[firstMatchIndex].Name)
        if util.GetDaysUntil(task.UserDueDate) == util.GetDaysUntil(taskToMatchForDate.UserDueDate) {
            log.Println(task.Name, "and", taskToMatchForDate.Name, "have matching due dates. Continuing.")
            continue
        } else {
            log.Println(task.Name, "and", taskToMatchForDate.Name, "have different dates. Optimizing all preceding dates")

            numMatchingTasks := (firstMatchIndex - i)
            log.Println("\t\t-> This amounts to", numMatchingTasks, "tasks, right?")

            daysBetweenDates := util.GetDaysUntil(taskToMatchForDate.UserDueDate) - util.GetDaysUntil(task.UserDueDate)

            log.Println("\tThere are", daysBetweenDates, "days between the two aforementioned tasks")

            // if numMatchingTasks > 1 {
            var x = 1 // used for multiplying daysBetweenDates / numMatchingTasks
            for j := i + 1; j <= firstMatchIndex; j++ {
                jTask := tasks[j]
                log.Println("\t\t-> looking at", jTask.Name, "due on day", jTask.UserDueDate)

                daysToAdd := (daysBetweenDates) * x / numMatchingTasks
                log.Println("\t\t-> also, we're going to add", daysToAdd, "days out from the due date of", task.Name, "which is due on day", task.UserDueDate)

                log.Println("Looking at task", jTask.Name)

                durationToAdd := time.Duration(daysToAdd) * time.Hour * 24
                log.Println("\tAdding", durationToAdd)


                jTask.OptimalDueDate = task.UserDueDate.Add(durationToAdd)
                log.Println("\t\t-> so, the new due day of", jTask.Name, "is", jTask.OptimalDueDate)
                x++

                // so i guess we have to send this back to the slice
                tasks[j] = jTask
            }
            firstMatchIndex = i
        }
    }
} // }}}

// returns the new slice
func deleteTask(task mivy.Task, tasks []mivy.Task) []mivy.Task {
    log.Printf("Deleting task with name '%s'", task.Name)
    for i, t := range tasks {
        if t == task {
            return append(tasks[:i], tasks[i+1:]...)
        }
    }
    log.Print("Task to delete wasn't found?")
    return tasks
}

// vim: foldmethod=marker
