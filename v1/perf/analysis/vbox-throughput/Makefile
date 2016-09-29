
.PHONY: depends clean display

display: output/plot.png
	open output/plot.png

output/plot.png: depends output/plot.gp output/plot.normal
	gnuplot output/plot.gp

output/plot.normal output/plot.gp: DESCRIBE

.PHONY: DESCRIBE
DESCRIBE: depends filter.sh describe.ml
	rm -f describe.ml.exe
	./filter.sh
	./describe.ml output/*.dat

depends:
	which gnuplot || brew install gnuplot
	which opam || brew install opam
	(eval `opam config env` && which ocamlscript) || (opam install ocamlscript stringext -y|| (opam init -y && opam install ocamlscript stringext -y))

clean:
	rm -rf output
