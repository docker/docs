FROM python:2.7-slim

WORKDIR /app

ADD . /app

RUN pip install --no-cache-dir -r requirements.txt

# TO RUN, cd to root of GitHub repo and run this:
# docker run -it --rm --name pullsuccess -v "$PWD":/app -w /app docs/pullsuccess python pullsuccess.py

CMD [ "python", "pullsuccess.py" ]
