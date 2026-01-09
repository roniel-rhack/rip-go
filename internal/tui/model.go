package tui

import (
	"fmt"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/roniel-rhack/rip-go/internal/process"
	"github.com/roniel-rhack/rip-go/internal/signal"
)

type Model struct {
	processes    []process.Info
	filtered     []int
	selected     map[int32]bool
	cursor       int
	signal       syscall.Signal
	filterText   string
	filtering    bool
	quitting     bool
	results      []string
	width        int
	height       int
	nameWidth    int
	scrollOffset int
	version      string
	sortField    process.SortField
}

func New(processes []process.Info, sig syscall.Signal, version string, sortField process.SortField) Model {
	filtered := make([]int, len(processes))
	for i := range processes {
		filtered[i] = i
	}

	return Model{
		processes: processes,
		filtered:  filtered,
		selected:  make(map[int32]bool),
		signal:    sig,
		width:     80,
		height:    24,
		nameWidth: 35,
		version:   version,
		sortField: sortField,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "q":
			if !m.filtering {
				m.quitting = true
				return m, tea.Quit
			}
			m.filterText += "q"
			m.applyFilter()

		case "esc":
			if m.filtering {
				m.filtering = false
				m.filterText = ""
				m.applyFilter()
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		case "/":
			if !m.filtering {
				m.filtering = true
				m.filterText = ""
			}

		case "backspace":
			if m.filtering && len(m.filterText) > 0 {
				m.filterText = m.filterText[:len(m.filterText)-1]
				m.applyFilter()
			}

		case "enter":
			if m.filtering {
				m.filtering = false
			} else if len(m.selected) > 0 {
				m.killSelected()
				m.quitting = true
				return m, tea.Quit
			}

		case "up", "k":
			if !m.filtering && m.cursor > 0 {
				m.cursor--
				m.ensureVisible()
			}

		case "down", "j":
			if !m.filtering && m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.ensureVisible()
			}

		case " ":
			if !m.filtering && len(m.filtered) > 0 {
				idx := m.filtered[m.cursor]
				pid := m.processes[idx].PID
				if m.selected[pid] {
					delete(m.selected, pid)
				} else {
					m.selected[pid] = true
				}
			}

		case "a":
			if !m.filtering {
				for _, idx := range m.filtered {
					m.selected[m.processes[idx].PID] = true
				}
			} else {
				m.filterText += "a"
				m.applyFilter()
			}

		case "n":
			if !m.filtering {
				m.selected = make(map[int32]bool)
			} else {
				m.filterText += "n"
				m.applyFilter()
			}

		case "1":
			if !m.filtering {
				m.sortField = process.SortByCPU
				m.reSort()
			}

		case "2":
			if !m.filtering {
				m.sortField = process.SortByMem
				m.reSort()
			}

		case "3":
			if !m.filtering {
				m.sortField = process.SortByPID
				m.reSort()
			}

		case "4":
			if !m.filtering {
				m.sortField = process.SortByName
				m.reSort()
			}

		default:
			if m.filtering && len(msg.String()) == 1 {
				m.filterText += msg.String()
				m.applyFilter()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.nameWidth = m.calculateNameWidth()
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting && len(m.results) > 0 {
		return strings.Join(m.results, "\n") + "\n"
	}

	if m.quitting {
		return ""
	}

	var b strings.Builder

	banner := m.renderBanner()
	b.WriteString(banner)
	b.WriteString("\n\n")

	header := m.renderHeader()
	b.WriteString(header)
	b.WriteString("\n")

	visibleRows := m.height - 7
	if visibleRows < 1 {
		visibleRows = 10
	}

	for i := m.scrollOffset; i < len(m.filtered) && i < m.scrollOffset+visibleRows; i++ {
		idx := m.filtered[i]
		proc := m.processes[idx]
		row := m.renderRow(proc, i == m.cursor, m.selected[proc.PID])
		b.WriteString(row)
		b.WriteString("\n")
	}

	rendered := len(m.filtered) - m.scrollOffset
	if rendered < visibleRows {
		for i := rendered; i < visibleRows; i++ {
			b.WriteString("\n")
		}
	}

	if m.filtering {
		filterLine := FilterStyle.Render("/ ") + m.filterText + CursorStyle.Render("█")
		b.WriteString("\n" + filterLine)
	} else if m.filterText != "" {
		filterLine := FilterStyle.Render("Filter: ") + DimStyle.Render(m.filterText)
		b.WriteString("\n" + filterLine)
	} else {
		b.WriteString("\n")
	}

	selectedCount := len(m.selected)
	status := ""
	if selectedCount > 0 {
		status = SelectedCountStyle.Render(fmt.Sprintf(" %d selected", selectedCount))
	}

	help := m.renderHelp()
	b.WriteString("\n" + help + status)

	return b.String()
}

func (m Model) renderBanner() string {
	logo := "█▀█ █ █▀█ ─ █▀▀ █▀█"
	sub := "█▀▄ █ █▀▀   █▄█ █▄█"
	version := VersionStyle.Render(" v" + m.version)
	return BannerStyle.Render(logo) + version + "\n" + BannerStyle.Render(sub)
}

func (m Model) renderHeader() string {
	pidLabel := "PID"
	nameLabel := "NAME"
	cpuLabel := "CPU %"
	memLabel := "MEMORY"

	switch m.sortField {
	case process.SortByCPU:
		cpuLabel = "CPU %▼"
	case process.SortByMem:
		memLabel = "MEMORY▼"
	case process.SortByPID:
		pidLabel = "PID▼"
	case process.SortByName:
		nameLabel = "NAME▼"
	}

	pidCol := HeaderStyle.Render(fmt.Sprintf("%-7s", pidLabel))
	nameCol := HeaderStyle.Render(fmt.Sprintf("%-*s", m.nameWidth, nameLabel))
	cpuCol := HeaderStyle.Render(fmt.Sprintf("%7s", cpuLabel))
	memCol := HeaderStyle.Render(fmt.Sprintf("%9s", memLabel))
	return fmt.Sprintf("      %s %s %s %s", pidCol, nameCol, cpuCol, memCol)
}

func (m Model) renderRow(proc process.Info, isCursor, isSelected bool) string {
	cursor := "  "
	if isCursor {
		cursor = CursorIndicatorStyle.Render("> ")
	}

	checkbox := "[ ] "
	if isSelected {
		checkbox = CheckedStyle.Render("[") + CheckmarkStyle.Render("✓") + CheckedStyle.Render("] ")
	}

	pid := PIDStyle.Render(fmt.Sprintf("%-7d", proc.PID))

	name := proc.Name
	if len(name) > m.nameWidth {
		name = name[:m.nameWidth-3] + "..."
	}
	nameStyled := NameStyle.Render(fmt.Sprintf("%-*s", m.nameWidth, name))

	cpuText := fmt.Sprintf("%5.1f%%", proc.CPU)
	var cpuStyled string
	if proc.CPU > 50.0 {
		cpuStyled = CPUHighStyle.Render(cpuText)
	} else if proc.CPU > 10.0 {
		cpuStyled = CPUMediumStyle.Render(cpuText)
	} else {
		cpuStyled = CPULowStyle.Render(cpuText)
	}

	memText := fmt.Sprintf("%6d MB", proc.Memory)
	memStyled := MemoryStyle.Render(memText)

	row := fmt.Sprintf("%s%s%s %s %s %s", cursor, checkbox, pid, nameStyled, cpuStyled, memStyled)

	if isCursor {
		row = RowHighlightStyle.Render(row)
	}

	return row
}

func (m Model) renderHelp() string {
	if m.filtering {
		return HelpStyle.Render("Type to filter • Enter to confirm • Esc to cancel")
	}
	return HelpStyle.Render("↑↓ navigate • Space select • Enter kill • / filter • 1-4 sort • q quit")
}

func (m *Model) applyFilter() {
	if m.filterText == "" {
		m.filtered = make([]int, len(m.processes))
		for i := range m.processes {
			m.filtered[i] = i
		}
	} else {
		m.filtered = m.filtered[:0]
		filterLower := strings.ToLower(m.filterText)
		for i, proc := range m.processes {
			if strings.Contains(strings.ToLower(proc.Name), filterLower) {
				m.filtered = append(m.filtered, i)
			}
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.scrollOffset = 0
}

func (m *Model) reSort() {
	process.SortProcesses(m.processes, m.sortField)
	m.applyFilter()
	m.cursor = 0
	m.scrollOffset = 0
}

func (m *Model) ensureVisible() {
	visibleRows := m.height - 7
	if visibleRows < 1 {
		visibleRows = 10
	}

	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+visibleRows {
		m.scrollOffset = m.cursor - visibleRows + 1
	}
}

func (m Model) calculateNameWidth() int {
	fixed := 2 + 4 + 8 + 7 + 10 + 4
	available := m.width - fixed
	if available < 20 {
		return 20
	}
	if available > 60 {
		return 60
	}
	return available
}

func (m *Model) killSelected() {
	for _, proc := range m.processes {
		if m.selected[proc.PID] {
			err := signal.Kill(int(proc.PID), m.signal)
			if err != nil {
				m.results = append(m.results,
					ErrorStyle.Render("Failed")+" "+
						lipgloss.NewStyle().Bold(true).Render(proc.Name)+" "+
						DimStyle.Render(fmt.Sprintf("(PID: %d)", proc.PID))+": "+
						err.Error())
			} else {
				m.results = append(m.results,
					SuccessStyle.Render("Killed")+" "+
						lipgloss.NewStyle().Bold(true).Render(proc.Name)+" "+
						DimStyle.Render(fmt.Sprintf("(PID: %d)", proc.PID)))
			}
		}
	}
}
