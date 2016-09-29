FROM ibmcom/db2express-c:latest

ENV DB2INST1_PASSWORD=password
ENV LICENSE=accept
USER root
RUN /bin/bash -c 'su - db2inst1 -c "db2start && db2 create database test"'
