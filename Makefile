fmt:
	go mod tidy
	go fmt ./...
build: build-ui copy-embed-ui
	go build -o bin/env-hub .
build-ui:
	cd console/environment-hub && npm i && npm run build-only
copy-embed-ui:
	rm -rf cmd/data/assets
	cp -r console/environment-hub/dist/* cmd/data
build-image:
	docker build . -t ghcr.io/linuxsuren/environment-hub:master
test-e2e:
	cd e2e && ./start.sh
init-env:
	curl https://linuxsuren.github.io/tools/install.sh|bash
	hd i cli/cli
	hd i k3d
goreleaser:
	goreleaser build --rm-dist --snapshot
install: build
	sudo install bin/env-hub /usr/local/bin