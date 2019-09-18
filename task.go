package mivy

import (
    "bufio"
	"fmt"
	"strings"
	"time"

    "github.com/muni-corn/mivy/util"
)

// Task is a task. this comment is just here to make gofmt happy
type Task struct {
    Name           string    `json:"name"`
	Group          string    `json:"group"`
    Complete       bool      `json:"complete"`
	UserDueDate    time.Time `json:"userDueDate"`
	OptimalDueDate time.Time `json:"-"`
    SnoozedUntil   time.Time `json:"snoozedUntil"`
}

// NewTask returns a new Task, parsed from a string
func NewTask(from string, groupSuggestions []string) Task {
    t := Task{}

    t.Set(from)

    if t.UserDueDate.IsZero() {
        t.PromptDueDate()
    }
    if t.Group == "" {
        t.PromptGroup(groupSuggestions)
    }

    return t
}

// parses metadata tags and records them in the task
func parseMetadataTags(tokens []string, task *Task) {
    for i, t := range tokens {
        switch t {
        case "@snooze", "@snoozed", "@start", "@open":
            var err error
            dueDateStr := getStringToNextTag(tokens[i+1:])
            task.SnoozedUntil, err = time.ParseInLocation(util.DueDateFormat, dueDateStr, time.Local)
            if err != nil {
                task.SnoozedUntil = time.Time{}
            }
        case "@due", "@date", "@duedate":
            var err error
            dueDateStr := getStringToNextTag(tokens[i+1:])
            task.UserDueDate, err = time.ParseInLocation(util.DueDateFormat, dueDateStr, time.Local)
            if err != nil {
                task.UserDueDate = time.Time{}
            }
        case "@group", "@class":
            task.Group = getStringToNextTag(tokens[i+1:])
        default:
            continue
        }
        parseMetadataTags(tokens[i+1:], task)
        return
    }
}

// returns true if the string is actually a metadata tag; starts with @
func isMetadataTag(token string) bool {
    if len(token) <= 0 {
        return false
    }
    return token[0] == '@'
}

// iterates through the slice and returns
func getStringToNextTag(tokens []string) string {
    for i, s := range tokens {
        if isMetadataTag(s) {
            // found a metadata tag
            return strings.Join(tokens[:i], " ")
        }
    }
    return strings.Join(tokens, " ")
}

func (t *Task) PromptName() {
    val, err := getValueFromRofi("What's this task's title?", "Title:")
    if err != nil {
        util.RofiShowError(err)
    }

    t.Name = val
}

func (t *Task) PromptGroup(suggestions []string) {
    val, err := getValueFromRofi("Which group should this task be a part of?", "Group", suggestions...)
    if err != nil {
        util.RofiShowError(err)
    }

    t.Group = val
}

func (t *Task) PromptDueDate() {
    val, err := getValueFromRofi("When is this task due? (Leave blank for no date. Keep in mind that tasks with due dates take priority over those without)", "m/d/yyyy")
    if err != nil {
        util.RofiShowError(err)
    }

    t.UserDueDate, err = time.ParseInLocation(util.DueDateFormat, val, time.Local)
    if err != nil {
        t.UserDueDate = time.Time{}
    }
}

// Set changes the metadata of the task based on the string passed in.
func (t *Task) Set(from string) {
    split := strings.Split(from, " ")

    // get the name by getting the string up to the first metadata tag
    t.Name = getStringToNextTag(split)
    // then, parse the rest of the tags
    parseMetadataTags(split, t)
}

// MarkDoneForToday postpones visibility of the task until tomorrow
func (t *Task) MarkDoneForToday() {
    t.SnoozedUntil = floorDate(time.Now().Add(time.Hour * 24))
}

// MarkComplete marks a task as complete
func (t *Task) MarkComplete() {
    t.Complete = true
}

// IsDoneNow returns true if the task is either complete or done for today
func (t Task) IsDoneNow() bool {
    return t.Complete || time.Now().Before(t.SnoozedUntil)
}

// DisplayString returns the string to be displayed in rofi
func (t Task) DisplayString() string {
    var nextTimeString string
    if !t.Complete && time.Now().Before(t.SnoozedUntil) {
        nextTimeString = "(" + getSnoozeString(t.SnoozedUntil) + ") "
    } else if t.Complete {
        nextTimeString = "(Done!) "
    }
    var dueString string
    if !t.UserDueDate.IsZero() && strings.TrimSpace(util.GetDueDateString(t.UserDueDate)) != "" {
        dueString = ", " + util.GetDueDateString(t.UserDueDate) 
    }
    var groupString string
    if strings.TrimSpace(t.Group) != "" {
        groupString = ", in group " + t.Group
    }
    return nextTimeString + t.Name + dueString + groupString
}

func getSnoozeString(date time.Time) string {
    if date.IsZero() || date.Before(time.Now()) {
        return ""
    }

    daysUntil := util.GetDaysUntil(date)
    switch {
    case daysUntil == 1:
        return "Snoozed until tomorrow"
    case daysUntil < 7:
        dayOfWeek := date.Weekday().String()
        return fmt.Sprint("Snoozed until ", dayOfWeek)
    case daysUntil <= 14:
        return fmt.Sprint("Snoozed for ", daysUntil, " days")
    default:
        return fmt.Sprint("Snoozed until ", date.Format(util.DueDateFormat))
    }
}

// IsLessThan returns true if t should appear above t2 in the list
func (t Task) IsLessThan(t2 Task) bool {
    if t.IsDoneNow() && !t2.IsDoneNow() {
        return false
    } else if !t.IsDoneNow() && t2.IsDoneNow() {
        return true
    }

    switch {
    case t.UserDueDate != t2.UserDueDate:
        return t.UserDueDate.Before(t2.UserDueDate)
    case t.Group != t2.Group:
        return t.Group < t2.Group
    default:
        return t.Name < t2.Name
    }
}

func floorDate(date time.Time) time.Time {
    return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func getValueFromRofi(mesg, prompt string, suggestions ...string) (string, error) {
    in, out, err := util.OpenRofiWithMesg(mesg, prompt)
    if err != nil {
        util.RofiShowError(err)
    }
    defer out.Close()

    for _, s := range suggestions {
        in.Write([]byte(s + "\n"))
    }
    in.Close()

    bufout := bufio.NewReader(out)
    o, _, err := bufout.ReadLine()
    return string(o), err
}
