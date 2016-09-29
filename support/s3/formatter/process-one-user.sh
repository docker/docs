#!/bin/sh -e

IN=$1
OUT=$2

if [ ! -d "$1" ]; then
  echo INDIR not found
  echo Usage: $0 INDIR OUTDIR UUID
  exit 1
fi

if [ -z "$2" ]; then
  echo OUTDIR not specified
  echo Usage: $0 INDIR OUTDIR UUID
  exit 1
fi

if [ -z "$3" ]; then
  echo UUID not specified
  echo Usage: $0 INDIR OUTDIR UUID
  exit 1
fi

user=$3

mkdir -p ${OUT}
cat >${OUT}/index.html <<Header-fragment
<html>
<body>
<h1>Docker.App logs for user $user</h1>
<h2>Uploads by Date</h2>
<p><ul>
Header-fragment

cd ${IN}
for d in 20*; do
  /formatter/process-one-upload.sh ${IN}/$d ${OUT}/$d $d
  echo $d
  dt=`echo $d | awk -F- '{print $1}'`
  tm=`echo $d | awk -F- '{print $2}' | fold -w2 | paste -sd':' -`
  echo "<li><a href=\"$d/index.html\">$dt $tm</a></li>" >> ${OUT}/index.html
done

cat >>${OUT}/index.html <<Footer-fragment
</ul></p>
</body>
</html>
Footer-fragment

