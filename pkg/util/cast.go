/*
Copyright 2019 Cloudical Deutschland GmbH. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
)

// CastToString cast interface{} to string
func CastToString(val interface{}) string {
	if cVal, ok := val.(float32); ok {
		return fmt.Sprintf("%f", cVal)
	} else if cVal, ok := val.(float64); ok {
		return fmt.Sprintf("%f", cVal)
	} else if cVal, ok := val.(int); ok {
		return fmt.Sprintf("%d", cVal)
	} else if cVal, ok := val.(int8); ok {
		return fmt.Sprintf("%d", cVal)
	} else if cVal, ok := val.(int16); ok {
		return fmt.Sprintf("%d", cVal)
	} else if cVal, ok := val.(int32); ok {
		return fmt.Sprintf("%d", cVal)
	} else if cVal, ok := val.(int64); ok {
		return fmt.Sprintf("%d", cVal)
	}
	return fmt.Sprintf("%v", val)
}

// CastNumberToFloat64 cast a number from interface{} to float64
func CastNumberToFloat64(in interface{}) (float64, error) {
	switch in.(type) {
	case float64:
		return in.(float64), nil
	case float32:
		return in.(float64), nil
	case int:
		return float64(in.(int)), nil
	case int8:
		return float64(in.(int8)), nil
	case int16:
		return float64(in.(int16)), nil
	case int32:
		return float64(in.(int32)), nil
	case int64:
		return float64(in.(int64)), nil
	default:
		return 0.0, fmt.Errorf("non number character can't be casted to float64")
	}
}
