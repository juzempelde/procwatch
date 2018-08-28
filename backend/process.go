package procwatch

import (
	"sort"
)

type Processes []*Process

// Process is a single process running on a machine.
type Process struct {
	Name string
	PID  int
}

// ProcessFilter is used by Processes.Filtered.
type ProcessFilter interface {
	// Allows returns wether the process passes this filter.
	Allows(proc *Process) bool
}

// ProcessFilterNameList is a filter which allows only processes whose name is
// contained in the list of names.
type ProcessFilterNameList []string

// Allows returns true for processes whose name appears in the filter's list.
func (filter ProcessFilterNameList) Allows(proc *Process) bool {
	for _, allowedName := range filter {
		if proc.Name == allowedName {
			return true
		}
	}
	return false
}

// Filtered returns all processes allowed by the filter.
func (procs Processes) Filtered(filter ProcessFilter) Processes {
	filteredProcs := Processes{}
	for _, proc := range procs {
		if filter.Allows(proc) {
			filteredProcs = append(filteredProcs, proc)
		}
	}
	return filteredProcs
}

// SortBy returns the processes in the order determined by the less parameter.
func (procs Processes) SortBy(less func(left, right *Process) bool) Processes {
	sorter := &sortProcessesBy{
		procs: procs,
		less:  less,
	}
	sort.Sort(sorter)
	return procs
}

// AscendingProcessNames is a possible parameter for Processes.SortBy. It sorts
// the processes by name.
func AscendingProcessNames(left, right *Process) bool {
	return left.Name < right.Name
}

type sortProcessesBy struct {
	procs Processes
	less  func(left, right *Process) bool
}

func (procs sortProcessesBy) Len() int {
	return len(procs.procs)
}

func (procs sortProcessesBy) Less(i, j int) bool {
	return procs.less(procs.procs[i], procs.procs[j])
}

func (procs sortProcessesBy) Swap(i, j int) {
	procs.procs[i], procs.procs[j] = procs.procs[j], procs.procs[i]
}
