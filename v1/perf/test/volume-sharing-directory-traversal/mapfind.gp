set terminal png size 1024,768
set output 'mapfind.png'
set title "osxfs directory traversal latency"
set xlabel "Test"
set ylabel "Time (s)"
set boxwidth 0.5
set style fill solid
set xtics rotate

plot 'mapfind.dat' using 2:xtic(1) with boxes
