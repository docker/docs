.PHONY: docs osx windows

default: osx windows
	@true

clean: clean-osx clean-windows
	@true

osx: clean-osx
	./script/build-osx

windows: clean-windows
	./script/build-windows

clean-osx:
	rm -f dist/DockerToolbox-*.pkg

clean-windows:
	rm -f dist/DockerToolbox-*.exe

docs:
	$(MAKE) -C docs docs
