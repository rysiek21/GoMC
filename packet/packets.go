package packet

type State int
type ID int

const (
	StateHandshaking   State = 0
	StateStatus        State = 1
	StateLogin         State = 2
	StateConfiguration State = 3
)

const (
	// State 0
	Handshake ID = 0x00

	// State 1
	StatusRequest ID = 0x00
	StatusPing    ID = 0x01

	// State 2
	LoginHello        ID = 0x00
	LoginAcknowledged ID = 0x03

	// State 3
	ConfigurationClientInfo ID = 0x00
)
