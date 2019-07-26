# SQL-based

For SQL based outputs (e.g., `mysql`, `sqlite`) these are example queries for (the querries assume that the table name with the results is named `iperf3results`):

## `CREATE TABLE`

> **NOTE** This is only necessary if you, e.g., want to import existing CSV results.

```sql
CREATE TABLE `iperf3results` (
  `test_time` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `round` int(11) DEFAULT NULL,
  `tester` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `server_host` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `client_host` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `socket` int(11) DEFAULT NULL,
  `start` float DEFAULT NULL,
  `end` float DEFAULT NULL,
  `seconds` float DEFAULT NULL,
  `bytes` bigint(20) DEFAULT NULL,
  `bits_per_second` float DEFAULT NULL,
  `retransmits` bigint(20) DEFAULT NULL,
  `snd_cwnd` bigint(20) DEFAULT NULL,
  `rtt` bigint(20) DEFAULT NULL,
  `rttvar` bigint(20) DEFAULT NULL,
  `pmtu` int(11) DEFAULT NULL,
  `omitted` tinyint(1) DEFAULT NULL,
  `iperf3_version` varchar(100) COLLATE utf8_bin DEFAULT NULL,
  `system_info` varchar(300) COLLATE utf8_bin DEFAULT NULL,
  `additional_info` varchar(255) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT=' ';
```

## `SELECT` averaged data

```sql
SELECT
	test_time,
	server_host,
	client_host,
	AVG(bits_per_second) AS bits_per_second_avg,
	AVG(bits_per_second / 1000000000) AS gbps_avg,
	AVG(rtt) AS rtt_avg,
	SUM(retransmits) AS total_retransmits
FROM
	`iperf3results`
WHERE
	server_host != client_host
GROUP BY
	round,
	server_host,
	client_host
ORDER BY
	gbps_avg DESC;
```
