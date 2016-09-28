FROM dockerautomation/dind-make:latest
MAINTAINER Mike Dougherty <mike.dougherty@docker.com>

# For dind, default stdout to /var/log/docker.log
ENV LOG=file
WORKDIR /usr/src/app
ENTRYPOINT ["/usr/src/app/docker-entrypoint.sh"]

COPY . /usr/src/app
