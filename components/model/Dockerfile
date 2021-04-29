FROM ubuntu:16.04

RUN apt-get update \ 
    && apt-get install -y curl \ 
    python3 \
    unzip \
    wget 

RUN wget https://bootstrap.pypa.io/get-pip.py \
    && python3 get-pip.py

RUN pip3 install opencv-python

RUN apt-get install -y  libglib2.0-dev

RUN apt-get install -y build-essential libgtk2.0-dev libavcodec-dev libavformat-dev libjpeg.dev libtiff4.dev libswscale-dev libjasper-dev  

RUN pip3 install pillow tensorflow==1.3
WORKDIR /

COPY  main main
COPY config.yaml config.yaml
COPY detect_rate /data/detect_rate
RUN mkdir /data/go
ENV JXS_CONFIG /

ENTRYPOINT ["./main"]
