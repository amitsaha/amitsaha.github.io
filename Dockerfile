FROM ubuntu:latest

ARG git_username
ARG git_useremail
ARG gid
ARG group
ARG uid
ARG user
RUN groupadd -g $gid  $group || true && useradd -u $uid -g $gid -d /home/$user $user

RUN apt-get update && apt-get -y install python3-pip wget git make bash
RUN git config --global user.name $git_username && git config --global user.email $git_useremail
RUN pip3 install pelican pelican-youtube 
USER $user
RUN git clone https://github.com/gfidente/pelican-svbhack /tmp/pelican-svbhack
RUN ls -lrt /tmp/pelican-svbhack
RUN git clone --recursive https://github.com/getpelican/pelican-plugins /tmp/pelican-plugins
WORKDIR /site
ENTRYPOINT ["make", "build"]
