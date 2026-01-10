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

const Version = "0.2.2"

var (
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
	model := tui.New(sig, Version, sortField)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}
