package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

func PrettyJson(v interface{}) string {
	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)
	e.SetIndent("", " ")
	err := e.Encode(v)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}

func Json(v interface{}) string {
	return string(JsonBin(v))
}

func JsonBin(v interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)
	err := e.Encode(v)
	if err != nil {
		return nil
	}
	return bytes.TrimSpace(buf.Bytes())
}
