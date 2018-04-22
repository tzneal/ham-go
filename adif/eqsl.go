package adif

import "fmt"

func ValidateEqsl(l Log) error {
	for _, record := range l.Records {
		if err := ValidateEqslRecord(record); err != nil {
			return err
		}
	}
	return nil
}

func ValidateEqslRecord(r Record) error {
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
