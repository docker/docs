FROM fish/haproxy

ADD ./haproxy.cfg /haproxy/haproxy.cfg
ADD ./run /haproxy/run

EXPOSE 80 443

ENTRYPOINT ["/haproxy/run"]

