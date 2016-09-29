set output 'mapfind.png'

set terminal png size 1024,768
set title "osxfs directory traversal latency"
set xlabel "Test"
set ylabel "Time (s)"
set boxwidth 0.5
set style fill solid
set xtics rotate

plot 'mapfind-nocache.dat' using 2:xtic(1) title "no dentry cache" with boxes,\
     'mapfind-cache.dat' using 2:xtic(1) title "10s lookup cache" with boxes
