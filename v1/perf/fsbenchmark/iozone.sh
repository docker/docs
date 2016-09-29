#!/bin/bash
while getopts ":o:f:g:" opt; do
  case $opt in
    o)
      output_file=$OPTARG
      ;;
    f)
      test_file=$OPTARG
      ;;
    g)
      max_size=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      echo "Syntax $0 -o output_file -f test_file" >&2
      echo "Example: $0 -o $(pwd)/iozone/results.iozone -f /f1" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done
if  [ ! -n "$output_file" ]; then
	if [ ! -d "iotest" ]; then
	  mkdir iotest
	fi
	output_dir=$(pwd)/iotest
	output_file=${output_dir}/results.iozone
else
	output_dir=`dirname ${output_file}`
fi
if  [ ! -n "$test_file" ]; then
	test_file=/iotest/f1
fi
if  [ ! -n "$max_size" ]; then
	max_size=10m
fi
docker run -i -v ${output_dir}:/iotest threadx/docker-ubuntu-iozone iozone -a -n 4k -q 31k -g ${max_size} -f ${test_file} > ${output_file}
