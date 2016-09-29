#!/bin/sh

if [ $# = 0 ]; then
  echo "Usage:"
  echo "$0 <dir0> ... <dirN>"
  echo "  -- write a LICENSE file to stdout by combining the LICENSE.packagename"
  echo "     files from each of the supplied directories"
  exit 1
fi

for DIR in $*; do
  pushd $DIR > /dev/null
  for FILE in $(find . | grep "LICENSE\|COPYING" | grep -v ".skip$"); do
    echo "# Begin $FILE"
    if $(file $FILE | grep -q ISO-8859-1) ; then
      iconv -c -f ISO-8859-1 -t UTF-8 < $FILE
    else
      iconv -c -f UTF-8 -t UTF-8 < $FILE
    fi
    echo "# End $FILE"
    echo
  done
  popd > /dev/null
done
