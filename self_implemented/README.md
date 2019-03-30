# Self Implemented Approach


## Requirements

* Linux

* [Docker swarm](https://www.python.org/download/releases/3.0/) running on the machine

## How to run

Run the script `run_swarm.sh` in the `self_implemented` directory. 
This will build all required docker containers for the application, deploy the stack to the docker swarm 
and start a script to scale the system's services according to their current load. 

## How to use

### Input/Output Files

Image files to be processed must be put into the directory `./imgs/images`, 
where `./imgs` is the default mounted directory by the minio container.

Each processing service has it's own output directory where it will put the resulting file of the image transformation:

Crop -> `./imgs/crop/`  
Most Significant Image -> `./imgs/mostsignificantimage/`  
Optimization -> `./imgs/optimization/`  
Portrait -> `./imgs/portrait/`  
Screenshot -> `./imgs/screenshot`  


### Service Endpoints

Each service has it's own endpoint. A call to the endpoint will trigger an image transformation with the given parameters:

1. Cropping an image
    ```
    curl -X POST \
      http://127.0.0.1:8080/crop \
      -H 'Content-Type: application/json' \
      -H 'cache-control: no-cache' \
      -d '{
        "data": {
            "type": "crop_task",
            "attributes": {
                "image_id": "example.jpg",
                "width": 50,
                "height": 10
            }
        }
    }'
    ```
    
    This will create a crop_task, that   
    takes the image *example.jpg* (which must be present in the `./imgs/images` folder)  
    and crops the image to a new image with the given *width* and *height* in the attributes.  
    Once the operation is complete the corresponding output folder `./images/crop` contains the resulting image 
    (with a naming pattern of `<original_file_name>_<parameters>_<timestamp>.jpg`).

2. Getting the most significant image of a web page
    ```
    curl -X POST \
      http://127.0.0.1:8080/most_significant_image \
      -H 'Content-Type: application/json' \
      -H 'cache-control: no-cache' \
      -d '{
        "data": {
            "attributes": {
                "url": "https://www.xkcd.com"
            }
        }
    }'
    ```

3. Performing several optimization tasks on an image
    ```
    curl -X POST \
      http://127.0.0.1:8080/optimization \
      -H 'Content-Type: application/json' \
      -H 'Postman-Token: d07af64a-d6eb-4006-a4de-d55726d3d1b5' \
      -H 'cache-control: no-cache' \
      -d '{
        "data": {
            "attributes": {
                "image_id": "example.jpg"
            }
        }
    }'
    ```

4. Crop a image with face detection
    ```
    curl -X POST \
      http://127.0.0.1:8080/portrait \
      -H 'Content-Type: application/json' \
      -H 'cache-control: no-cache' \
      -d '{
        "data": {
            "attributes": {
                "image_id": "surf_cat.jpg",
                "width": 50,
                "height": 50
            }
        }
    }'
    ```

5. Take a screenshot of a webpage
    ```
    curl -X POST \
      http://127.0.0.1:8080/screenshot \
      -H 'Content-Type: application/json' \
      -H 'cache-control: no-cache' \
      -d '{
        "data": {
            "attributes": {
                "url": "https://www.example.com"
            }
        }
    }'
    ```

### Monitoring

A [Grafana dashboard](`http://127.0.0.1:3000`) displays simple metrics provided by the swarms [Prometheus](http://127.0.0.1:9000) server.

### Scaling

Is achieved by a simple python script (`./scripts/scale_service.py`) running on the host machine.  
This script periodically polls the stack for the number of pending tasks for each service, prints the result to the console and scales a service up if the number exceeds a certain threshold (10 * #services_for__task). 
