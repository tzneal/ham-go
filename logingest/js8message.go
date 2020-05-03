package logingest

import (
	"bytes"
	"encoding/json"
)

type JS8Message struct {
	Type   string                 `json:"type"`
	Value  string                 `json:"value"`
	Params map[string]interface{} `json:"params"`
}

func JS8Decode(msg []byte) (JS8Message, error) {
	r := bytes.NewReader(msg)
	dec := json.NewDecoder(r)

	var jmsg JS8Message
	if err := dec.Decode(&jmsg); err != nil {
		return JS8Message{}, err
	}
	return jmsg, nil
}
