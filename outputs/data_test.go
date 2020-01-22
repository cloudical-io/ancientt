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

package outputs

import (
	"fmt"
	"testing"

	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestDataTableTransform(t *testing.T) {
	dataTable := Table{
		Headers: []*Row{
			&Row{Value: "bits_per_second"},
			&Row{Value: "willremain"},
			&Row{Value: "replacedwithwillremain"},
		},
		Rows: [][]*Row{
			[]*Row{
				&Row{Value: float64(123.0)},
				&Row{Value: "nope"},
				&Row{Value: int64(50)},
			},
			[]*Row{
				&Row{Value: int64(15)},
				&Row{Value: "nope"},
				&Row{Value: int64(30)},
			},
			[]*Row{
				&Row{Value: int64(15)},
				&Row{Value: "nope"},
				&Row{Value: int64(75)},
			},
		},
	}

	transformations := []*config.Transformation{
		{
			Action:         config.TransformationActionAdd,
			Source:         "bits_per_second",
			Destination:    "gigabits_per_second",
			Modifier:       util.FloatPointer(float64(100)),
			ModifierAction: config.ModifierActionDivison,
		},
		{
			Source: "bits_per_second",
			Action: config.TransformationActionDelete,
		},
		{
			Action:         config.TransformationActionReplace,
			Source:         "replacedwithwillremain",
			Destination:    "tb_per_second",
			Modifier:       util.FloatPointer(float64(1000)),
			ModifierAction: config.ModifierActionMultiply,
		},
	}

	fmt.Println("BEFORE TRANSFORMATION:")
	pp.Println(dataTable)

	err := dataTable.Transform(transformations)
	assert.Nil(t, err)

	fmt.Println("===\nAFTER TRANSFORMATION:")
	pp.Println(dataTable)
}
