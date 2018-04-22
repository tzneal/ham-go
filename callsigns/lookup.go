package callsigns

type Lookup interface {
	Lookup(call string) (*Response, error)
}

type Response struct {
	Name      *string
	Grid      *string
	Latitude  *float64
	Longitude *float64
	Country   *string
	DXCC      *int
	CQZone    *int
	ITUZone   *int
}
