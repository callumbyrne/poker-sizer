.PHONY: dev build clean tailwind templ templ-watch

all: build

tailwind:
	tailwindcss -i web/static/css/input.css -o web/static/css/style.css --minify

tailwind-watch:
	tailwindcss -i web/static/css/input.css -o web/static/css/style.css --watch

templ:
	templ generate

templ-watch:
	templ generate --watch

build: templ tailwind
	go build -o build/poker-sizer cmd/server/main.go

clean:
	rm -rf build
	rm -f web/static/css/output.css
	rm -f web/templates/*.go

dev:
	@echo "Starting development server with air, templ watch, and tailwind watch"
	@$(MAKE) -j3 air templ-watch tailwind-watch

air:
	air

