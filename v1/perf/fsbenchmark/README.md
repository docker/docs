# File system mount Performance benchmark

TL;DR
```
./bench.sh
```

The ./bench.sh script runs the benchmark for hyperkit local vs volume mounted and generates a html visual report of the comparison. bench.sh takes one argument, the max size of the files it tests. By default it is set at 10m. If you are interested in larger sized tests, which are much slower, set it explicitly with the -g option.
```
./bench.sh -g 1024m
```

If you add the -m `machinename` option it will run a third set for VirtualBox shared folders using the docker-machine provisioned machine named `machinename`. The machine needs to be up when you run.

```
./bench.sh -m default
```

If you use this option, the comparison of the 3 sets (native, vbox, pinata) will be in html_all, comparison of vbox vs pinata will be in html_out. The second one has more detailed data to understand the difference between the 2. The first one is useful to see how we perform compared to native.

## IOzone filesystem benchmark

We're using the [IOzone Filesystem Benchmark](http://iozone.org/) to test the performance of filesystems when they are mounted in containers (inspired by Mitchell H's http://mitchellh.com/comparing-filesystem-performance-in-virtual-machines). We are running iozone in a container, using the [threadx/iozone](https://hub.docker.com/r/threadx/iozone/) image. The iozone shell script is a convenience wrapper to call the container with iozone in automatic mode from 4k to 1Gb file sizes. All the options of iozone can be found in the [iozone reference doc](http://iozone.org/docs/IOzone_msword_98.pdf).

```
./iozone.sh
ls iotest
```

-o option lets you specify the the output file path on your local machine. The directory where this file is will be mounted in the container. -f allows to specify path for working file.
Here is an example where we run a baseline for performance inside virtualbox native fs:
```
./iozone.sh -o ~/iozone/baseline_1.iozone -f /f1
```

## iozone-results-comparator report generator

We use [iozone-results-comparator](https://github.com/cinterloper/iozone-results-comparator) ([docs](https://code.google.com/p/iozone-results-comparator/wiki/Tutorial) to produce a html visualization of the comparison. I created a Dockerfile for iozone-results-comparator. To build it yourself do:
```
docker build -t chanezon/iozone-rc .
```
Beware, it takes more than 15 minutes to build.

To generate a report from your test results:
```
cd iotest
docker run -ti -v $(pwd):/iotest chanezon/iozone-rc \
--baseline $(find ./ -name "baseline*.iozone" | sort) \
--set1 $(find ./ -name "sharedfolder*.iozone" | sort)
open iotest/html_out/index.html
```

Data analysis is left as an exercise to the reader:-)

## TODO
* build a leaner image based on python:2 or debian image
* try it in pinata with run
