package procwatch

type Processes []*Process

type Process struct {
	Name string
	PID  int
}
