
.PHONY: depends clean display

display: output/docker-proxy-socket.png
	open output/docker-proxy-socket.png

output/docker-proxy-socket.png: depends output/docker-proxy-socket.gp output/docker-proxy-socket-upload.dat output/docker-proxy-socket-download.dat
	gnuplot output/docker-proxy-socket.gp

output/docker-proxy-socket-upload.dat output/docker-proxy-socket-download.dat output/docker-proxy-socket.gp: DESCRIBE

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
