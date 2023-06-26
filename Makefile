.PHONY: run
run:
	mkdir -m777 -p .keycloak-data
	docker run -p8080:8080 -v $(PWD)/.keycloak-data:/opt/keycloak/data --env-file .env -e TZ=Europe/London --network=host quay.io/keycloak/keycloak:21.1.1 start-dev
