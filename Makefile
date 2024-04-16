docker.build:
	npx tailwindcss -i ./static/input.css -o ./static/output.css --minify
	docker build -t gofiber-library:latest .
docker.push:
	docker tag gofiber-library:latest $(IMAGE_URL)
	docker push $(IMAGE_URL)

air:
	air
db:
	docker-compose up
tailwind:
	npx tailwindcss -i ./static/input.css -o ./static/output.css --watch
