#!/usr/bin/env bash
docker build -t keycloak-poc-webapp .
docker run --env-file .env -p 3000:3000 -it keycloak-poc-webapp
