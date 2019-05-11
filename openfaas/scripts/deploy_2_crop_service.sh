cd crop
faas-cli deploy --label com.openfaas.scale.min=2 --label com.openfaas.scale.max=2 -f crop.yml
cd ..
