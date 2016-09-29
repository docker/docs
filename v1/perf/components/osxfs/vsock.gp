Min = 0
Max = 500
n = 250

width = (Max - Min)/n
bin(x) = width*(floor((x-Min)/width)+0.5) + Min

set terminal png size 1024,768
set output 'vsock.png'
set xrange [Min:Max]
set title "osxfs-transfused vsock latency overhead"
set xlabel "μs"
set ylabel "Count"

set xtics 50
set xtics add ("50%%" (fifty / 1000))
set xtics add ("90%%" (ninety / 1000))

plot 'vsock.dat' using \
     (bin($1/1000)):(1.0) smooth freq with boxes title "vsock RTT"
