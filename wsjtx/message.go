package wsjtx

// Message is a wsjtx decoded message
type Message interface {
	Code() MessageCode
}
