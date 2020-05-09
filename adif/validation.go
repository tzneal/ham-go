package adif

import "fmt"

// ValidateADIFRecord determines if the record has the minimum required fields to be considered valid.
func ValidateADIFRecord(r Record) error {
	hasFields := map[Identifier]bool{}
	for _, f := range r {
		hasFields[f.Name] = true
	}
	for _, req := range []Identifier{QSODateStart, TimeOn, Call, AMode} {
		if _, ok := hasFields[req]; !ok {
			return fmt.Errorf("missing field: %s", req)
		}
	}
	_, hasBand := hasFields[ABand]
	_, hasFreq := hasFields[Frequency]
	_, hasSatMode := hasFields[SatelliteMode]
	if !hasBand && !hasFreq && !hasSatMode {
		return fmt.Errorf("must have one of %s, %s or %s", ABand, Frequency, SatelliteMode)
	}
	return nil
}

// IsValid returns true if  the record has the minimum required fields to be considered valid.
func IsValid(r Record) bool {
	return ValidateADIFRecord(r) == nil
}
