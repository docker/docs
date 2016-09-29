
.PHONY: depends clean display

display: output/volume-sharing-dd-read.png output/volume-sharing-dd-write.png
	open output/volume-sharing-dd-read.png
	open output/volume-sharing-dd-write.png

output/volume-sharing-dd-read.png: depends output/volume-sharing-dd-read.gp output/volume-sharing-dd-read.dat
	gnuplot output/volume-sharing-dd-read.gp

output/volume-sharing-dd-write.png: depends output/volume-sharing-dd-write.gp output/volume-sharing-dd-write.dat
	gnuplot output/volume-sharing-dd-write.gp

output/volume-sharing-dd-read.dat output/volume-sharing-dd-read.gp output/volume-sharing-dd-write.dat output/volume-sharing-dd-write.gp: DESCRIBE

.PHONY: DESCRIBE
DESCRIBE: depends describe.ml
	rm -f describe.ml.exe
	./describe.ml ../../results

depends:
	which gnuplot || brew install gnuplot
	which opam || brew install opam
	(eval `opam config env` && which ocamlscript) || (opam install ocamlscript stringext -y|| (opam init -y && opam install ocamlscript stringext -y))

clean:
	rm -rf output
