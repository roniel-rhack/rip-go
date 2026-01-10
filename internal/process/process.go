package process

import (
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type SortField int

const (
	SortByCPU SortField = iota
	SortByMem
	SortByPID
	SortByName
)

func ParseSortField(s string) SortField {
	switch strings.ToLower(s) {
	case "mem", "memory":
		return SortByMem
	case "pid":
		return SortByPID
	case "name":
		return SortByName
	default:
		return SortByCPU
	}
}

type Info struct {
	PID    int32
	Name   string
	CPU    float64
	Memory uint64
}

func GetProcesses(filter string, sortBy SortField) ([]Info, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	for _, p := range procs {
		_, _ = p.CPUPercent()
	}
	time.Sleep(200 * time.Millisecond)

	filterLower := strings.ToLower(filter)
	var result []Info

	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if filter != "" && !strings.Contains(strings.ToLower(name), filterLower) {
			continue
		}

		cpu, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()

		var memMB uint64
		if memInfo != nil {
			memMB = memInfo.RSS / 1024 / 1024
		}

		result = append(result, Info{
			PID:    p.Pid,
			Name:   name,
			CPU:    cpu,
			Memory: memMB,
		})
	}

	SortProcesses(result, sortBy)
	return result, nil
}

func SortProcesses(procs []Info, sortBy SortField) {
	switch sortBy {
	case SortByCPU:
		sort.Slice(procs, func(i, j int) bool { return procs[i].CPU > procs[j].CPU })
	case SortByMem:
		sort.Slice(procs, func(i, j int) bool { return procs[i].Memory > procs[j].Memory })
	case SortByPID:
		sort.Slice(procs, func(i, j int) bool { return procs[i].PID < procs[j].PID })
	case SortByName:
		sort.Slice(procs, func(i, j int) bool {
			return strings.ToLower(procs[i].Name) < strings.ToLower(procs[j].Name)
		})
	}
}
