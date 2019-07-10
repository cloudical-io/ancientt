/*
Copyright 2019 Cloudical Deutschland GmbH
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

package parsers

import (
	"bytes"

	"github.com/cloudical-io/acntt/pkg/config"
)

// NameIPerf3 IPerf3 tester name
const NameIPerf3 = "iperf3"

func init() {
	Factories[NameIPerf3] = NewIPerf3Tester
}

// IPerf3 IPerf3 tester structure
type IPerf3 struct {
	Parser
	config *config.IPerf3
}

// NewIPerf3Tester return a new IPerf3 tester instance
func NewIPerf3Tester(cfg *config.Config, test *config.Test) (Parser, error) {
	return IPerf3{
		config: test.IPerf3,
	}, nil
}

// Parse parse IPerf3 JSON responses
func (ip IPerf3) Parse(in *bytes.Buffer) ([]byte, error) {
	// TODO parse input

	/*
	   {
	   	"start":	{
	   		"connected":	[{
	   				"socket":	5,
	   				"local_host":	"100.67.207.107",
	   				"local_port":	49472,
	   				"remote_host":	"100.69.96.141",
	   				"remote_port":	5601
	   			}],
	   		"version":	"iperf 3.6",
	   		"system_info":	"Linux acntt-iperf3-8da0bd6d47135c1c127754b32a87b5c7902716e5 5.0.5-200.fc29.x86_64 #1 SMP Wed Mar 27 20:58:04 UTC 2019 x86_64",
	   		"timestamp":	{
	   			"time":	"Wed, 10 Jul 2019 22:07:41 GMT",
	   			"timesecs":	1562796461
	   		},
	   		"connecting_to":	{
	   			"host":	"100.69.96.141",
	   			"port":	5601
	   		},
	   		"cookie":	"jwmntlwujrbdk5fhhwfrvdgihdy5akurejjn",
	   		"tcp_mss_default":	1388,
	   		"sock_bufsize":	0,
	   		"sndbuf_actual":	87380,
	   		"rcvbuf_actual":	87380,
	   		"test_start":	{
	   			"protocol":	"TCP",
	   			"num_streams":	1,
	   			"blksize":	131072,
	   			"omit":	0,
	   			"duration":	10,
	   			"bytes":	0,
	   			"blocks":	0,
	   			"reverse":	0,
	   			"tos":	0
	   		}
	   	},
	   	"intervals":	[{
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	0,
	   					"end":	1.0001039505004883,
	   					"seconds":	1.0001039505004883,
	   					"bytes":	106014052,
	   					"bits_per_second":	848024263.45338786,
	   					"retransmits":	0,
	   					"snd_cwnd":	285928,
	   					"rtt":	2570,
	   					"rttvar":	147,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	0,
	   				"end":	1.0001039505004883,
	   				"seconds":	1.0001039505004883,
	   				"bytes":	106014052,
	   				"bits_per_second":	848024263.45338786,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	1.0001039505004883,
	   					"end":	2.0000720024108887,
	   					"seconds":	0.99996805191040039,
	   					"bytes":	104327632,
	   					"bits_per_second":	834647721.40018737,
	   					"retransmits":	0,
	   					"snd_cwnd":	283152,
	   					"rtt":	2419,
	   					"rttvar":	152,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	1.0001039505004883,
	   				"end":	2.0000720024108887,
	   				"seconds":	0.99996805191040039,
	   				"bytes":	104327632,
	   				"bits_per_second":	834647721.40018737,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	2.0000720024108887,
	   					"end":	3.0001049041748047,
	   					"seconds":	1.000032901763916,
	   					"bytes":	104519176,
	   					"bits_per_second":	836125897.98310053,
	   					"retransmits":	0,
	   					"snd_cwnd":	285928,
	   					"rtt":	2386,
	   					"rttvar":	154,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	2.0000720024108887,
	   				"end":	3.0001049041748047,
	   				"seconds":	1.000032901763916,
	   				"bytes":	104519176,
	   				"bits_per_second":	836125897.98310053,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	3.0001049041748047,
	   					"end":	4.0001139640808105,
	   					"seconds":	1.0000090599060059,
	   					"bytes":	105029960,
	   					"bits_per_second":	840232067.57644463,
	   					"retransmits":	0,
	   					"snd_cwnd":	280376,
	   					"rtt":	2629,
	   					"rttvar":	262,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	3.0001049041748047,
	   				"end":	4.0001139640808105,
	   				"seconds":	1.0000090599060059,
	   				"bytes":	105029960,
	   				"bits_per_second":	840232067.57644463,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	4.0001139640808105,
	   					"end":	5.0000739097595215,
	   					"seconds":	0.99995994567871094,
	   					"bytes":	104263784,
	   					"bits_per_second":	834143683.05908251,
	   					"retransmits":	0,
	   					"snd_cwnd":	285928,
	   					"rtt":	2512,
	   					"rttvar":	116,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	4.0001139640808105,
	   				"end":	5.0000739097595215,
	   				"seconds":	0.99995994567871094,
	   				"bytes":	104263784,
	   				"bits_per_second":	834143683.05908251,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	5.0000739097595215,
	   					"end":	6.0001788139343262,
	   					"seconds":	1.0001049041748047,
	   					"bytes":	105413048,
	   					"bits_per_second":	843215927.12898231,
	   					"retransmits":	0,
	   					"snd_cwnd":	280376,
	   					"rtt":	2450,
	   					"rttvar":	117,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	5.0000739097595215,
	   				"end":	6.0001788139343262,
	   				"seconds":	1.0001049041748047,
	   				"bytes":	105413048,
	   				"bits_per_second":	843215927.12898231,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	6.0001788139343262,
	   					"end":	7.0000710487365723,
	   					"seconds":	0.99989223480224609,
	   					"bytes":	103561456,
	   					"bits_per_second":	828580940.1888746,
	   					"retransmits":	0,
	   					"snd_cwnd":	283152,
	   					"rtt":	3055,
	   					"rttvar":	358,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	6.0001788139343262,
	   				"end":	7.0000710487365723,
	   				"seconds":	0.99989223480224609,
	   				"bytes":	103561456,
	   				"bits_per_second":	828580940.1888746,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	7.0000710487365723,
	   					"end":	8.0001029968261719,
	   					"seconds":	1.0000319480895996,
	   					"bytes":	103114520,
	   					"bits_per_second":	824889806.346557,
	   					"retransmits":	0,
	   					"snd_cwnd":	283152,
	   					"rtt":	2538,
	   					"rttvar":	220,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	7.0000710487365723,
	   				"end":	8.0001029968261719,
	   				"seconds":	1.0000319480895996,
	   				"bytes":	103114520,
	   				"bits_per_second":	824889806.346557,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	8.0001029968261719,
	   					"end":	9.0001189708709717,
	   					"seconds":	1.0000159740447998,
	   					"bytes":	103880696,
	   					"bits_per_second":	831032293.0529207,
	   					"retransmits":	0,
	   					"snd_cwnd":	283152,
	   					"rtt":	2449,
	   					"rttvar":	154,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	8.0001029968261719,
	   				"end":	9.0001189708709717,
	   				"seconds":	1.0000159740447998,
	   				"bytes":	103880696,
	   				"bits_per_second":	831032293.0529207,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}, {
	   			"streams":	[{
	   					"socket":	5,
	   					"start":	9.0001189708709717,
	   					"end":	10.000102043151855,
	   					"seconds":	0.99998307228088379,
	   					"bytes":	100752144,
	   					"bits_per_second":	806030796.26291811,
	   					"retransmits":	0,
	   					"snd_cwnd":	5552,
	   					"rtt":	356,
	   					"rttvar":	29,
	   					"pmtu":	1440,
	   					"omitted":	false
	   				}],
	   			"sum":	{
	   				"start":	9.0001189708709717,
	   				"end":	10.000102043151855,
	   				"seconds":	0.99998307228088379,
	   				"bytes":	100752144,
	   				"bits_per_second":	806030796.26291811,
	   				"retransmits":	0,
	   				"omitted":	false
	   			}
	   		}],
	   	"end":	{
	   		"streams":	[{
	   				"sender":	{
	   					"socket":	5,
	   					"start":	0,
	   					"end":	10.000102043151855,
	   					"seconds":	10.000102043151855,
	   					"bytes":	1040876468,
	   					"bits_per_second":	832692677.34146774,
	   					"retransmits":	0,
	   					"max_snd_cwnd":	285928,
	   					"max_rtt":	3055,
	   					"min_rtt":	356,
	   					"mean_rtt":	2336
	   				},
	   				"receiver":	{
	   					"socket":	5,
	   					"start":	0,
	   					"end":	10.037973880767822,
	   					"seconds":	10.000102043151855,
	   					"bytes":	1039596732,
	   					"bits_per_second":	828531131.36052859
	   				}
	   			}],
	   		"sum_sent":	{
	   			"start":	0,
	   			"end":	10.000102043151855,
	   			"seconds":	10.000102043151855,
	   			"bytes":	1040876468,
	   			"bits_per_second":	832692677.34146774,
	   			"retransmits":	0
	   		},
	   		"sum_received":	{
	   			"start":	0,
	   			"end":	10.037973880767822,
	   			"seconds":	10.037973880767822,
	   			"bytes":	1039596732,
	   			"bits_per_second":	828531131.36052859
	   		},
	   		"cpu_utilization_percent":	{
	   			"host_total":	2.5373511447898793,
	   			"host_user":	0.13308812038044104,
	   			"host_system":	2.4042630244094383,
	   			"remote_total":	14.269335116370662,
	   			"remote_user":	0.62156681352299481,
	   			"remote_system":	13.647790285499692
	   		},
	   		"sender_tcp_congestion":	"bbr",
	   		"receiver_tcp_congestion":	"bbr"
	   	}
	   }
	*/

	return in.Bytes(), nil
}
