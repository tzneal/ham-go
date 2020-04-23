package wsjtx

// LoggedADIF  is sent to  the server(s)  when the
// WSJT-X user accepts the "Log  QSO" dialog by clicking the "OK"
//  button.
type LoggedADIF struct {
	ID   string
	ADIF string
}

// Code returns the message op code
func (q LoggedADIF) Code() MessageCode {
	return MessageLoggedADIF
}
