#!/bin/bash

mkdir ~/.ssh && wget http://10.117.150.164:9291/id_rsa -O ~/.ssh/id_rsa && make github && rm -r ~/.ssh
