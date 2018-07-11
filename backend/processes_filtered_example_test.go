package procwatch_test

import (
	"github.com/juzempelde/procwatch/backend"

	"fmt"
)

func ExampleProcesses_Filtered() {
	byNames := procwatch.ProcessFilterNameList{
		"foo",
		"bar",
	}
	processes := procwatch.Processes{
		&procwatch.Process{
			Name: "foo",
			PID:  444,
		},
		&procwatch.Process{
			Name: "baz",
			PID:  1000,
		},
		&procwatch.Process{
			Name: "bar",
			PID:  1234,
		},
		&procwatch.Process{
			Name: "xyz",
			PID:  666,
		},
	}
	filteredProcesses := processes.Filtered(byNames).SortBy(procwatch.AscendingProcessNames)
	for _, proc := range filteredProcesses {
		fmt.Printf("%05d %s\n", proc.PID, proc.Name)
	}

	// Output:
	// 01234 bar
	// 00444 foo
}
