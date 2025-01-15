package client

const (
	OperationTypeGet = "get"
	OperationTypeSet = "set"
	OperationTypeDel = "del"
	OperationTypeMix = "mixed"

	ServerTypeRest = "http"
	ServerTypeTCP  = "tcp"
)

type Client interface {
	Do(c *Operation) error
	PipelinedDo(c []*Operation) error
}

type Operation struct {
	Name  string
	Key   string
	Value string
}

func NewClient(serverType, serverAddr string) Client {
	switch serverType {
	case ServerTypeRest:
		return newHttpClient(serverAddr)
	case ServerTypeTCP:
		return newTCPClient(serverAddr)
	default:
		panic("unknown client type " + serverType)
	}
}
