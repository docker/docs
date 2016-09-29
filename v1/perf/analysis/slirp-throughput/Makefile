
.PHONY: depends clean display

display: output/plot.png
	open output/plot.png

output/plot.png: depends output/plot.gp
	gnuplot output/plot.gp

output/plot.gp: run.sh
	./run.sh

depends:
	which gnuplot || brew install gnuplot

clean:
	rm -rf output
