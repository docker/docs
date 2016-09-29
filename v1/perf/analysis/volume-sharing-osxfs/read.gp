
  set terminal png
  set output 'output/read.png'
  set title 'volume sharing read throughput comparison'
  set ylabel 'Throughput in MiB/s'
  set xlabel 'Released software version'

  set xrange [0: 12]
  set xtics rotate by -45 ("alpha 7" 0, "alpha 8" 1, "alpha 9" 2, "alpha 10" 3, "alpha 11" 4, "alpha 12" 5, "beta 1" 6, "beta 2" 7, "virtualbox" 8, "fusion" 9)
  plot 'output/volume-sharing-dd-read.dat' using 0:1 title '512' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:2 title '1K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:3 title '2K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:4 title '4K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:5 title '8K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:6 title '16K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:7 title '32K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:8 title '64K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:9 title '128K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:10 title '256K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:11 title '512K' linewidth 2 with points, \
       'output/volume-sharing-dd-read.dat' using 0:12 title '1M' linewidth 2 with points
