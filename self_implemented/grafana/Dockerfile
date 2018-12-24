FROM grafana/grafana

COPY datasource.yaml $GF_PATHS_PROVISIONING/datasources/datasource.yaml
COPY dashboards.yaml $GF_PATHS_PROVISIONING/dashboards/dashboards.yaml
COPY dashboards /var/lib/grafana/dashboards/
