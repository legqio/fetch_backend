# Fetch backend challenge

This repository is a submission of [Fetch's backend receipt processor challenge](https://github.com/fetch-rewards/receipt-processor-challenge).

## Overview
This project implements a webservice following the given [specification](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/api.yml).

## How to run it
This project is configured to run as a docker container. To run it, follow these steps:

1. From the project root folder run `docker build -t legqio/fetch_backend .` to build the docker image

2. Now run `docker run -p 8080:8080 legqio/fetch_backend` to create and run the docker container

3. The webservice is now up and running on your localhost port 8080. (i.e. [http://localhost:8080/receipts/process](http://localhost:8080/receipts/process))

4. You can terminate the webservice by simply stopping or killing the docker container process.

## Credits

This challenge was developed by [Luis Eduardo Gonzalez Quiroz](https://github.com/legqio) using [this](https://github.com/JetBrains/go-code-samples/tree/main/go-rest-demo) repository as a base.
