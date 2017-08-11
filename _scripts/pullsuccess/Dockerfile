FROM python:2.7-slim

WORKDIR /app

ADD . /app

RUN pip install --no-cache-dir -r requirements.txt

CMD [ "python", "pullsuccess.py" ]
