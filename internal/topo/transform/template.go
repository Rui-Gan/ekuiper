// Copyright 2022-2023 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lf-edge/ekuiper/internal/conf"
	"github.com/lf-edge/ekuiper/internal/converter"
	"github.com/lf-edge/ekuiper/pkg/ast"
	"github.com/lf-edge/ekuiper/pkg/message"
	"text/template"
)

type TransFunc func(interface{}) ([]byte, bool, error)
type SelectFunc func([]byte, []string) ([]byte, bool, error)

func GenTransform(dt string, format string, schemaId string, delimiter string) (TransFunc, SelectFunc, error) {
	var (
		tp  *template.Template = nil
		c   message.Converter
		out message.Converter
		err error
	)
	switch format {
	case message.FormatProtobuf, message.FormatCustom:
		c, err = converter.GetOrCreateConverter(&ast.Options{FORMAT: format, SCHEMAID: schemaId})
		if err != nil {
			return nil, nil, err
		}
		out, _ = converter.GetOrCreateConverter(&ast.Options{FORMAT: format, SCHEMAID: schemaId})
	case message.FormatDelimited:
		c, err = converter.GetOrCreateConverter(&ast.Options{FORMAT: format, DELIMITER: delimiter})
		if err != nil {
			return nil, nil, err
		}
		out, _ = converter.GetOrCreateConverter(&ast.Options{FORMAT: format, DELIMITER: delimiter})
	}

	if dt != "" {
		temp, err := template.New("sink").Funcs(conf.FuncMap).Parse(dt)
		if err != nil {
			return nil, nil, err
		}
		tp = temp
	}
	return func(d interface{}) ([]byte, bool, error) {
			var (
				bs          []byte
				transformed bool
			)
			if tp != nil {
				var output bytes.Buffer
				err := tp.Execute(&output, d)
				if err != nil {
					return nil, false, fmt.Errorf("fail to encode data %v with dataTemplate for error %v", d, err)
				}
				bs = output.Bytes()
				transformed = true
			}
			switch format {
			case message.FormatJson:
				if transformed {
					return bs, transformed, nil
				}
				j, err := json.Marshal(d)
				return j, false, err
			case message.FormatProtobuf, message.FormatCustom, message.FormatDelimited:
				if transformed {
					m := make(map[string]interface{})
					err := json.Unmarshal(bs, &m)
					if err != nil {
						return nil, false, fmt.Errorf("fail to decode data %s after applying dataTemplate for error %v", string(bs), err)
					}
					d = m
				}
				b, err := c.Encode(d)
				return b, transformed, err
			default: // should not happen
				return nil, false, fmt.Errorf("unsupported format %v", format)
			}
		}, func(bytes []byte, fields []string) ([]byte, bool, error) {
			if fields == nil {
				return bytes, false, fmt.Errorf("unsupported fields %v", fields)
			}
			var m interface{}
			switch format {
			case message.FormatJson:
				err = json.Unmarshal(bytes, &m)
				if err != nil {
					return bytes, false, err
				}
				switch m.(type) {
				case []interface{}:
					mm := m.([]interface{})
					outputs := make([]map[string]interface{}, len(mm))
					for i, v := range mm {
						if out, ok := v.(map[string]interface{}); !ok {
							return bytes, false, fmt.Errorf("fail to decode json, unsupported type %v", mm)
						} else {
							outputs[i] = selectMap(out, fields)
						}
					}
					jsonBytes, err := json.Marshal(outputs)
					return jsonBytes, true, err
				case map[string]interface{}:
					mm := m.(map[string]interface{})
					jsonBytes, err := json.Marshal(selectMap(mm, fields))
					return jsonBytes, true, err
				default:
					return bytes, false, fmt.Errorf("fail to decode json, unsupported type %v", m)
				}

			case message.FormatProtobuf, message.FormatCustom, message.FormatDelimited:
				m, err = c.Decode(bytes)
				if err != nil {
					return bytes, false, err
				}
				mm, ok := m.(map[string]interface{})
				if !ok {
					return bytes, false, fmt.Errorf("expect map[string]interface{} but got %T", m)
				}
				outBytes, err := out.Encode(selectMap(mm, fields))
				if err != nil {
					return bytes, false, err
				}
				return outBytes, true, nil

			default:
				return bytes, false, fmt.Errorf("unsupported format %v", format)
			}

		}, nil
}

func GenTp(dt string) (*template.Template, error) {
	return template.New("sink").Funcs(conf.FuncMap).Parse(dt)
}

func selectMap(input map[string]interface{}, fields []string) map[string]interface{} {
	output := make(map[string]interface{})
	for _, field := range fields {
		if v, ok := input[field]; ok {
			output[field] = v
		}
	}
	return output

}
