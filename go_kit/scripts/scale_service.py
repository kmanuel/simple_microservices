import urllib.request
import json
import time
import subprocess
from subprocess import PIPE
import re
import signal
import sys

def signal_handler(sig, frame):
    subprocess.run(["docker", "stack", "rm", "kit_swarm"])
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)

def fetch_queues():
    response = urllib.request.urlopen("http://127.0.0.1:8080/info")
    data = json.load(response)
    queues = data['data']['attributes']['queues']
    return queues


def get_current_instances_of_service(service_name):
    out = subprocess.run(["docker", "stack", "services", "--format", "\"{{.Replicas}}\"", "-f", "name=kit_swarm_" + service_name, "kit_swarm"], stdout=PIPE)
    string_output = out.stdout.decode("utf-8")
    return int(re.search("\/(\d+)", string_output).groups(1)[0])


def scale_service_to(service, instances):
    print("scale service " + service + " to " + str(instances) + " instances")
    subprocess.run(["docker", "service", "scale", "-d", "kit_swarm_" + service + "=" + str(instances)])


service_names = ["crop", "most_significant_image", "optimization", "portrait", "screenshot"]


while True:
    try:
        queues = fetch_queues()
    except:
        print("scale_service.py: could not fetch task queue data")
        time.sleep(1)
        continue

    print(queues)

    for service in service_names:
        if queues.get(service) == None:
            continue

        instances = get_current_instances_of_service(service)

        open_tasks = queues[service]

        if (instances * 10) < open_tasks or (instances * 10) > open_tasks:
            target_instances = min(5, int(open_tasks / 10))
            if target_instances > 0 and target_instances != instances:
                scale_service_to(service, target_instances)

    time.sleep(1)
