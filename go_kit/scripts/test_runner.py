import bs4
import requests
import time

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
        if response.status_code != 201 and response.status_code != 200:
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


def perform_run(request_count, logfile):
    create_tasks(request_count)
    processing_time = time_till_all_tasks_done()
    run_log_string = str(request_count) + ';' + str(processing_time)
    write_to_log_file(logfile, run_log_string)


def write_to_log_file(logfile, text):
    print(text)
    with open(logfile, 'a', buffering=1) as log_file:
        log_file.write(text + '\n')


def do_run(logfile):
    for count in request_counts:
        print('request_count: ' + str(count))
        for i in range(0, number_of_runs):
            print('run nr: ' + str(i))
            perform_run(request_count=count, logfile=logfile)
