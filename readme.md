# Helm Live

Helm live is a project to render helm charts in the browser. 
It is a tool to help helm developers to visualize their charts and values in a more interactive way.

## Usage


```bash
helm-live \
  --static-path=./ui/dist/helm-live/browser \
  --chart-path=./helm-demo \
  --value-file=./helm-demo/values-example.yaml 
```

The open your browser and go to `http://localhost:8085/ui/`


## Development

```bash
go mod tidy                
go mod vendor
GIT_TERMINAL_PROMPT=1 GO111MODULE=on go build -o helm-live   ./cmd/live/live.go 
```