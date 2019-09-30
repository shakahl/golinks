package main

import (
	"bytes"
	"encoding/xml"
	"strconv"
)

// SafeParseInt ...
func SafeParseInt(s string, d int) int {
	n, e := strconv.Atoi(s)
	if e != nil {
		return d
	}
	return n
}

// EscapeXML makes string XML-safe.
func EscapeXML(s string) (string, error) {
	b := &bytes.Buffer{}
	if err := xml.EscapeText(b, []byte(s)); err != nil {
		return "", err
	}
	return b.String(), nil
}
