Generate groups:

```bash
./generate.py groups --count 1 > groups.ndjson
```

Create groups in Keycloak from generated file:

```bash
cat groups.ndjson | ./create.py \
  --username=$KEYCLOAK_ADMIN \
  --password=$KEYCLOAK_ADMIN_PASSWORD \
  --realm='master' \
  --server_url='http://localhost:8080/' groups
```


Generate users:
```bash
./generate.py users --count 1 > users.ndjson
```

Generate users with assigned group:
```bash
./generate.py users --count 1 --group 'Personal assistant' > users_g.ndjson
```

Create users in Keycloak from generated file:

```bash
cat users.ndjson | ./create.py \
    --username=$KEYCLOAK_ADMIN \
    --password=$KEYCLOAK_ADMIN_PASSWORD \
    --realm='master' \
    --server_url='http://localhost:8080/' users
```
