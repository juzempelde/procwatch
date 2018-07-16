package rpc

type ProcessNameFilter struct {
	Names           []string
	RefreshDeadline RefreshDeadlineFunc
}

func (procs *ProcessNameFilter) Expose(request ProcessNameFilterRequest, response *ProcessNameFilterResponse) error {
	refreshDeadline(procs.RefreshDeadline)
	response.Names = procs.Names
	return nil
}

type ProcessNameFilterRequest struct{}

type ProcessNameFilterResponse struct {
	Names []string
}

const processFilter = "processFilter"
