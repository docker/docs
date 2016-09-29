#!/bin/sh -e

IN=$1
OUT=$2

if [ ! -d "$1" ]; then
  echo Usage: $0 INDIR OUTDIR
  exit 1
fi

if [ -z "$2" ]; then
  echo Usage: $0 INDIR OUTDIR
  exit 1
fi

d=$3
dt=`echo $d | awk -F- '{print $1}'`
tm=`echo $d | awk -F- '{print $2}' | fold -w2 | paste -sd':' -`

rm -rf ${OUT}/*.html
mkdir -p ${OUT}
cd ${OUT}

tar --strip-components=2 -xf ${IN}/report.tar

cat >index.html <<Header-fragment
<html>
<body>
<h1>User `cat ${OUT}/user-id`</h1>
<p><a href="../index.html"><i>back to date index</i></a></p>
<p><pre>
Date:           $dt $tm
Version:        `cat ${OUT}/version`
`cat ${OUT}/sw_vers`
</pre></p>
<h2>Logs</h2>
<p><ul>
Header-fragment

FILES=`find . -type f`
echo ${FILES}
for i in `find . -type f`; do
  file=`echo $i | sed -e 's,^\./,,g'`
  case $i in
  *DS_Store)
    echo Skipping DS_Store
    ;;
  *.swp)
    echo Ignoring dot files
    ;;
  *.html)
    echo Skipping $i as HTML already
    ;;
  *)
    echo Processing $i
    pygmentize -f html -O full,encoding=latin1 -l text -o $i.html $i
    echo "<li><a href=\"$file.html\">$file</a></li>" >> index.html
    ;;
  esac
done

cat >>index.html <<Footer-fragment
</ul></p>
</body>
</html>
Footer-fragment
