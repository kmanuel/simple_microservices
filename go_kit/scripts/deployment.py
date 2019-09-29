import os
import subprocess
import time

debug = True
cmd_postfix = ''


def set_swarm_name(swarm_name):
    global service_swarm_name
    service_swarm_name = swarm_name

def set_compose_file_dir(dir):
    global compose_file_dir
    compose_file_dir = dir

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
    deploy_stack(compose_file_dir + 'base_infra.yml')


def deploy_scale():
    deploy_stack(compose_file_dir + 'scale_infra.yml')


def deploy_services(service_count):
    print('deploy services count=' + str(service_count))
    deploy_stack(compose_file_dir + 'services.yml')
    os.system("docker service scale " + service_swarm_name + "_crop=" + str(service_count))
    print('wait for stack convergence')
    await_stack_convergence()
    print('stack converged')



def remove_service(name):
    os.system('docker service rm ' + service_swarm_name + '_' + name + cmd_postfix)


def remove_services():
    remove_service('crop')


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
