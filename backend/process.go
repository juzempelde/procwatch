package procwatch

import (
	"sort"
)

type Processes []*Process

type Process struct {
	Name string
	PID  int
}

type ProcessNameFilter []string

func (filter ProcessNameFilter) Apply(procs Processes) Processes {
	filteredProcs := Processes{}
	for _, proc := range procs {
		for _, allowedName := range filter {
			if proc.Name == allowedName {
				filteredProcs = append(filteredProcs, proc)
				break
			}
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
