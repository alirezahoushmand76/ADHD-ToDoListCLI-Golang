package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/todolist/internal/ui"
)

// PomodoroConfig holds configuration for a Pomodoro session
type PomodoroConfig struct {
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	LongBreakInterval  int
	TaskName           string
}

// DefaultPomodoroConfig returns the default Pomodoro configuration
func DefaultPomodoroConfig() PomodoroConfig {
	return PomodoroConfig{
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		LongBreakInterval:  4,
		TaskName:           "Current Task",
	}
}

// PomodoroSession represents a Pomodoro session
type PomodoroSession struct {
	Config       PomodoroConfig
	CurrentCycle int
	IsWorking    bool
	StartTime    time.Time
	EndTime      time.Time
}

// NewPomodoroSession creates a new Pomodoro session with the given configuration
func NewPomodoroSession(config PomodoroConfig) *PomodoroSession {
	return &PomodoroSession{
		Config:       config,
		CurrentCycle: 1,
		IsWorking:    true,
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(config.WorkDuration),
	}
}

// Start starts the Pomodoro session
func (p *PomodoroSession) Start() {
	// Handle Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nPomodoro session interrupted.")
		os.Exit(0)
	}()

	for {
		if p.IsWorking {
			ui.PrintInfo("🍅 WORK CYCLE %d: Focus on '%s' for %s",
				p.CurrentCycle, p.Config.TaskName, formatDuration(p.Config.WorkDuration))
		} else {
			breakType := "SHORT BREAK"
			breakDuration := p.Config.ShortBreakDuration

			if p.CurrentCycle%p.Config.LongBreakInterval == 0 {
				breakType = "LONG BREAK"
				breakDuration = p.Config.LongBreakDuration
			}

			ui.PrintInfo("☕ %s: Take a break for %s", breakType, formatDuration(breakDuration))
		}

		duration := p.Config.WorkDuration
		if !p.IsWorking {
			if p.CurrentCycle%p.Config.LongBreakInterval == 0 {
				duration = p.Config.LongBreakDuration
			} else {
				duration = p.Config.ShortBreakDuration
			}
		}

		// Run the timer
		ui.CountdownTimer(duration, func(remaining time.Duration) {
			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60
			fmt.Printf("\r%02d:%02d remaining", minutes, seconds)
		})

		fmt.Println()

		// Cycle completed
		if p.IsWorking {
			ui.PrintSuccess("✅ Work cycle completed!")
		} else {
			ui.PrintSuccess("✅ Break completed!")
			p.CurrentCycle++
		}

		// Toggle between work and break
		p.IsWorking = !p.IsWorking

		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}
