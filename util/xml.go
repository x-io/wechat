package util

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

//XMLToMap XMLToMap
func XMLToMap(xmlStr string) map[string]interface{} {

	params := make(map[string]interface{})
	// the current value stack
	values := make([]string, 0)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.CharData: // 标签内容
			values = append(values, string([]byte(token)))
		case xml.EndElement:
			if token.Name.Local == "xml" || token.Name.Local == "langs" {
				continue
			}
			params[token.Name.Local] = values[len(values)-1]
			// pop
			values = values[:len(values)]
		}
	}

	return params
}

//MapToXML MapToXML
func MapToXML(params map[string]interface{}) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`>`)

		switch vv := v.(type) {
		case string:
			buf.WriteString(fmt.Sprintf(`<![CDATA[%s]]>`, vv))
			break
		case int:
			buf.WriteString(fmt.Sprintf(`%d`, vv))
			break
		case int64:
			buf.WriteString(fmt.Sprintf(`%d`, vv))
			break
		default:
			buf.WriteString(fmt.Sprintf(`%s`, vv))
			break
		}

		buf.WriteString(`</`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)
	return buf.String()
}
