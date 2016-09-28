# Temporary image built by the UCP CI system with docker 1.12 experimental CLI included in the image
FROM dockerorcadev/compose:latest

MAINTAINER Tom Barlow "<thomas.barlow@docker.com>"

ADD compose-wrapper.sh /bin/compose-wrapper.sh
ADD bundle-wrapper.sh /bin/bundle-wrapper.sh
RUN chmod +x /bin/compose-wrapper.sh && \
	chmod +x /bin/bundle-wrapper.sh && \
    mkdir /ucp-compose

WORKDIR /ucp-compose

ENTRYPOINT ["compose-wrapper.sh"]
