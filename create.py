#!/usr/bin/env python

import json
import argparse
import sys
from keycloak import KeycloakAdmin
from keycloak import KeycloakOpenIDConnection

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--username", help="username of the admin user", required=True)
    parser.add_argument("--password", help="password of the admin user", required=True)
    parser.add_argument("--realm", help="Realm name. Example: master", required=True)
    parser.add_argument("--server_url", help="Server url. Example: http://localhost:8080/", required=True)

    subparsers = parser.add_subparsers(dest="command", help='sub-command help')
    parser_users = subparsers.add_parser('users', help='Users help')
    parser_groups = subparsers.add_parser('groups', help='Groups help')
    args = parser.parse_args()

    keycloak_connection = KeycloakOpenIDConnection(
        server_url=args.server_url,
        username=args.username,
        password=args.password,
        realm_name=args.realm,
        verify=True,
    )

    keycloak_admin = KeycloakAdmin(connection=keycloak_connection)

    for line in sys.stdin:
        if args.command == 'users':
            new_user = keycloak_admin.create_user(json.loads(line))
        elif args.command == 'groups':
            new_group = keycloak_admin.create_group(json.loads(line))
        else:
            print("Unknown command")
            sys.exit(1)
