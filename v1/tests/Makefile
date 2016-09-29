PACKAGES=fmt alcotest re.glob threads.posix ezjsonm stringext
LWT_PACKAGES=lwt.unix lwt cstruct.syntax sexplib.syntax cstruct.lwt ipaddr io-page

PACKAGES_=$(PACKAGES:%=-package %)
LWT_PACKAGES_=$(LWT_PACKAGES:%=-package %)

.PHONY: all clean depends

all: pinata-rt
	@

depends:
	opam install -y --color=always alcotest fmt re ezjsonm stringext

clean:
	rm -f pinata-rt *~
	rm -rf _tests
	rm -f src/*.cmi src/*.cmx src/*.o src/*~

pinata-rt: install-symlinks-helper src/test.ml
	ocamlfind ocamlopt $(PACKAGES_) -g src/test.ml -linkpkg -thread -o pinata-rt

install-symlinks-helper: src/install_symlinks_helper.ml
	ocamlfind ocamlopt $(LWT_PACKAGES_) -syntax camlp4o -g src/install_symlinks_helper.ml -linkpkg -thread -o install-symlinks-helper
