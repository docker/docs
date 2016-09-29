# The exact list is `ls packages/local` + $PKGS
PKGS = alcotest depext ocamlscript ezxmlm ezjsonm

MIRRORS=

update:
	./update-upstream.sh $(PKGS)
	$(MAKE) mirror

mirror:
	mkdir -p packages/upstream/oasis.0.4.7
	echo "mirrors: [\n\
	  \"http://gazagnaire.org/packages/oasis.0.4.7/oasis-0.4.7.tar.gz\"\n\
	]" >> packages/upstream/oasis.0.4.7/url
	mkdir -p packages/upstream/ocamlify.0.0.1
	echo "mirrors: [\n\
	  \"http://gazagnaire.org/packages/ocamlify.0.0.1/ocamlify-0.0.1.tar.gz\"\n\
	]" >> packages/upstream/ocamlify.0.0.1/url
	mkdir -p packages/upstream/ocamlmod.0.0.8
	echo "mirrors: [\n\
	  \"http://gazagnaire.org/packages/ocamlmod.0.0.8/ocamlmod-0.0.8.tar.gz\"\n\
	]" >> packages/upstream/ocamlmod.0.0.8/url


commit:
	$(MAKE) update
	git add packages/upstream
	git commit packages/upstream -s -m "AUTO: Update upstream packages"
