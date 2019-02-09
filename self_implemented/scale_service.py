import urllib.request
import json
import time
import subprocess
from subprocess import PIPE
import re

def fetch_queues():
    response = urllib.request.urlopen("http://127.0.0.1:8080/info")
    data = json.load(response)
    queues = data['data']['attributes']['queues']
    return queues


def get_current_instances_of_service(service_name):
    out = subprocess.run(["docker", "stack", "services", "--format", "\"{{.Replicas}}\"", "-f", "name=self_impl_swarm_" + service_name, "self_impl_swarm"], stdout=PIPE)
    string_output = out.stdout.decode("utf-8")
    return int(re.search("\/(\d+)", string_output).groups(1)[0])


def scale_service_to(service, instances):
    print("scale service " + service + " to " + str(instances) + " instances")
    subprocess.run(["docker", "service", "scale", "self_impl_swarm_" + service + "=" + str(instances)])


service_names = ["crop", "most_significant_image", "optimization", "portrait", "screenshot"]


while True:
    queues = fetch_queues()
    print(queues)

    for service in service_names:
        instances = get_current_instances_of_service(service)

        open_tasks = queues[service]

        if (instances * 10) < open_tasks or (instances * 10) > open_tasks:
            target_instances = min(5, int(open_tasks / 10))
            if target_instances > 0 and target_instances != instances:
                scale_service_to(service, target_instances)

    time.sleep(1)
