Min = 0
Max = 500
n = 250

width = (Max - Min)/n
bin(x) = width*(floor((x-Min)/width)+0.5) + Min

set terminal png size 1024,768
set output 'ping.png'
set xrange [Min:Max]
set title "osxfs latency overhead"
set xlabel "μs"
set ylabel "Count"

set xtics 50

plot 'ping.dat' using \
     (bin($1/1000)):(1.0) smooth freq with boxes title "event RTT",\
     'ping.dat' using \
     (bin($2/1000)):(1.0) smooth freq with boxes title "error RTT",\
     'ping.dat' using \
     (bin(($1+$2)/1000)):(1.) smooth freq with boxes title "event+error RTT",\
     'ping.dat' using \
     (bin(($1 + $2 - (2 * vsock))/1000)):(1.0) \
     smooth freq with boxes title "FUSE symlink error RTT"
