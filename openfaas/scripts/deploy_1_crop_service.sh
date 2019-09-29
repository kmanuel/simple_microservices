cd crop
faas-cli deploy --label com.openfaas.scale.min=1 --label com.openfaas.scale.max=1 -f crop.yml
cd ..
