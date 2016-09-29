.PHONY: all clean com.docker.osxfs

all: com.docker.osxfs
	@

clean:
	ocamlbuild -clean
	rm -f com.docker.osxfs *~
	rm -rf _tests

com.docker.osxfs:
	ocamlbuild -use-ocamlfind -classic-display osxfs.native
	cp osxfs.native com.docker.osxfs
