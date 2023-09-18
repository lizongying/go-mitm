.PHONY: all

all: ui tidy mitm-web

shell:
	@echo 'SHELL='$(SHELL)

ui:
	npm run build --prefix ./web/ui
	cp -R ./web/ui/dist ./static

tidy:
	go mod tidy

mitm-web:
	go vet ./cmd/mitm-web
	go build -ldflags "-s -w" -o ./releases/mitm ./cmd/mitm-web
