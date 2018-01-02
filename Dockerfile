FROM ubuntu:latest

RUN apt-get update && apt-get -y install python3-pip make bash
RUN pip3 install pelican pelican-youtube markdown pelican-gist
USER $user
RUN git clone https://github.com/gfidente/pelican-svbhack /tmp/pelican-svbhack
RUN git clone --recursive https://github.com/getpelican/pelican-plugins /tmp/pelican-plugins
WORKDIR /site
ENTRYPOINT ["make", "build"]
