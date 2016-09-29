FROM ubuntu:16.04
RUN apt-get update && apt-get install -y g++
RUN mkdir /src
COPY Main.cpp /src/Main.cpp
RUN g++ /src/Main.cpp -std=c++14 -o3 -pthread -o /PerfTest
