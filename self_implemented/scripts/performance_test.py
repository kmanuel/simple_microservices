import os
import requests
import time
import bs4
import sys
import time
import subprocess

verbose = False
number_of_runs = 10
request_counts = [1, 10, 50, 100, 150, 250]
with_resource_constraints = False
with_auto_scale = False
service_swarm_name = 'self_impl_swarm'

def create_tasks(request_count):
    data = """{
        "data": {
            "type": "crop_task",
            "attributes": {
                "image_id": "surf_cat.jpg",
                "width": 50,
                "height": 10
            }
        }
    }"""

    for i in range(0, request_count):
        response = requests.post('http://127.0.0.1:8080/crop', data=data)
        if response.status_code != 201:
            print('received unexpected statuscode when creating task=' + str(response.status_code))


def time_till_all_tasks_done(request_count):
    start = time.time()
    while not is_all_tasks_finished():
        if time.time() - start >= 60:
            raise Exception("couldnt finish all task")
        time.sleep(0.25)

    return (time.time() - start)


def is_all_tasks_finished():
    return get_busy_task_count() == 0 and get_queue_task_count() == 0

def get_processed_task_count():
    response = requests.get('http://127.0.0.1:7420/queues')
    soup = bs4.BeautifulSoup(response.content, 'html.parser')
    task_count = soup.select('li.processed.col-sm-1 span.count')
    return int(task_count[0].text.strip().replace(',', ''))


def get_busy_task_count():
    response = requests.get('http://127.0.0.1:7420/queues')
    soup = bs4.BeautifulSoup(response.content, 'html.parser')
    task_count = soup.select('li.busy.col-sm-1 > a > span.count')
    return int(task_count[0].text.strip().replace(',', ''))

def get_queue_task_count():
    response = requests.get('http://127.0.0.1:7420/queues')
    soup = bs4.BeautifulSoup(response.content, 'html.parser')
    task_count = soup.select('li.enqueued.col-sm-1 > a > span.count')
    return int(task_count[0].text.strip().replace(',', ''))



def deploy_stack(with_resource_constraints,
                 with_auto_scale,
                 remove_old_stack=False):
    if remove_old_stack:
        os.system('docker stack rm ' + service_swarm_name + ' > /dev/null')

    if with_resource_constraints:
        compose_file = "scripts/res-limit-service-swarm-compose.yml"
    else:
        compose_file = "scripts/service-swarm-compose.yml "


    res = -1
    while res != 0:
        res = os.system('docker stack deploy --compose-file scripts/infrastructure-swarm-compose.yml --compose-file ' + compose_file + ' ' + service_swarm_name + ' > /dev/null')
    all_ready = False

    while not all_ready:
        time.sleep(0.25)
        all_ready = True
        replicas = subprocess.check_output("docker service ls | tr -s ' ' | cut -d' ' -f4", shell=True);
        replicas_split = replicas.decode('utf-8').split('\n')
        for replica in replicas_split:
            if replica.startswith('0'):
                all_ready = False

    if not with_auto_scale:
        os.system('docker service rm ' + service_swarm_name + '_scale-service')

    print('deployed stack')

def get_scale_to_of_service(service_name):
    replicas = subprocess.check_output("docker service ls | grep \'" + service_name + "' | tr -s ' ' | cut -d' ' -f4", shell=True);
    target_replicas_str = replicas.decode('utf-8').split('/')[1]
    target_replicas_str = target_replicas_str.replace('\n', '')
    return int(target_replicas_str.strip())

def clear_img_dir():
    os.system('rm -rf ./imgs/crop')

def setup_img_dir():
    os.system('mkdir ./imgs/crop')

def scale_back():
    if get_scale_to_of_service('crop') > 1:
        os.system('docker service scale ' + service_swarm_name + '_crop=1')


def remove_services():
    os.system('docker service rm ' + service_swarm_name + '_gateway ' \
              + service_swarm_name + '_crop ' \
              + service_swarm_name + '_most_significant_image ' \
              + service_swarm_name + '_optimization ' \
              + service_swarm_name + '_portrait ' \
              + service_swarm_name + '_scale-service ' \
              + service_swarm_name + '_screenshot')
    print('purging service layer')

def perform_run(request_count, with_resource_constraints, with_auto_scale):
    deploy_stack(with_resource_constraints, with_auto_scale)
    print('running request_count=' + str(request_count))
    clear_img_dir()
    setup_img_dir()
    create_tasks(request_count)
    try:
        processing_time = time_till_all_tasks_done(request_count)
        res = str(request_count) + ';' + str(processing_time)
        print('took: ' + str(processing_time))
    except:
        print('took: timed out')
        res = str(request_count) + ';' + 'timeout'
    write_to_log_file(res)
    remove_services()


def write_to_log_file(text):
    file_name = 'performance_test_resconstraint_' + str(with_resource_constraints) + '_autoscale_' + str(with_auto_scale)
    with open(file_name + '.txt', 'a', buffering=1) as log_file:
        log_file.write(text + '\n')


with_resource_constraints = False
with_auto_scale = False
deploy_stack(with_resource_constraints, with_auto_scale, True)
for i in range(0, number_of_runs):
    print('run nr: ' + str(i))
    for count in request_counts:
        perform_run(request_count=count,
                    with_resource_constraints=with_resource_constraints,
                    with_auto_scale=with_auto_scale)
