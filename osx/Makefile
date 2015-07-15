DOCKER_IMAGE := osx-installer
DOCKER_CONTAINER := build-osx-installer

default: dockerbuild
	@true

dockerbuild: clean
	docker build -t $(DOCKER_IMAGE) .
	docker run --name "$(DOCKER_CONTAINER)" "$(DOCKER_IMAGE)"
	docker cp "$(DOCKER_CONTAINER)":/DockerToolbox.pkg .
	docker rm "$(DOCKER_CONTAINER)" 2>/dev/null || true
	productsign --sign "3rd Party Mac Developer Installer: Docker Inc" DockerToolbox.pkg DockerToolboxSigned.pkg || true
	mv DockerToolboxSigned.pkg DockerToolbox.pkg || true

clean:
	rm -f DockerToolbox.pkg
	docker rm "$(DOCKER_CONTAINER)" 2>/dev/null || true

