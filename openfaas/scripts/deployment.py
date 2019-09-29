import os
import subprocess
import time


def set_resource_limit(cpu, memory):
    os.system("docker service update --limit-cpu=%s --limit-memory=%sM crop" % (str(cpu), str(memory)))


def deploy_base_stack():
    os.system('sh deploy_faas_and_minio.sh')
    await_stack_convergence()


def deploy_crop_service(single_instance):
    if single_instance:
        os.system('sh scripts/deploy_1_crop_service.sh')
        os.system("docker service scale func_queue-worker=1")
    else:
        os.system('sh scripts/deploy_2_crop_service.sh')
        os.system("docker service scale func_queue-worker=2")
    set_resource_limit(0.5, 350)
    await_stack_convergence()


def await_stack_convergence():
    all_ready = False
    while not all_ready:
        time.sleep(0.25)
        all_ready = True
        replicas = subprocess.check_output("docker service ls | tr -s ' ' | cut -d' ' -f4", shell=True)
        replicas_split = replicas.decode('utf-8').split('\n')
        for replica in replicas_split:
            if replica.startswith('0'):
                all_ready = False


def remove_services():
    os.system('faas-cli remove crop')

def shut_down_stack():
    os.system('docker stack rm func')
    os.system('docker service rm `docker service ls -q`')

    all_down = False
    while not all_down:
        time.sleep(0.25)
        all_down = True
        docker_processes = int(subprocess.check_output("docker ps | wc -l", shell=True)) - 1
        if docker_processes > 0:
            all_down = False
