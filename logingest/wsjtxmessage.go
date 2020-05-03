package logingest

// WSJTXMessage is a wsjtx decoded message
type WSJTXMessage interface {
	// Code returns the message op code
	Code() WSJTXMessageCode
}
