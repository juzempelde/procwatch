package rpc

type ProcessNameFilter struct {
	Names []string
}

func (procs *ProcessNameFilter) Expose(request ProcessNameFilterRequest, response *ProcessNameFilterResponse) error {
	response.Names = procs.Names
	return nil
}

type ProcessNameFilterRequest struct{}

type ProcessNameFilterResponse struct {
	Names []string
}

const processFilter = "processFilter"
