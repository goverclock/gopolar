go install ./cmd/...
cd webui
vite build
mv ./dist $GOPATH/bin/gpwebui
