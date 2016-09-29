.PHONY: all build clean
all: build

build: \
    machines/test-osx-10.10-template \
    machines/test-osx-10.11-template \
    machines/test-osx-10.12-template \
    machines/test-win-10-pro-10586-template \
    machines/test-win-10-ent-10586-template \
    machines/test-win-10-pro-14393-template
clean:
	rm -r machines/*-template

get-isos:
	aws s3 sync s3://pinata-test-isos/isos iso

machines/test-osx-10.10-template:
	packer build -var-file=packer/osx-10.10.json packer/osx.json

machines/test-osx-10.11-template:
	packer build -var-file=packer/osx-10.11.json packer/osx.json

machines/test-osx-10.12-template:
	packer build -var-file=packer/osx-10.12.json packer/osx.json

machines/test-win-10-pro-10586-template:
	packer build -var-file=packer/win-10-pro-10586.json packer/win-10.json

machines/test-win-10-ent-10586-template:
	packer build -var-file=packer/win-10-ent-10586.json packer/win-10.json

machine/test-win-10-pro-14393-template:
	packer build -var-file=packer/win-10-pro-14393.json packer/win-10.json
