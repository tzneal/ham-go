package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
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

	src := &bytes.Buffer{}

	fmt.Fprint(src, "package dxcc\n")
	fmt.Fprint(src, "import \"regexp\"\n")
	fmt.Fprint(src, "type Entity struct {\n")
	fmt.Fprint(src, "  Entity string\n")
	fmt.Fprint(src, "  DXCC int\n")
	fmt.Fprint(src, "  Continent string\n")
	fmt.Fprint(src, "  CQZone int\n")
	fmt.Fprint(src, "  ITUZone int\n")
	fmt.Fprint(src, "  Latitude float64\n")
	fmt.Fprint(src, "  Longitude float64\n")
	fmt.Fprint(src, "  Prefixes []string\n")
	fmt.Fprint(src, "  Score int\n")
	fmt.Fprint(src, "  PrefixRegexp *regexp.Regexp\n")
	fmt.Fprint(src, "}\n")
	fmt.Fprint(src, "var Entities = []Entity{\n")
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

		fmt.Fprintf(src, "{\n")
		fmt.Fprintf(src, `  Entity: "%s",`+"\n", entity)
		fmt.Fprintf(src, `  DXCC: %s,`+"\n", dxcc)
		fmt.Fprintf(src, `  Continent: "%s",`+"\n", continent)
		fmt.Fprintf(src, `  CQZone: %s,`+"\n", cqZone)
		fmt.Fprintf(src, `  ITUZone: %s,`+"\n", ituZone)
		fmt.Fprintf(src, `  Latitude: %f,`+"\n", lat)
		fmt.Fprintf(src, `  Longitude: %f,`+"\n", lon)
		fmt.Fprintf(src, "  Prefixes: []string{%s},\n", splitPrefixes(prefixes))
		fmt.Fprintf(src, `  PrefixRegexp: regexp.MustCompile("%s"),`+"\n", prefixRegexp(prefixes))
		fmt.Fprintf(src, "},")
	}
	fmt.Fprintf(src, "}\n")

	b := src.Bytes()
	b, err = format.Source(b)
	if err != nil {
		fmt.Printf("%s\n", string(src.Bytes()))
		panic(err)
	}
	if err := ioutil.WriteFile("../db.go", b, 0666); err != nil {
		log.Fatalf("can't write output: %v\n", err)
	}

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
	var sorted []byte
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
