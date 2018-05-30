package wsjtx

type QSOLogged struct {
	Id uint32 // unique key
}

func (q QSOLogged) Code() MessageCode {
	return MessageQSOLogged
}
