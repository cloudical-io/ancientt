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

package iperf3

// ClientResult IPerf3 client result output
type ClientResult struct {
	Start     Start      `json:"start"`
	Intervals []Interval `json:"intervals"`
	End       End        `json:"end"`
}

// Start
type Start struct {
	Connected     []ConnectedEntry `json:"connected"`
	Version       string           `json:"version"`
	SystemInfo    string           `json:"system_info"`
	Timestamp     Timestamp        `json:"timestamp"`
	ConnectingTo  ConnectingTo     `json:"connecting_to"`
	Cookie        string           `json:"cookie"`
	TCPMSSDefault int64            `json:"tcp_mss_default"`
	SockBufsize   int64            `json:"sock_bufsize"`
	SndbufActual  int64            `json:"sndbuf_actual"`
	RcvbufActual  int64            `json:"rcvbuf_actual"`
	PlannedTime   PlannedTime      `json:"test_start"`
}

// ConnectedEntry
type ConnectedEntry struct {
	Socket     int    `json:"socket"`
	LocalHost  string `json:"local_host"`
	LocalPort  int    `json:"local_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
}

// Timestamp
type Timestamp struct {
	Time     string `json:"time"`
	Timesecs int64  `json:"timesecs"`
}

// ConnectingTo
type ConnectingTo struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

// PlannedTime
type PlannedTime struct {
	Protocol   string `json:"protocol"`
	NumStreams int64  `json:"num_streams"`
	BlkSize    int64  `json:"blksize"`
	Omit       int64  `json:"omit"`
	Duration   int64  `json:"duration"`
	Bytes      int64  `json:"bytes"`
	Blocks     int64  `json:"blocks"`
	Reverse    int64  `json:"reverse"`
	Tos        int64  `json:"tos"`
}

// Interval
type Interval struct {
	Streams []Stream `json:"streams"`
	Sum     Sum      `json:"sum"`
}

// Stream
type Stream struct {
	Socket        int64   `json:"socket"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int64   `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int64   `json:"retransmits"`
	SndCwnd       int64   `json:"snd_cwnd"`
	RTT           int64   `json:"rtt"`
	RTTVar        int64   `json:"rttvar"`
	PMTU          int64   `json:"pmtu"`
	Omitted       bool    `json:"omitted"`
}

// Sum
type Sum struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int64   `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int64   `json:"retransmits"`
	Omitted       bool    `json:"omitted"`
}

// End
type End struct {
	Streams               []EndStream           `json:"streams"`
	SumSent               SumSent               `json:"sum_sent"`
	SumReceived           SumReceived           `json:"sum_received"`
	CPUUtilizationPercent CPUUtilizationPercent `json:"cpu_utilization_percent"`
	SenderTCPCongestion   string                `json:"sender_tcp_congestion"`
	ReceiverTCPCongestion string                `json:"receiver_tcp_congestion"`
}

// EndStream
type EndStream struct {
	Sender   Sender   `json:"sender"`
	Receiver Receiver `json:"receiver"`
}

// Sender
type Sender struct {
	Socket        int64   `json:"socket"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int64   `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int64   `json:"retransmits"`
	MaxSndCwnd    int64   `json:"max_snd_cwnd"`
	MaxRTT        int64   `json:"max_rtt"`
	MinRTT        int64   `json:"min_rtt"`
	MeanRTT       int64   `json:"mean_rtt"`
}

// Receiver
type Receiver struct {
	Socket        int64   `json:"socket"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int64   `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

// SumSent
type SumSent struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int64   `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int64   `json:"retransmits"`
}

// SumReceived
type SumReceived struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int64   `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

// CPUUtilizationPercent
type CPUUtilizationPercent struct {
	HostTotal    float64 `json:"host_total"`
	HostUser     float64 `json:"host_user"`
	HostSystem   float64 `json:"host_system"`
	RemoteTotal  float64 `json:"remote_total"`
	RemoteUser   float64 `json:"remote_user"`
	RemoteSystem float64 `json:"remote_system"`
}
