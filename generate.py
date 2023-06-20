#!/usr/bin/env python

import json
import sys
import argparse
from faker import Faker

fake = Faker()


def random_user():
    """
    Generate a random user.
    User representation: https://www.keycloak.org/docs-api/21.1.1/rest-api/index.html#_userrepresentation
    Faker providers: https://faker.readthedocs.io/en/master/providers.html
    """
    return {
        'email': fake.email(),
        'emailVerified': True,
        'username': fake.user_name(),
        'enabled': True,
        'firstName': fake.first_name(),
        'lastName': fake.last_name()
    }


def random_group():
    """
    Generate a random user.
    Group representation: https://www.keycloak.org/docs-api/21.1.1/rest-api/index.html#_grouprepresentation
    Faker providers: https://faker.readthedocs.io/en/master/providers.html
    """
    return {
        'name': fake.job(),
    }


# Run: python generate_users.py 100 > users.json
# Output: 100 random users in JSON format
if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog="generate.py")
    subparsers = parser.add_subparsers(dest="command", help='sub-command help')

    parser_users = subparsers.add_parser('users', help='users help')
    parser_users.add_argument('--count', type=int, help='Number of users to generate', required=True)
    parser_users.add_argument('--group', type=str, help='Group name to assign to generated users', required=True)

    parser_groups = subparsers.add_parser('groups', help='groups help')
    parser_groups.add_argument('--count', type=int, help='Number of groups to generate', required=True)

    args = parser.parse_args()

    if args.command == 'users':
        for i in range(args.count):
            user = random_user()
            if args.group:
                user['groups'] = [args.group]
            print(json.dumps(user))
    elif args.command == 'groups':
        for i in range(args.count):
            print(json.dumps(random_group()))
    else:
        print("Unknown command")
        sys.exit(1)
