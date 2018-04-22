package rigcontrol

// Mode is the mode of operation
//go:generate stringer -type=Mode
type Mode byte

// Mode constants
const (
	ModeLSB Mode = 0x00
	ModeUSB Mode = 0x01
	ModeCW  Mode = 0x02
	ModeCWR Mode = 0x03
	ModeAM  Mode = 0x04
	ModeWFM Mode = 0x06
	ModeFM  Mode = 0x08
	ModeDIG Mode = 0x0A
	ModePKT Mode = 0x0C
	ModeCWN Mode = 0x82
	ModeNFM Mode = 0x88
)

// SquelchStatus is the status of the squelch (off = signal present, on = no signal)
//go:generate stringer -type=SquelchStatus
type SquelchStatus byte

// SquelchStatus constants
const (
	SquelchUnknown SquelchStatus = iota
	SquelchOff
	SquelchOn
)

type Status struct {
	Frequency float64 // frequency in MHz
	Mode      Mode
	Squelch   SquelchStatus
	SMeter    byte // SMeter (i.e. 9 = S9)
}
type Rig interface {
	Close() error
	ReadStatus() (*Status, error)
	Tune(freq float64) error
	SetMode(mode Mode) error
}
