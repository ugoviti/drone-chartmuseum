docker build -t izdock/$(basename $PWD) . 
docker push izdock/$(basename $PWD)

# test
#docker run -it --rm --name dev izdock/drone-chartmuseum bash
