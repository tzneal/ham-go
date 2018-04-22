package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// reads cty.csv from http://www.country-files.com/
func main() {
	f, err := os.Open("cty.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cr := csv.NewReader(f)

	op, err := os.Create("../db.go")
	if err != nil {
		panic(err)
	}
	defer op.Close()

	fmt.Fprintf(op, "package dxcc\n")
	fmt.Fprintf(op, "type Entity struct {\n")
	fmt.Fprintf(op, "  Entity string\n")
	fmt.Fprintf(op, "  DXCC int\n")
	fmt.Fprintf(op, "  Continent string\n")
	fmt.Fprintf(op, "  CQZone int\n")
	fmt.Fprintf(op, "  ITUZone int\n")
	fmt.Fprintf(op, "  Latitude float64\n")
	fmt.Fprintf(op, "  Longitude float64\n")
	fmt.Fprintf(op, "  Prefixes []string\n")
	fmt.Fprintf(op, "  PrefixRegexp *regexp.Regexp\n")
	fmt.Fprintf(op, "}\n")
	fmt.Fprintf(op, "var Entities = []Entity{\n")
	for {
		record, err := cr.Read()
		if err != nil {
			break
		}

		entity := record[1]
		prefixes := record[9]
		continent := record[3]
		dxcc := record[2]
		cqZone := record[4]
		ituZone := record[5]
		lat, _ := strconv.ParseFloat(record[6], 64)
		lon, _ := strconv.ParseFloat(record[7], 64)
		// positive is west, so convert to normal form
		lon *= -1

		fmt.Fprintf(op, "{\n")
		fmt.Fprintf(op, `  Entity: "%s",`+"\n", entity)
		fmt.Fprintf(op, `  DXCC: %s,`+"\n", dxcc)
		fmt.Fprintf(op, `  Continent: "%s",`+"\n", continent)
		fmt.Fprintf(op, `  CQZone: %s,`+"\n", cqZone)
		fmt.Fprintf(op, `  ITUZone: %s,`+"\n", ituZone)
		fmt.Fprintf(op, `  Latitude: %f,`+"\n", lat)
		fmt.Fprintf(op, `  Longitude: %f,`+"\n", lon)
		fmt.Fprintf(op, "  Prefixes: []string{%s},\n", splitPrefixes(prefixes))
		fmt.Fprintf(op, `  PrefixRegexp: regexp.MustCompile("%s"),`+"\n", prefixRegexp(prefixes))
		fmt.Fprintf(op, "},")
	}
	fmt.Fprintf(op, "}\n")
}

func splitPrefixes(pfx string) string {
	sb := strings.Builder{}

	pfx = strings.Replace(pfx, ";", "", -1)
	for _, p := range strings.Split(pfx, " ") {
		sb.WriteByte('"')
		sb.WriteString(p)
		sb.WriteByte('"')
		sb.WriteByte(',')
	}
	return sb.String()
}

func prefixRegexp(pfx string) string {

	initialChars := map[byte]struct{}{}
	pfx = strings.Replace(pfx, ";", "", -1)
	for _, p := range strings.Split(pfx, " ") {
		switch p[0] {
		case '=':
			initialChars[p[1]] = struct{}{}
		default:
			initialChars[p[0]] = struct{}{}
		}
	}
	sorted := []byte{}
	for c := range initialChars {
		sorted = append(sorted, c)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	sb := strings.Builder{}
	sb.WriteString("^[")
	for _, v := range sorted {
		sb.WriteByte(v)
	}
	sb.WriteByte(']')
	return sb.String()
}
