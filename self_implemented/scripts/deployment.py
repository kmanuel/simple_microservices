import os
import subprocess
import time

service_swarm_name = 'self_impl_swarm'

debug = True
cmd_postfix = ''


def set_quiet():
    global cmd_postfix
    global debug
    debug = False
    cmd_postfix = ' > /dev/null'


def deploy_stack(compose_file):
    res = -1
    while res != 0:
        res = os.system('docker stack deploy --compose-file ' + compose_file + ' ' + service_swarm_name + cmd_postfix)
    await_stack_convergence()


def shutdown_stack():
    os.system('docker stack rm ' + service_swarm_name + cmd_postfix)


def deploy_infrastructure():
    deploy_stack('scripts/base_infra.yml')


def deploy_scale():
    deploy_stack('scripts/scale_infra.yml')


def scale_service(service_name, service_instance_count):
    if debug:
        print('scaling service ' + service_name + ' to ' + str(service_instance_count))
    os.system('docker service scale ' + service_swarm_name + '_' + service_name + '=' + str(service_instance_count))
    if debug:
        print('scaled service ' + service_name + ' to ' + str(service_instance_count))


def deploy_services(service_instance_count):
    deploy_stack('scripts/services-res-limit.yml')
    scale_service('crop', service_instance_count)


def remove_service(name):
    os.system('docker service rm ' + service_swarm_name + '_' + name + cmd_postfix)


def remove_services():
    remove_service('alertmanager')
    remove_service('crop')
    remove_service('faktory')
    remove_service('scale-service')


def await_stack_convergence():
    all_ready = False
    while not all_ready:
        time.sleep(0.25)
        all_ready = True
        replicas = subprocess.check_output("docker service ls | tr -s ' ' | cut -d' ' -f4", shell=True);
        replicas_split = replicas.decode('utf-8').split('\n')
        for replica in replicas_split:
            if replica.startswith('0'):
                all_ready = False
