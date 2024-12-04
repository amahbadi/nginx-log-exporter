# NGINX Log Exporter

A Go-based log exporter that parses NGINX access logs and exposes metrics for Prometheus, including upstream response duration by URI. This exporter watches the log file for new entries, parses them, and updates the metrics accordingly.
Features

- Parses NGINX log entries to extract upstream response time and URI.
- Exposes Prometheus metrics via HTTP at /metrics.
- Configurable log file path, metrics server port, and histogram buckets using environment variables.
- Supports log file watching for real-time metric updates.

## Requirements

- Go 1.18+ (for building and running the exporter)
- Prometheus (for scraping metrics)
- NGINX logs in the specified format

## Installation
Clone the repository

```bash
git clone https://github.com/amahbadi/nginx-log-exporter.git
cd nginx-log-exporter
```
### Build the project

```bash
go build -o nginx-log-exporter cmd/main.go
```

### Configuration
Environment Variables

The exporter is fully configurable via environment variables:
Variable	Description	Default Value

```bash
LOG_FILE_PATH	Path to the NGINX log file to be parsed	(Required)
METRICS_PORT	Port for exposing Prometheus metrics	8080
METRICS_BUCKETS	Comma-separated list of histogram buckets for response times	0.005,0.01,0.025,0.05,0.1,0.25,0.5,1.0
Example
```

### Create a .env file or export the variables directly:

```bash
export LOG_FILE_PATH="/var/log/nginx/access.log"
export METRICS_PORT="8080"
export METRICS_BUCKETS="0.005,0.01,0.025,0.05,0.1,0.25,0.5,1.0"
```
## Usage
Start the exporter

Once you have configured the environment variables, start the exporter:

```bash
go run cmd/main.go
```

The exporter will:

- Watch the NGINX log file for new entries.
- Parse the logs to extract relevant information.
- Update Prometheus metrics for each URI and its corresponding upstream response time.
- Expose metrics at http://localhost:8080/metrics.

## Scraping Metrics with Prometheus

### Add the following scrape job to your prometheus.yml configuration:

```yaml
scrape_configs:
  - job_name: "nginx-log-exporter"
    static_configs:
      - targets: ["localhost:8080"]
```

### Example Prometheus Query

You can query the exported metrics using Prometheus to view upstream response durations for each URI:

```json
upstream_response_duration_seconds_bucket{uri="/example", le="0.005"}
```

- This will show the count of requests to /example with upstream response durations less than or equal to 5ms.
File Parsing and Metrics

- The log entries are parsed from the NGINX access log format:

```nginx
log_format access_log_prometheus $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent
    "$http_referer" "$http_x_forwarded_for" "$http_user_agent"
    $request_time $request_length $upstream_response_time $upstream_addr $upstream_status
    $server_name $server_addr $server_port $uri
```
- The key metric is the upstream response time ($upstream_response_time) and URI ($uri).
