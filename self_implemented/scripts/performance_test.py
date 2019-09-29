import os
import requests
import bs4
import time
import subprocess
import deployment

verbose = False
number_of_runs = 10
request_counts = [10, 100, 250, 500]
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


def time_till_all_tasks_done():
    start = time.time()
    while not is_all_tasks_finished():
        time.sleep(0.25)

    return time.time() - start


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
    print('deploying stack with_resource_constriants=%s and with_auto_scale=%s' % (str(with_resource_constraints), str(with_auto_scale)))

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

def clear_img_dir():
    os.system('rm -rf ./imgs/crop')

def setup_img_dir():
    os.system('mkdir ./imgs/crop')

def remove_services():
    os.system('docker service rm ' + service_swarm_name + '_gateway ' \
              + service_swarm_name + '_crop ' \
              + service_swarm_name + '_most_significant_image ' \
              + service_swarm_name + '_optimization ' \
              + service_swarm_name + '_portrait ' \
              + service_swarm_name + '_scale-service ' \
              + service_swarm_name + '_screenshot')
    print('purging service layer')

def remove_stack():
    os.system('docker stack rm ' + service_swarm_name)

def perform_run(request_count, logfile):
    create_tasks(request_count)
    processing_time = time_till_all_tasks_done()
    res = str(request_count) + ';' + str(processing_time)
    print('took: ' + str(processing_time))
    write_to_log_file(logfile, res)

def write_to_log_file(logfile, text):
    with open(logfile, 'a', buffering=1) as log_file:
        log_file.write(text + '\n')


def do_run(logfile):
    for count in request_counts:
        print('request_count: ' + str(count))
        for i in range(0, number_of_runs):
            print('run nr: ' + str(i))
            perform_run(request_count=count, logfile=logfile)


deployment.deploy_infrastructure()

deployment.deploy_services(1)
do_run('1_services.csv')

deployment.scale_service('crop', 2)
do_run('2_services.csv')

deployment.shutdown_stack()
