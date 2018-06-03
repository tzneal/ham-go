package wsjtx

// Message is a wsjtx decoded message
type Message interface {
	// Code returns the message op code
	Code() MessageCode
}
