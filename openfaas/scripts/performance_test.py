import time

import deployment
import helpers


runs_count = 10
task_counts = [10, 100, 250, 500]


def perform_run(count):
    start_imgs_count = helpers.get_image_count()
    print(f'start_imgs_count={start_imgs_count}')
    helpers.create_tasks(count)
    start = time.time()
    helpers.wait_until_imgs_increased_by(start_imgs_count, count)
    end = time.time()
    return end - start


def write_to_log_file(logfile, text):
    with open(logfile, 'a', buffering=1) as log_file:
        log_file.write(text + '\n')




def run(log_file_name):
    for count in task_counts:
        for run_nr in range(0, runs_count):
            print('run nr=%s count=%s' % (str(run_nr), str(count)))
            execution_time = perform_run(count)
            write_to_log_file(log_file_name, (str(count) + ';' + str(execution_time)))
            print('took: ' + str(execution_time))



deployment.deploy_base_stack()
# deployment.deploy_crop_service(True)
# run('1_services.csv')
# deployment.remove_services()

deployment.deploy_crop_service(False)
run('2_services.csv')

deployment.shut_down_stack()
