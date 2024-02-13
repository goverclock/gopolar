go install ./cmd/...
cd webui
npm install
vite build
rm -rf $GOPATH/bin/gpwebui
mv ./dist $GOPATH/bin/gpwebui
