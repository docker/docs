GIT_VERSION := $(shell git describe --abbrev=6 --dirty --always --tags)
GIT_VERSION_SHA1 := $(shell git rev-parse HEAD)

ifeq ($V, 1)
	VERBOSE =
else
	VERBOSE = @
endif

include config.mk

VMM_SRC := \
	src/vmm/x86.c \
	src/vmm/vmm.c \
	src/vmm/vmm_host.c \
	src/vmm/vmm_mem.c \
	src/vmm/vmm_lapic.c \
	src/vmm/vmm_instruction_emul.c \
	src/vmm/vmm_ioport.c \
	src/vmm/vmm_callout.c \
	src/vmm/vmm_stat.c \
	src/vmm/vmm_util.c \
	src/vmm/vmm_api.c \
	src/vmm/intel/vmx.c \
	src/vmm/intel/vmx_msr.c \
	src/vmm/intel/vmcs.c \
	src/vmm/io/vatpic.c \
	src/vmm/io/vatpit.c \
	src/vmm/io/vhpet.c \
	src/vmm/io/vioapic.c \
	src/vmm/io/vlapic.c \
	src/vmm/io/vpmtmr.c \
	src/vmm/io/vrtc.c

XHYVE_SRC := \
	src/acpitbl.c \
	src/atkbdc.c \
	src/block_if.c \
	src/consport.c \
	src/dbgport.c \
	src/inout.c \
	src/ioapic.c \
	src/md5c.c \
	src/mem.c \
	src/mevent.c \
	src/mptbl.c \
	src/pci_ahci.c \
	src/pci_emul.c \
	src/pci_hostbridge.c \
	src/pci_irq.c \
	src/pci_lpc.c \
	src/pci_uart.c \
	src/pci_virtio_9p.c \
	src/pci_virtio_block.c \
	src/pci_virtio_net_tap.c \
	src/pci_virtio_net_vmnet.c \
	src/pci_virtio_net_vpnkit.c \
	src/pci_virtio_rnd.c \
	src/pci_virtio_sock.c \
	src/pm.c \
	src/post.c \
	src/rtc.c \
	src/smbiostbl.c \
	src/task_switch.c \
	src/uart_emul.c \
	src/xhyve.c \
	src/virtio.c \
	src/xmsr.c

FIRMWARE_SRC := \
	src/firmware/bootrom.c \
	src/firmware/kexec.c \
	src/firmware/fbsd.c

HAVE_OCAML_QCOW := $(shell if ocamlfind query qcow uri >/dev/null 2>/dev/null ; then echo YES ; else echo NO; fi)

ifeq ($(HAVE_OCAML_QCOW),YES)
CFLAGS += -DHAVE_OCAML=1 -DHAVE_OCAML_QCOW=1 -DHAVE_OCAML=1

# prefix vsock file names if PRI_ADDR_PREFIX
# is defined. (not applied to aliases)
ifneq ($(PRI_ADDR_PREFIX),)
CFLAGS += -DPRI_ADDR_PREFIX=\"$(PRI_ADDR_PREFIX)\"
endif

# override default connect socket name if 
# CONNECT_SOCKET_NAME is defined 
ifneq ($(CONNECT_SOCKET_NAME),)
CFLAGS += -DCONNECT_SOCKET_NAME=\"$(CONNECT_SOCKET_NAME)\"
endif

OCAML_SRC := \
	src/mirage_block_ocaml.ml

OCAML_C_SRC := \
	src/mirage_block_c.c

OCAML_WHERE := $(shell ocamlc -where)
OCAML_PACKS := cstruct cstruct.lwt io-page io-page.unix uri mirage-block mirage-block-unix qcow unix threads lwt lwt.unix
OCAML_LDLIBS := -L $(OCAML_WHERE) \
	$(shell ocamlfind query cstruct)/cstruct.a \
	$(shell ocamlfind query cstruct)/libcstruct_stubs.a \
	$(shell ocamlfind query io-page)/io_page.a \
	$(shell ocamlfind query io-page)/io_page_unix.a \
	$(shell ocamlfind query io-page)/libio_page_unix_stubs.a \
	$(shell ocamlfind query lwt.unix)/liblwt-unix_stubs.a \
	$(shell ocamlfind query lwt.unix)/lwt-unix.a \
	$(shell ocamlfind query lwt.unix)/lwt.a \
	$(shell ocamlfind query threads)/libthreadsnat.a \
	$(shell ocamlfind query mirage-block-unix)/libmirage_block_unix_stubs.a \
	-lasmrun -lbigarray -lunix

build/xhyve.o: CFLAGS += -I$(OCAML_WHERE)
endif

SRC := \
	$(VMM_SRC) \
	$(XHYVE_SRC) \
	$(FIRMWARE_SRC) \
	$(OCAML_C_SRC)

OBJ := $(SRC:src/%.c=build/%.o) $(OCAML_SRC:src/%.ml=build/%.o)
DEP := $(OBJ:%.o=%.d)
INC := -Iinclude

CFLAGS += -DVERSION=\"$(GIT_VERSION)\" -DVERSION_SHA1=\"$(GIT_VERSION_SHA1)\"

TARGET = build/com.docker.hyperkit

all: $(TARGET) | build

.PHONY: clean all test
.SUFFIXES:

-include $(DEP)

build:
	@mkdir -p build

include/xhyve/dtrace.h: src/dtrace.d
	@echo gen $<
	$(VERBOSE) $(DTRACE) -h -s $< -o $@

$(SRC): include/xhyve/dtrace.h

build/%.o: src/%.c
	@echo cc $<
	@mkdir -p $(dir $@)
	$(VERBOSE) $(ENV) $(CC) $(CFLAGS) $(INC) $(DEF) -MMD -MT $@ -MF build/$*.d -o $@ -c $<

$(OCAML_C_SRC:src/%.c=build/%.o): CFLAGS += -I$(OCAML_WHERE)
build/%.o: src/%.ml
	@echo ml $<
	@mkdir -p $(dir $@)
	$(VERBOSE) $(ENV) ocamlfind ocamlopt -thread -package "$(OCAML_PACKS)" -c $< -o build/$*.cmx
	$(VERBOSE) $(ENV) ocamlfind ocamlopt -thread -linkpkg -package "$(OCAML_PACKS)" -output-obj -o $@ build/$*.cmx

$(TARGET).sym: $(OBJ)
	@echo ld $(notdir $@)
	$(VERBOSE) $(ENV) $(LD) $(LDFLAGS) -Xlinker $(TARGET).lto.o -o $@ $(OBJ) $(LDLIBS) $(OCAML_LDLIBS)
	@echo dsym $(notdir $(TARGET).dSYM)
	$(VERBOSE) $(ENV) $(DSYM) $@ -o $(TARGET).dSYM

$(TARGET): $(TARGET).sym
	@echo strip $(notdir $@)
	$(VERBOSE) $(ENV) $(STRIP) $(TARGET).sym -o $@

clean:
	@rm -rf build
	@rm -f include/xhyve/dtrace.h
	@rm -f test/vmlinuz test/initrd.gz

test/vmlinuz test/initrd.gz:
	@cd test; ./tinycore.sh

test: $(TARGET) test/vmlinuz test/initrd.gz
	@./test_linux.exp
