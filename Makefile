.PHONY: all

all: ui tidy mitm-web

shell:
	@echo 'SHELL='$(SHELL)

ui:
	npm run build --prefix ./web/ui
	rm -rf ./static/dist
	cp -R ./web/ui/dist ./static

tidy:
	go mod tidy

mitm-web:
	go vet ./cmd/mitm-web

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mitm_linux_amd64 ./cmd/mitm-web

	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mitm_linux_arm64 ./cmd/mitm-web

	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mitm_darwin_amd64 ./cmd/mitm-web

	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mitm_darwin_arm64 ./cmd/mitm-web

	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mitm_windows_amd64.exe ./cmd/mitm-web