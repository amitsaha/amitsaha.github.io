FROM fedora:24
RUN dnf -y install python3-pip wget git make bash
RUN pip3 install pelican pelican-youtube
RUN git clone https://github.com/getpelican/pelican-themes.git /pelican-themes
RUN git config --global user.name "Amit Saha"
RUN git config --global user.email "amitsaha.in@gmail.com"
WORKDIR /site
ADD publish.sh /publish.sh
ENTRYPOINT ["/publish.sh"]
