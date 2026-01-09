package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roniel-rhack/rip-go/internal/process"
	"github.com/roniel-rhack/rip-go/internal/signal"
	"github.com/roniel-rhack/rip-go/internal/tui"
	"github.com/spf13/cobra"
)

const Version = "0.1.2"

var (
	filter    string
	signalStr string
	sortBy    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "rip-go",
		Short:   "Fuzzy find and kill processes",
		Long:    "A terminal-based fuzzy process finder and killer",
		Version: Version,
		RunE:    run,
	}

	rootCmd.Flags().StringVarP(&filter, "filter", "f", "", "Pre-filter processes by name")
	rootCmd.Flags().StringVarP(&signalStr, "signal", "s", "KILL", "Signal to send (KILL, TERM, INT, HUP, QUIT)")
	rootCmd.Flags().StringVar(&sortBy, "sort", "cpu", "Sort by: cpu, mem, pid, name")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	sig, err := signal.Parse(signalStr)
	if err != nil {
		return fmt.Errorf("invalid signal: %w", err)
	}

	sortField := process.ParseSortField(sortBy)
	processes, err := process.GetProcesses(filter, sortField)
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	if len(processes) == 0 {
		fmt.Println("No processes found")
		return nil
	}

	model := tui.New(processes, sig, Version, sortField)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}
