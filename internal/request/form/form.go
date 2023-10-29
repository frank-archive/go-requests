// package form implements form encoder
package form

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func toString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}

// encodeUrlValues acts the same as [url.Values.Encode],
// except it's written to [bytes.Buffer]
func encodeUrlValues(v url.Values, w *bytes.Buffer) {
	if v == nil {
		return
	}
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for i, v := range vs {
			if i > 0 {
				w.WriteByte('&')
			}
			w.WriteString(keyEscaped)
			w.WriteByte('=')
			w.WriteString(url.QueryEscape(v))
		}
	}
}

func (d *Encoder) Encode(val map[string]interface{}) error {
	data := url.Values{}
	for name, value := range val {
		switch rv := value.(type) {
		case []interface{}:
			for _, sv := range rv {
				data.Add(name+"[]", fmt.Sprint(sv))
			} // two layers at most, that's where most web frameworks stop
		case map[string]interface{}:
			for mname, mv := range rv {
				kname := fmt.Sprintf("%s[%s]", name, mname)
				data.Add(kname, fmt.Sprint(mv))
			} // two layers at most, that's where most web frameworks stop
		default:
			data.Add(name, fmt.Sprint(rv))
		}
	}
	buf := new(bytes.Buffer)
	encodeUrlValues(data, buf)
	_, err := buf.WriteTo(d.w)
	return err
}
