

# run keycloak with argv[1] as env-file
run:
	docker run -p8080:8080 --env-file .env quay.io/keycloak/keycloak:21.1.1 start-dev
