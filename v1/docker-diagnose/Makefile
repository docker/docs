VERSION=$(shell git rev-parse --short HEAD)
PACKAGES=fmt.tty re.glob re.str threads.posix stringext astring cmdliner \
  ezxmlm ezjsonm

OFLAGS=$(shell \
	ocamlfind query $(PACKAGES) -predicates native,link,mt,mt_posix,re.str -r \
	-format "-I %d %a") -w A-4-41-45-44 -warn-error +1..49

.PHONY: all clean depends src/version.ml

# Due to the behaviour of "go build" and the cp in the docker diagnose
# rule below building docker-diagnose-upload must always come strictly
# before building docker-diagnose.
all:
	$(MAKE) docker-diagnose

depends:
	opam install -y --color=always fmt re stringext cmdliner

.PHONY: src/version.ml

src/version.ml:
	echo "let git = \"$(VERSION)\"" > $@

clean:
	rm -f docker-diagnose *~
	rm -rf _tests
	rm -f src/*.cmi src/*.cmx src/*.o src/*~
	rm -f src/version.ml
	rm -f diagnose docker-diagnose-upload

docker-diagnose: src/misc.ml src/moby.ml src/diagnose.ml src/version.ml
	ocamlopt -linkall $(OFLAGS) -g -I src \
	  src/version.ml src/misc.ml src/common.ml src/user.ml src/logs.ml \
		src/driver.ml src/db.ml src/menubar.ml src/vmnetd.ml src/slirp.ml \
		src/env.ml src/osxfs.ml src/moby.ml src/dockerCli.ml src/system.ml \
		src/app.ml src/virtualization.ml src/disk.ml src/dns.ml \
		src/diagnose.ml -o docker-diagnose
