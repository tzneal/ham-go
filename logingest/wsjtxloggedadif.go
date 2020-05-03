package logingest

// WSJTXLoggedAdif  is sent to  the server(s)  when the
// WSJT-X user accepts the "Log  QSO" dialog by clicking the "OK"
//  button.
type WSJTXLoggedAdif struct {
	ID   string
	ADIF string
}

// Code returns the message op code
func (q WSJTXLoggedAdif) Code() WSJTXMessageCode {
	return MessageLoggedADIF
}
