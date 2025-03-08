
SIZE="1280px"
HOST=dtr.domain.com
HOSTIP:=192.168.18.131
HOSTURL=https://$(HOST)
IMAGE=svendowideit/screenshot
#IMAGE=screenshot
RUNOPTS=docker run --rm -v $(CURDIR):/srv --add-host $(HOST):$(HOSTIP)
RUN=$(RUNOPTS) $(IMAGE)

all:
	# Needs to be run with no authentication
	$(RUN) $(HOSTURL)/admin/ admin-metrics.png $(SIZE)
	$(RUN) $(HOSTURL)/admin/logs/ admin-logs.png $(SIZE)
	$(RUN) $(HOSTURL)/admin/settings/auth/ admin-settings-authentication.png $(SIZE)
	$(RUN) $(HOSTURL)/admin/settings/general admin-settings.png $(SIZE)
	$(RUN) $(HOSTURL)/admin/settings/security admin-settings-security.png $(SIZE)
	$(RUN) $(HOSTURL)/admin/settings/storage admin-settings-storage.png $(SIZE)
	$(RUN) $(HOSTURL)/admin/settings/license admin-settings-license.png $(SIZE)

manual:
	#  --env-file=passwords.env
	$(RUNOPTS)  --env-file=passwords.env \
		$(IMAGE) \
		$(HOSTURL)/admin/settings#auth admin-settings-authentication-basic.png $(SIZE)

	#admin-settings-authentication-ldap.png


# TODO: needs manual intervention atm
qs1:
	# this requires the hostname to resolve.
	echo "$HOST $HOSTIP" >> /etc/hosts
	# and requires that host to be added to the docker daemon's insecure registry settings
	# maybe use a dind daemon?
	
	docker pull jenkins
	docker tag -f jenkins t400s.home.org.au/ci-infrastructure/jnkns-img
	docker push t400s.home.org.au/ci-infrastructure/jnkns-img
	sleep 1
	$(RUN) $(HOSTURL)/admin/ console-push.png $(SIZE)
	sleep 10
qs2:
	docker rmi jenkins t400s.home.org.au/ci-infrastructure/jnkns-img
	docker pull t400s.home.org.au/ci-infrastructure/jnkns-img
	sleep 1
	$(RUN) $(HOSTURL)/admin/ console-pull.png $(SIZE)
