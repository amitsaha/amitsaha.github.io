#!/bin/bash

mkdir ~/.ssh && wget http://10.117.150.164:9291/id_rsa -O ~/.ssh/id_rsa && wget -O ~/.ssh/id_rsa.pub http://10.117.150.164:9291/id_rsa.pub && chmod 7000 ~/.ssh && chmod 644 ~/.ssh/id_rsa.pub && chmod 600 ~/.ssh/id_rsa & make github && rm -r ~/.ssh
