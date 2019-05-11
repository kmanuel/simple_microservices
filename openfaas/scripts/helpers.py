import subprocess
import time
import requests


def get_image_count():
    return int(subprocess.check_output("ls imgs/crop | wc -l", shell=True))


def wait_until_imgs_increased_by(start_imgs_count, task_count):
    run_nr = 0
    finished = False
    while not finished:
        run_nr += 1
        finished = True
        curr_img_count = get_image_count()
        if not curr_img_count - start_imgs_count >= task_count:
            if run_nr % 20 == 0:
                print('tasks done: ' + str(curr_img_count - start_imgs_count))
            time.sleep(0.1)
            finished = False
    print('done')


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
        response = requests.post('http://127.0.0.1:8080/async-function/crop', data=data)
        if (response.status_code - 200) > 99:
            print('received unexpected statuscode when creating task=' + str(response.status_code))
