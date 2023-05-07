build:
	docker buildx build -t sloonz/drone-env-signed .

publish:
	docker push sloonz/drone-env-signed
