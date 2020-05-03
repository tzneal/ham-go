package logingest

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"time"
)

type FLDIGIServer struct {
	Messages chan string
	server   *http.Server
	shutdown chan struct{}
}

func NewFLDIGIServer(addr string) (*FLDIGIServer, error) {
	mux := http.NewServeMux()
	h := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s := &FLDIGIServer{
		server:   h,
		shutdown: make(chan struct{}),
		Messages: make(chan string),
	}
	mux.Handle("/RPC2", s)
	return s, nil
}

const (
	listMethods = "system.listMethods"
	addRecord   = "log.add_record"
	methodHelp  = "system.methodHelp"
	checkDup    = "log.check_dup"
	getRecord   = "log.get_record"
)

func (s *FLDIGIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	dec := xml.NewDecoder(r.Body)
	msg := &MethodCall{}
	if err := dec.Decode(msg); err != nil {
		log.Printf("fllog error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	switch msg.Method {
	case listMethods:
		s.listMethods(w, msg)
	case methodHelp:
		s.methodHelp(w, msg)
	case getRecord:
		s.getRecord(w, msg)
	case addRecord:
		s.addRecord(w, msg)
	default:
		log.Println("unhandled fldigi method", msg.Method)
	}
}

func newFldigiRsp() *MethodResponse {
	rsp := &MethodResponse{}
	rsp.Params = &MethodParams{}
	return rsp
}
func newFldigiParam() MethodParam {
	param := MethodParam{}
	param.Value = &ParamValue{}
	param.Value.Array = &ParamValueArray{}
	param.Value.Array.Data = &ParamValueData{}
	return param
}

func (s *FLDIGIServer) addRecord(w http.ResponseWriter, msg *MethodCall) {
	if !fldigiVerifySingleValue(w, msg) {
		log.Printf("FAILED!: %#v", msg.Params.Param[0].Value)
		return
	}

	rsp := newFldigiRsp()
	param := newFldigiParam()
	// fllog returns an empty string
	param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, "")
	rsp.Params.Param = append(rsp.Params.Param, param)
	enc := xml.NewEncoder(w)
	enc.Encode(&rsp)

	s.Messages <- msg.Params.Param[0].Value.Data
}

func (s *FLDIGIServer) getRecord(w http.ResponseWriter, msg *MethodCall) {
	rsp := newFldigiRsp()
	param := newFldigiParam()
	param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, "NO_RECORD")
	rsp.Params.Param = append(rsp.Params.Param, param)
	enc := xml.NewEncoder(w)
	enc.Encode(&rsp)
}

func fldigiVerifySingleValue(w http.ResponseWriter, msg *MethodCall) bool {
	if len(msg.Params.Param) != 1 {
		http.Error(w, "expected a single param", http.StatusBadRequest)
		return false
	}

	if msg.Params.Param[0].Value == nil {
		http.Error(w, "expected a value", http.StatusBadRequest)
		return false
	}
	// single element value
	if len(msg.Params.Param[0].Value.Data) > 0 {
		return true
	}
	// array with a single value
	if msg.Params.Param[0].Value.Array == nil {
		http.Error(w, "expected an array", http.StatusBadRequest)
		return false
	}
	if msg.Params.Param[0].Value.Array.Data == nil {
		http.Error(w, "expected data", http.StatusBadRequest)
		return false
	}
	if len(msg.Params.Param[0].Value.Array.Data.Value) != 1 {
		http.Error(w, "expected a single string", http.StatusBadRequest)
		return false
	}
	return true
}

func (s *FLDIGIServer) methodHelp(w http.ResponseWriter, msg *MethodCall) {
	if msg.Params == nil {
		http.Error(w, "expected params", http.StatusBadRequest)
		return
	}
	if !fldigiVerifySingleValue(w, msg) {
		return
	}

	rsp := newFldigiRsp()
	param := newFldigiParam()
	switch msg.Params.Param[0].Value.Array.Data.Value[0] {
	case addRecord:
		param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, "log.add_record ADIF RECORD")
	case checkDup:
		param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, "log.check_dup CALL, MODE(0), TIME_SPAN(0), FREQ_HZ(0), STATE(0), XCHG_IN(0)")
	case getRecord:
		param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, "log.get_record CALL")
	}

	rsp.Params.Param = append(rsp.Params.Param, param)
	enc := xml.NewEncoder(w)
	enc.Encode(&rsp)
}

func (s *FLDIGIServer) listMethods(w http.ResponseWriter, msg *MethodCall) {
	rsp := newFldigiRsp()
	param := newFldigiParam()
	param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, listMethods)
	param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, addRecord)
	param.Value.Array.Data.Value = append(param.Value.Array.Data.Value, methodHelp)
	rsp.Params.Param = append(rsp.Params.Param, param)
	enc := xml.NewEncoder(w)
	enc.Encode(&rsp)
}

// Run is a non-blocking call that starts the server
func (s *FLDIGIServer) Run() {
	go s.server.ListenAndServe()
}

// Close gracefully shuts down the server
func (s *FLDIGIServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
