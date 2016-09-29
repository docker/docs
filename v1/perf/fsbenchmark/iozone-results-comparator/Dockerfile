FROM scivm/scientific-python-2.7

RUN git clone https://github.com/cinterloper/iozone-results-comparator.git
WORKDIR /iozone-results-comparator
COPY requirements.txt /iozone-results-comparator/requirements.txt
RUN pip install -r requirements.txt
VOLUME /iotest
WORKDIR /iotest
ENTRYPOINT ["/iozone-results-comparator/src/iozone_results_comparator.py"]


