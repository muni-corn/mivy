package util

import (
    "fmt"
    "log"
    "time"
    "math"
    "io"
    "os/exec"
)

// DueDateFormat is the standard date format, all-American
const DueDateFormat = "1/2/2006"

// GetDaysUntil returns the amount of days until
// the given date
func GetDaysUntil(date time.Time) int {
	untilDueDate := math.Ceil((float64(time.Until(date)) / float64(time.Hour*24)))

	return int(untilDueDate)
}

func GetDueDateString(dueDate time.Time) string {
    if dueDate.IsZero() {
        return ""
    }

    daysUntil := GetDaysUntil(dueDate)
    switch {
    case daysUntil == 1:
        return "due tomorrow"
    case daysUntil == 0:
        return "due today"
    case daysUntil == -1:
        return "due yesterday"
    case daysUntil < 0:
        return fmt.Sprint("overdue by ", -daysUntil, " days")
    case daysUntil < 7:
        dayOfWeek := dueDate.Weekday().String()
        return fmt.Sprint("due on ", dayOfWeek)
    case daysUntil <= 14:
        return fmt.Sprint("due in ", daysUntil, " days")
    default:
        return fmt.Sprint("due on ", dueDate.Format(DueDateFormat))
    }
}

func OpenRofiWithMesg(mesg, prompt string) (io.WriteCloser, io.ReadCloser, error) {
    log.Printf("open rofi with mesg: %s", mesg)
    cmd := exec.Command("rofi", "-dmenu", "-p", prompt, "-mesg", mesg)

    // get pipes
    in, err := cmd.StdinPipe()
    if err != nil {
        log.Printf("err getting stdin pipe for rofi: %s", err.Error())
        exec.Command("rofi", "-e", err.Error())
        return in, nil, err 
    }
    out, err := cmd.StdoutPipe()
    if err != nil {
        log.Printf("err getting stdout pipe for rofi: %s", err.Error())
        exec.Command("rofi", "-e", err.Error())
        return in, out, err
    }

    err = cmd.Start()
    if err != nil {
        log.Println(err.Error())
        RofiShowError(err)
    }

    return in, out, err
}

func RofiShowError(err error) {
    log.Printf("show error with rofi: %s", err)
    exec.Command("rofi", "-e", err.Error()).Start()
}
