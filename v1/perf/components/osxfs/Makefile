.PHONY: vsock.dat ping.dat all

LOGS=~/Library/Containers/com.docker.docker/Data/logs
DIR=$(shell ls $(LOGS) | sort -r | head -1)
VSOCK=$(shell ls $(LOGS)/$(DIR)/*.vsock.ping | sort -r | head -1)
PING=$(shell ls $(LOGS)/$(DIR)/*.ping | grep -v vsock | sort -r | head -1)

all: vsock.png ping.png

vsock.dat:
	ln -sf $(VSOCK) vsock.dat

vsock.png: vsock.dat vsock.gp
	gnuplot \
	-e "fifty=`../../lib/percentiles.py vsock.dat 50`" \
	-e "ninety=`../../lib/percentiles.py vsock.dat 90`" \
	vsock.gp

ping.dat:
	ln -sf $(PING) ping.dat

ping.png: ping.dat ping.gp vsock.dat
	gnuplot \
	-e "vsock=`../../lib/percentiles.py vsock.dat 50`" \
	ping.gp
