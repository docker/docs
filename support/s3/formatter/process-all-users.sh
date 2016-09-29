#!/usr/bin/env bash

set -ex

IN=$1
OUT=$2

if [ ! -d "$1" ]; then
  echo INDIR not found
  echo Usage: $0 INDIR OUTDIR
  exit 1
fi

if [ ! -d "$2" ]; then
  echo OUTDIR not found
  echo Usage: $0 INDIR OUTDIR
  exit 1
fi

mkdir -p ${OUT}
cat >${OUT}/index.html <<Header-fragment
<html>
<body>
<h1>Docker.App support logs</h1>
<h2>Upload By UUID</h2>
<p><ul>
Header-fragment

cd ${IN}
for uuid in *; do
  bash /formatter/process-one-user.sh ${IN}/$uuid ${OUT}/$uuid $uuid
  echo "<li><a href=\"$uuid/index.html\">$uuid</a></li>" >> ${OUT}/index.html
done

cat >>${OUT}/index.html <<Footer-fragment
</ul></p>
</body>
</html>
Footer-fragment
