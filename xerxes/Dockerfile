FROM python:3.9-slim

RUN apt-get update \
&& apt-get install gcc -y \
&& apt-get clean

# set working directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# add requirements (to leverage Docker cache)
ADD ./requirements.txt /usr/src/app/requirements.txt

# install requirements
RUN pip install -r requirements.txt

# add app
ADD . /usr/src/app/

ENTRYPOINT ["tail", "-f", "/dev/null"]
