### A simple project to track CPU, Memory and Network usage (sampling some values)
> Uses some DI, channels in a simple way

#### App's general flow is `Collectors → eventChannel → Agent (holding dispatcher logic - could be decoupled into `events` too) → Exporters`


##### Disclaimer: values used for most metrics are simulated trying to keep as real as possible, OS-specific code would require some platform-bound code. (it might be added later on still)


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
