go install ./cmd/...
cd webui
npm install
vite build
sudo rm -rf $GOPATH/bin/gpwebui
sudo mv ./dist $GOPATH/bin/gpwebui
