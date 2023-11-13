package Format

type (
	Kommunikation struct {
		KommType
		Data
	}

	KommType int

	Data map[string]string
)

const (
	SETUP KommType = iota
	GET
	GET_RESPONSE
	PING
	PING_RESPONSE
)
