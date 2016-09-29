.PHONY: all clean nurse

all: nurse
	@

clean:
	@ocamlbuild -clean
	@rm -f nurse *~

nurse:
	@ocamlbuild -use-ocamlfind 'nurse.native'
	@mv nurse.native nurse
