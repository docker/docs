.PHONY: all depends build configure clean distclean

all: build

QEMUV ?= 2.4.1

qemu-$(QEMUV):
	git clone -b v${QEMUV} https://github.com/qemu/qemu.git qemu-${QEMUV}

configure: qemu-$(QEMUV) qemu-$(QEMUV)/x86_64-softmmu/config-target.mak
qemu-$(QEMUV)/x86_64-softmmu/config-target.mak:
	cd qemu-$(QEMUV) \
	&& ./configure --target-list=x86_64-softmmu \
	     --disable-cocoa --disable-gtk --disable-vnc

depends: configure qemu-$(QEMUV)/x86_64-softmmu/qemu-system-x86_64
	cp qemu-$(QEMUV)/x86_64-softmmu/qemu-system-x86_64 $(CACHE_DIR)/

qemu-$(QEMUV)/x86_64-softmmu/qemu-system-x86_64:
	make -C qemu-$(QEMUV)

build:
	cp $(CACHE_DIR)/qemu-system-x86_64 .

clean:
	make -C qemu-$(QEMUV) clean
	$(RM) -r qemu-system-x86_64

distclean:
	$(RM) -r qemu-$(QEMUV)
