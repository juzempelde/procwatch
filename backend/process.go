package procwatch

import (
	"sort"
)

type Processes []*Process

type Process struct {
	Name string
	PID  int
}

type ProcessFilter interface {
	Allows(proc *Process) bool
}

type ProcessFilterNameList []string

func (filter ProcessFilterNameList) Allows(proc *Process) bool {
	for _, allowedName := range filter {
		if proc.Name == allowedName {
			return true
		}
	}
	return false
}

func (procs Processes) Filtered(filter ProcessFilter) Processes {
	filteredProcs := Processes{}
	for _, proc := range procs {
		if filter.Allows(proc) {
			filteredProcs = append(filteredProcs, proc)
		}
	}
	return filteredProcs
}

func (procs Processes) SortBy(less func(left, right *Process) bool) Processes {
	sorter := &sortProcessesBy{
		procs: procs,
		less:  less,
	}
	sort.Sort(sorter)
	return procs
}

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
