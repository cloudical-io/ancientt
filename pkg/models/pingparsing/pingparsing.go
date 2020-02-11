/*
Copyright 2020 Cloudical Deutschland GmbH. All rights reserved.
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

package pingparsing

// ClientResults PingParsing Client Result map
type ClientResults map[string]PingResult

// PingResult Ping target client results structure
type PingResult struct {
	Destination          string      `json:"destination"`
	PacketTransmit       int64       `json:"packet_transmit"`
	PacketReceive        int64       `json:"packet_receive"`
	PacketLossRate       float64     `json:"packet_loss_rate"`
	PacketLossCount      int64       `json:"packet_loss_count"`
	RTTMin               float64     `json:"rtt_min"`
	RTTAvg               float64     `json:"rtt_avg"`
	RTTMax               float64     `json:"rtt_max"`
	RTTMDev              float64     `json:"rtt_mdev"`
	PacketDuplicateRate  float64     `json:"packet_duplicate_rate"`
	PacketDuplicateCount int64       `json:"packet_duplicate_count"`
	ICMPReplies          []ICMPReply `json:"icmp_replies"`
}

// ICMPReply ICMP reply entry
type ICMPReply struct {
	Timestamp string  `json:"timestamp"`
	ICMPSeq   int64   `json:"icmp_seq"`
	TTL       int64   `json:"ttl"`
	Time      float64 `json:"time"`
	Duplicate bool    `json:"duplicate"`
}
