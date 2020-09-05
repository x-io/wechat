package param

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//Params Params
type Params map[string]interface{}

//Set Set
func (p Params) Set(k string, s interface{}) Params {

	// switch v := s.(type) {
	// case string:
	// 	p[k] = v
	// case int:
	// 	p[k] = strconv.Itoa(v)
	// case int64:
	// 	p[k] = strconv.FormatInt(v, 10)
	// }
	p[k] = s
	return p
}

//Get Get
func (p Params) Get(k string) string {
	s, _ := p[k]

	switch v := s.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	}

	return fmt.Sprintf("%v", s)
}

//GetInt GetInt
func (p Params) GetInt(k string) string {
	s, _ := p[k]

	switch v := s.(type) {
	case string:
		return v
	case int:

		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	}

	return "0"
}

//ContainsKey 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}

//JSON ToJSON
func (p Params) JSON() string {
	str, _ := json.Marshal(p)

	return string(str)
}
