package fldigi

import "encoding/xml"

type MethodCall struct {
	XMLName xml.Name      `xml:"methodCall"`
	Method  string        `xml:"methodName"`
	Params  *MethodParams `xml:"params"`
}

type MethodResponse struct {
	XMLName xml.Name      `xml:"methodResponse"`
	Params  *MethodParams `xml:"params"`
}

type MethodParams struct {
	Param []MethodParam `xml:"param"`
}
type MethodParam struct {
	Value *ParamValue `xml:"value"`
}
type ParamValue struct {
	Array *ParamValueArray `xml:"array"`
	Data  string           `xml:",chardata"`
}
type ParamValueArray struct {
	Data *ParamValueData `xml:"data"`
}
type ParamValueData struct {
	Value []string `xml:"value"`
}
