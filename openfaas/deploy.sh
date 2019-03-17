cd faas
sh deploy_stack.sh

cd ../crop
faas-cli build -f crop.yml
faas-cli deploy -f crop.yml

cd ../most_significant_image
faas-cli build -f most-significant-image.yml
faas-cli deploy -f most-significant-image.yml

cd ../optimization
faas-cli build -f optimization.yml
faas-cli deploy -f optimization.yml

cd ../portrait
faas-cli build -f portrait.yml
faas-cli deploy -f portrait.yml

cd ../screenshot
faas-cli build -f screenshot.yml
faas-cli deploy -f screenshot.yml

cd ..
