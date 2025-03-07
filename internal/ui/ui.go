package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/user/todolist/internal/models"
)

var (
	// Color functions
	titleColor          = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	descriptionColor    = color.New(color.FgWhite).SprintFunc()
	lowPriorityColor    = color.New(color.FgBlue).SprintFunc()
	mediumPriorityColor = color.New(color.FgYellow).SprintFunc()
	highPriorityColor   = color.New(color.FgRed, color.Bold).SprintFunc()
	completedColor      = color.New(color.FgGreen).SprintFunc()
	categoryColor       = color.New(color.FgCyan).SprintFunc()
	dateColor           = color.New(color.FgMagenta).SprintFunc()
	idColor             = color.New(color.FgHiBlack).SprintFunc()
	errorColor          = color.New(color.FgRed, color.Bold).SprintFunc()
	successColor        = color.New(color.FgGreen, color.Bold).SprintFunc()
	warningColor        = color.New(color.FgYellow, color.Bold).SprintFunc()
	infoColor           = color.New(color.FgCyan, color.Bold).SprintFunc()
)

// PrintTask prints a task with colors
func PrintTask(task *models.Task) {
	status := "[ ]"
	if task.Completed {
		status = completedColor("[✓]")
	} else if task.IsOverdue() {
		status = errorColor("[!]")
	}

	priorityStr := ""
	switch task.Priority {
	case models.PriorityLow:
		priorityStr = lowPriorityColor(string(task.Priority))
	case models.PriorityMedium:
		priorityStr = mediumPriorityColor(string(task.Priority))
	case models.PriorityHigh:
		priorityStr = highPriorityColor(string(task.Priority))
	}

	dueStr := "No due date"
	if !task.DueDate.IsZero() {
		dueStr = dateColor(task.DueDate.Format("2006-01-02 15:04"))
	}

	reminderStr := ""
	if !task.ReminderAt.IsZero() {
		reminderStr = fmt.Sprintf(" (Reminder: %s)", dateColor(task.ReminderAt.Format("2006-01-02 15:04")))
	}

	fmt.Printf("%s %s %s\n", status, titleColor(task.Title), idColor(fmt.Sprintf("(ID: %s)", task.ID)))

	if task.Description != "" {
		fmt.Printf("   %s\n", descriptionColor(task.Description))
	}

	fmt.Printf("   Priority: %s | Category: %s | Due: %s%s\n",
		priorityStr,
		categoryColor(string(task.Category)),
		dueStr,
		reminderStr)

	fmt.Println()
}

// PrintTaskList prints a list of tasks with colors
func PrintTaskList(tasks []*models.Task, title string) {
	if len(tasks) == 0 {
		fmt.Println(infoColor("No tasks found."))
		return
	}

	fmt.Println(titleColor(strings.ToUpper(title)))
	fmt.Println(strings.Repeat("-", len(title)))
	fmt.Println()

	for _, task := range tasks {
		PrintTask(task)
	}

	fmt.Printf(infoColor("Total: %d tasks\n\n"), len(tasks))
}

// PrintError prints an error message
func PrintError(format string, a ...interface{}) {
	fmt.Println(errorColor(fmt.Sprintf(format, a...)))
}

// PrintSuccess prints a success message
func PrintSuccess(format string, a ...interface{}) {
	fmt.Println(successColor(fmt.Sprintf(format, a...)))
}

// PrintWarning prints a warning message
func PrintWarning(format string, a ...interface{}) {
	fmt.Println(warningColor(fmt.Sprintf(format, a...)))
}

// PrintInfo prints an info message
func PrintInfo(format string, a ...interface{}) {
	fmt.Println(infoColor(fmt.Sprintf(format, a...)))
}

// ProgressBar displays a simple progress bar for the Pomodoro timer
func ProgressBar(current, total int, width int) {
	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	fmt.Printf("\r[%s] %d/%d (%d%%)", bar, current, total, int(percent*100))
}

// CountdownTimer displays a countdown timer for the Pomodoro timer
func CountdownTimer(duration time.Duration, onTick func(remaining time.Duration)) {
	start := time.Now()
	end := start.Add(duration)

	for {
		remaining := end.Sub(time.Now())
		if remaining <= 0 {
			break
		}

		if onTick != nil {
			onTick(remaining)
		}

		time.Sleep(time.Second)
	}
}
