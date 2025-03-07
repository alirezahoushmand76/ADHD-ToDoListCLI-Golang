package utils

import (
	"sort"
	"time"

	"github.com/user/todolist/internal/models"
)

// TaskScore represents a task with its calculated score for focus mode
type TaskScore struct {
	Task  *models.Task
	Score float64
}

// FocusMode suggests the next task to focus on based on priority, due date, and other factors
func FocusMode(tasks []*models.Task) []*models.Task {
	// Filter out completed tasks
	var incompleteTasks []*models.Task
	for _, task := range tasks {
		if !task.Completed {
			incompleteTasks = append(incompleteTasks, task)
		}
	}

	if len(incompleteTasks) == 0 {
		return nil
	}

	// Calculate scores for each task
	var scoredTasks []TaskScore
	now := time.Now()

	for _, task := range incompleteTasks {
		score := calculateTaskScore(task, now)
		scoredTasks = append(scoredTasks, TaskScore{
			Task:  task,
			Score: score,
		})
	}

	// Sort tasks by score (descending)
	sort.Slice(scoredTasks, func(i, j int) bool {
		return scoredTasks[i].Score > scoredTasks[j].Score
	})

	// Return tasks in order of priority
	result := make([]*models.Task, len(scoredTasks))
	for i, scoredTask := range scoredTasks {
		result[i] = scoredTask.Task
	}

	return result
}

// calculateTaskScore calculates a score for a task based on various factors
func calculateTaskScore(task *models.Task, now time.Time) float64 {
	var score float64

	// Priority factor (0-100)
	switch task.Priority {
	case models.PriorityHigh:
		score += 100
	case models.PriorityMedium:
		score += 50
	case models.PriorityLow:
		score += 25
	}

	// Due date factor
	if !task.DueDate.IsZero() {
		// Calculate hours until due
		hoursUntilDue := task.DueDate.Sub(now).Hours()

		if hoursUntilDue < 0 {
			// Overdue tasks get highest priority
			score += 200
		} else if hoursUntilDue < 24 {
			// Due within 24 hours
			score += 150
		} else if hoursUntilDue < 48 {
			// Due within 48 hours
			score += 100
		} else if hoursUntilDue < 72 {
			// Due within 72 hours
			score += 75
		} else if hoursUntilDue < 168 {
			// Due within a week
			score += 50
		} else {
			// Due later
			score += 25
		}
	}

	// Reminder factor
	if !task.ReminderAt.IsZero() && now.After(task.ReminderAt) {
		// Reminder is due
		score += 50
	}

	// Age factor - older tasks get a slight boost to prevent them from being forgotten
	ageHours := now.Sub(task.CreatedAt).Hours()
	ageFactor := ageHours / 24 // Days old
	score += ageFactor * 2     // 2 points per day old

	return score
}

// GetNextFocusTask returns the single highest priority task to focus on
func GetNextFocusTask(tasks []*models.Task) *models.Task {
	prioritizedTasks := FocusMode(tasks)
	if len(prioritizedTasks) == 0 {
		return nil
	}
	return prioritizedTasks[0]
}
