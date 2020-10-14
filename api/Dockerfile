FROM python:3.8-alpine
# set working directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# add requirements (to leverage Docker cache)
ADD ./requirements.txt /usr/src/app/requirements.txt

# install requirements
RUN pip install -r requirements.txt

EXPOSE 5000

# add app
ADD app.py /usr/src/app
ADD error_model.py /usr/src/app

RUN mkdir -p /usr/src/app/api
ADD misc /usr/src/app/api
ADD verba /usr/src/app/verba
ADD nomina /usr/src/app/nomina

# run server
#CMD python manage.py runserver -h 0.0.0.0
CMD gunicorn --bind 0.0.0.0:5000 app