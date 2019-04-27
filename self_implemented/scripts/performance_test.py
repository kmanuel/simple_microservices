import os
import requests
import time
import bs4
import sys
import time
import subprocess

verbose = False
swarm_name = 'self_impl_swarm'

def create_tasks(request_count):
    data = """{
        "data": {
            "type": "crop_task",
            "attributes": {
                "image_id": "surfcat.jpg",
                "width": 50,
                "height": 10
            }
        }
    }"""

    for i in range(0, request_count):
        response = requests.post('http://127.0.0.1:8080/crop', data=data)
        if verbose and response.status_code != 201:
            print('received unexpected statuscode when creating task=' + str(response.status_code))


def time_till_all_tasks_done(request_count):
    start = time.time()
    while not is_all_tasks_finished():
        time.sleep(0.25)
    return time.time() - start


def is_all_tasks_finished():
    return get_busy_task_count() == 0 and get_queue_task_count() == 0

def get_busy_task_count():
    response = requests.get('http://127.0.0.1:7420/queues')
    soup = bs4.BeautifulSoup(response.content, 'html.parser')
    task_count = soup.select('li.busy.col-sm-1 > a > span.count')
    return int(task_count[0].text.strip())

def get_queue_task_count():
    response = requests.get('http://127.0.0.1:7420/queues')
    soup = bs4.BeautifulSoup(response.content, 'html.parser')
    task_count = soup.select('li.enqueued.col-sm-1 > a > span.count')
    return int(task_count[0].text.strip())


def remove_old_stack():
    os.system('docker stack rm ' + swarm_name + ' > /dev/null')

def deploy_stack(with_resource_constraints):
    if with_resource_constraints:
        compose_file = "scripts/service-swarm-compose.yml "
    else:
        compose_file = "scripts/res-limit-service-swarm-compose.yml"

    res = -1
    while res != 0:
        res = os.system('docker stack deploy --compose-file scripts/infrastructure-swarm-compose.yml --compose-file ' + compose_file + ' ' + swarm_name + ' > /dev/null')
    all_ready = False

    while not all_ready:
        time.sleep(0.25)
        all_ready = True
        replicas = subprocess.check_output("docker service ls | tr -s ' ' | cut -d' ' -f4", shell=True);
        replicas_split = replicas.decode('utf-8').split('\n')
        for replica in replicas_split:
            if replica.startswith('0'):
                all_ready = False


def clear_img_dir():
    os.system('rm -rf ./imgs/crop > /dev/null 2>&1 &')

def setup_img_dir():
    os.system('mkdir ./imgs/crop > /dev/null 2>&1 &')


def perform_run(request_count,
                with_resource_constraints,
                with_auto_scale):
    clear_img_dir()
    setup_img_dir()
    remove_old_stack()
    deploy_stack(with_resource_constraints=with_resource_constraints)
    create_tasks(request_count)
    processing_time = time_till_all_tasks_done(request_count)
    print(str(request_count) + ";" + str(with_resource_constraints) + ";" + str(with_auto_scale) + ';' + str(processing_time))


number_of_runs = 10
request_counts = [1, 10, 100, 250]

for i in range(0, number_of_runs):
    for count in request_counts:
        perform_run(request_count=count, with_resource_constraints=False, with_auto_scale=False)
        perform_run(request_count=count, with_resource_constraints=False, with_auto_scale=True)
        perform_run(request_count=count, with_resource_constraints=True, with_auto_scale=False)
        perform_run(request_count=count, with_resource_constraints=True, with_auto_scale=True)