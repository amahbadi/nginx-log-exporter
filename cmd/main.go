package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"nginx-log-exporter/internal/config"
	"nginx-log-exporter/internal/metrics"

	"github.com/fsnotify/fsnotify"
)

type LogEntry struct {
	URI                  string
	UpstreamResponseTime float64
}

func parseLog(logLine string) (*LogEntry, error) {
	pattern := `^(?P<remote_addr>\S+) - (?P<remote_user>\S+) \[(?P<time_local>.+?)\] "(?P<request>.+?)" (?P<status>\d+) (?P<body_bytes_sent>\d+) "(?P<http_referer>.*?)" "(?P<http_x_forwarded_for>\S+)" "(?P<http_user_agent>.*?)" (?P<request_time>\S+) (?P<request_length>\d+) (?P<upstream_response_time>\S+) (?P<upstream_addr>\S+) (?P<upstream_status>\d+) (?P<server_name>\S+) (?P<server_addr>\S+) (?P<server_port>\d+) (?P<uri>\S+)$`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(logLine)

	if match == nil {
		return nil, nil // Skip lines that don't match
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	upstreamResponseTime, _ := strconv.ParseFloat(result["upstream_response_time"], 64)

	return &LogEntry{
		URI:                  result["uri"],
		UpstreamResponseTime: upstreamResponseTime,
	}, nil
}

func watchLogFile(filePath string, m *metrics.Metrics) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					for {
						line, err := reader.ReadString('\n')
						if err != nil {
							break
						}
						line = strings.TrimSpace(line)
						entry, err := parseLog(line)
						if err != nil {
							log.Printf("Error parsing log: %v", err)
							continue
						}
						if entry != nil {
							m.ObserveUpstreamDuration(entry.URI, entry.UpstreamResponseTime)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatalf("Error watching file: %v", err)
	}
	log.Printf("Watching file: %s", filePath)

	<-make(chan struct{})
}

func main() {

	cfg := config.LoadConfig()

	m := metrics.NewMetrics(cfg.DefaultBuckets)

	go metrics.ServeMetrics(cfg.MetricsPort)

	watchLogFile(cfg.LogFilePath, m)
}
