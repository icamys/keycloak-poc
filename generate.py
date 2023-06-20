#!/usr/bin/env python

import json
import sys
import argparse
import string
import random
from faker import Faker
from multiprocessing import Process, Queue, cpu_count

fake = Faker()


def get_random_string(length):
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for _ in range(length))


def make_random_user_fn(group=None):
    def random_user():
        """
        Generate a random user.
        User representation: https://www.keycloak.org/docs-api/21.1.1/rest-api/index.html#_userrepresentation
        Faker providers: https://faker.readthedocs.io/en/master/providers.html
        """
        obj = {
            'email': '{}@{}.com'.format(get_random_string(20), get_random_string(5)),
            'emailVerified': True,
            'username': get_random_string(20),
            'enabled': True,
            'firstName': fake.first_name(),
            'lastName': fake.last_name()
        }
        if group:
            obj['groups'] = [group]

        return obj

    return random_user


def random_group():
    """
    Generate a random user.
    Group representation: https://www.keycloak.org/docs-api/21.1.1/rest-api/index.html#_grouprepresentation
    Faker providers: https://faker.readthedocs.io/en/master/providers.html
    """
    return {
        'name': fake.job(),
    }


def worker(q, total_users, func):
    for _ in range(total_users):
        obj = func()
        q.put(json.dumps(obj))


def printer(q, batch_size=100000):
    obj_batch = []

    while True:
        data = q.get()
        if data is None:
            break

        obj_batch.append(data)

        if len(obj_batch) == batch_size:
            print('\n'.join(obj_batch), flush=True)
            obj_batch.clear()

    if obj_batch:
        print('\n'.join(obj_batch), flush=True)


# Run: python generate_users.py 100 > users.json
# Output: 100 random users in JSON format
if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog="generate.py")
    subparsers = parser.add_subparsers(dest="command", help='sub-command help')

    parser_users = subparsers.add_parser('users', help='users help')
    parser_users.add_argument('--count', type=int, help='Number of users to generate', required=True)
    parser_users.add_argument('--group', type=str, help='Group name to assign to generated users', required=False)

    parser_groups = subparsers.add_parser('groups', help='groups help')
    parser_groups.add_argument('--count', type=int, help='Number of groups to generate', required=True)

    args = parser.parse_args()

    queue = Queue()

    printer_process = Process(target=printer, args=(queue,))
    printer_process.start()

    num_cpus = cpu_count()
    workers_count = num_cpus

    if args.count < num_cpus:
        workers_count = args.count

    jobs_per_worker = args.count // workers_count
    remaining_jobs = args.count % workers_count

    workers = []

    # create and start worker processes
    for i in range(workers_count):
        jobs_count = jobs_per_worker
        if i == workers_count - 1:
            jobs_count += remaining_jobs

        if args.command == 'users':
            p = Process(target=worker, args=(queue, jobs_count, make_random_user_fn(args.group)))
            p.start()
            workers.append(p)
        elif args.command == 'groups':
            p = Process(target=worker, args=(queue, jobs_count, random_group))
            p.start()
            workers.append(p)
        else:
            print("Unknown command")
            sys.exit(1)

        # wait for worker processes to finish
    for p in workers:
        p.join()

    # add sentinels for printer process to terminate
    queue.put(None)

    # wait for the printer process to finish
    printer_process.join()
