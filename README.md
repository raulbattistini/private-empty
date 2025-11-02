## System Metrics Collector

> ####  Implemented in Go, it stores into a lightweight `.sqlite` and displays alerts with sound (WIP) too. Tracks (currently): CPU, Memory and Net Usage.
Also features: a light API in Gin also displays the data as well (WIP). Supports a (almost) graceful shutdown, communication with channels handling mutexes.

#### How to run the app:
1. To run the project
```bash
  go version # important to check the installed version
```

2. Install deps
```bash
  go mod tidy
```

3. Run the project (wait a few seconds till it logs the CLI)
```bash
  go run `cmd/agent/main.go`
```

#### Disclaimers:
1. Values for threshold as low as 0.5% are used to mock as in
```go
	alerter := alert.NewAlerter()
	alerter.AddRule(alert.Rule{
		Source:    "cpu",
		Condition: enum.Above,
		Threshold: 0.05, // here
		Duration:  1 * time.Second, // or time lasted
		Action: func(m metrics.Metric) {
			log.Warn("🚨 ALERT: CPU above 5%% for 6s: %.2f%%", m.Value)
			alert.DesktopNotifier(m)
		},
	})

	alerter.AddRule(alert.Rule{
		Source:    "memory",
		Condition: enum.Above,
		Threshold: 256.0, // 4GB
		Duration:  30 * time.Second,
		Action: func(m metrics.Metric) {
			log.Warn("🚨 ALERT: Memory above 250GB for 30s: %.2f MB", m.Value)
			if cfg.Notification.Annoying {
				alert.AnnoyingDesktopNotifier(m) // AnnoyingDesktopNotifier was made to resemble Windows behavior
			} else {
				alert.DesktopNotifier(m)  
			}
		},
	})
// rest of the code
```
2. and so is the timeout to collect new hardware information -- which should be more frequent than the current arbitrary values on `config.yaml` (adjust as per your needs):
```yaml
collectors:
  cpu:
    enabled: true
    interval: 7s
  memory:
    enabled: true
    interval: 7s
  net:
    enabled: true 
    interval: 7s
# other configs
```
3. And also values used for most metrics are simulated trying to keep as real as possible, OS-specific code would require some platform-bound code. (it might be added later on still)

#### TODO: 
1. Notify based on percentages instead of absolute values
2. Folder detailed description & docs 
3. Enhanced webpage
4. `Makefile`, Docker support
5. Clean up folder architecture


#### App's general flow is `Collectors → eventChannel → Agent (holding dispatcher logic - could be decoupled into `events` too) → Exporters`