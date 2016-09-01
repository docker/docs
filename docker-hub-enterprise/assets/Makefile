
RUN=docker run --rm -v $(CURDIR):/srv docker-screenshot
#  https://10.10.10.20/admin/metrics/ admin-metrics.png 
SIZE="984px*616px"

all:
	$(RUN) https://10.10.10.20/admin/metrics/ admin-metrics.png $(SIZE)
	$(RUN) https://10.10.10.20/admin/logs/ admin-logs.png $(SIZE)
	$(RUN) https://10.10.10.20/admin/settings/ admin-settings.png $(SIZE)
# sadly, the reactjs tabs we're using make it fall over.
#	$(RUN) https://10.10.10.20/admin/settings#security admin-settings-security.png $(SIZE)
