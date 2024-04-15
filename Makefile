air:
	air
db:
	docker-compose up
tailwind:
	npx tailwindcss -c ./tailwind.config.js -i ./static/input.css -o ./static/output.css --watch
