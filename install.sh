go install ./cmd/...
cd webui
vite build
rm -rf $GOPATH/bin/gpwebui
mv ./dist $GOPATH/bin/gpwebui
