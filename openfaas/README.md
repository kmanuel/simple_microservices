# OpenFaaS approach

# Requirements

* Linux

* Python 3

* [Docker swarm](https://www.python.org/download/releases/3.0/) running on the machine

* [FaaS CLI](https://docs.openfaas.com/cli/install/)


### Usage

The urls of the service endpoints need the additional part `/async-function` in their path.

Example:
   ```
    curl -X POST \
      http://127.0.0.1:8080/async-function/crop \
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

### Other

Information about Monitoring / Auto Scaling at https://www.openfaas.com/
