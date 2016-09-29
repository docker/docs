#!/bin/bash
while getopts ":g:m:" opt; do
  case $opt in
    g)
      max_size=$OPTARG
      ;;
    m)
      machine_name=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      echo "Syntax $0 -g max_size" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done
echo $max_size
if [ ! -d "iotest" ]; then
  mkdir iotest
else
  rm -rf iotest
  mkdir iotest
fi
if  [ ! -n "$max_size" ]; then
	max_size=10m
fi
set -o verbose
echo ${max_size}
./iozone.sh -o $(pwd)/iotest/baseline_1.iozone -f /f1 -g ${max_size}
./iozone.sh -o $(pwd)/iotest/baseline_2.iozone -f /f1 -g ${max_size}
./iozone.sh -o $(pwd)/iotest/volume_1.iozone -g ${max_size}
./iozone.sh -o $(pwd)/iotest/volume_2.iozone -g ${max_size}
if  [ -n "$machine_name" ]; then
	eval $(docker-machine env $machine_name)
	./iozone.sh -o $(pwd)/iotest/sharedfolder_1.iozone -g ${max_size}
	./iozone.sh -o $(pwd)/iotest/sharedfolder_2.iozone -g ${max_size}
	cd iotest
	multiset="--set2 $(find ./ -name "sharedfolder*.iozone" | sort) --multiset"
	cd ..
else
	multiset=""
fi
cd iotest
docker run -i -v $(pwd):/iotest chanezon/iozone-rc \
--baseline $(find ./ -name "baseline*.iozone" | sort) \
--set1 $(find ./ -name "volume*.iozone" | sort) \
$multiset
if  [ -n "$machine_name" ]; then
	mv html_out html_all
	docker run -i -v $(pwd):/iotest chanezon/iozone-rc \
	--baseline $(find ./ -name "sharedfolder*.iozone" | sort) \
	--set1 $(find ./ -name "volume*.iozone" | sort)
fi
