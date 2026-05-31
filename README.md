# Go Parallel DNS Subdomain Recon
Perform rapid DNS subdomain discoveries concurrently using parallel lookup routines in Go.

## Compile
```bash
go build -o dns-recon main.go
```

## Run
```bash
./dns-recon -d github.com -w subdomains.txt
```