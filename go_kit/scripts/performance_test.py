import deployment
import test_runner

deployment.set_swarm_name('kit_swarm')
deployment.set_compose_file_dir('scripts/compose_files/')

deployment.set_quiet()
deployment.deploy_infrastructure()

#start
deployment.deploy_services(1)
test_runner.do_run(logfile='1_service.csv')
deployment.remove_services()

deployment.deploy_services(2)
test_runner.do_run(logfile='2_service.csv')
# end
deployment.shutdown_stack()
