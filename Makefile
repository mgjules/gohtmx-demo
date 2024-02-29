dev:
	@air

templ:
	@templ generate --watch --proxy="http://localhost:8080"

tailwind:
	@npx tailwindcss -i ./assets/src/tailwind.css -o ./assets/dist/app.css --watch