import requests
import time
import bs4
import sys


verbose = False
request_count = int(sys.argv[1])

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

if verbose:
    print('creating ' + str(request_count) + ' tasks')
for i in range(0, request_count):
    response = requests.post('http://127.0.0.1:8080/crop', data=data)
    if response.status_code != 201:
        print('received unexpected statuscode when creating task=' + str(response.status_code))

if verbose:
    print('created tasks')

all_tasks_finished = False

start = time.time()

while not all_tasks_finished:
        all_tasks_finished = True
        response = requests.get('http://127.0.0.1:7420/queues')
        soup = bs4.BeautifulSoup(response.content, 'html.parser')
        task_rows = soup.select('table.queues tr') #> td:nth-child(2)
        for row in task_rows:
            tds = row.select('td')
            task_type = tds[0].text.strip()
            number_text = tds[1].text.replace(',', '')
            number_of_tasks_pending = int(number_text.strip())
            if number_of_tasks_pending > 0:
                all_tasks_finished = False
                time.sleep(0.25)

end = time.time()

if verbose:
    print('duration: ' + str(end - start))
else:
    print(str(end- start))

