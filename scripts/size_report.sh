#!/usr/bin/env sh

TARGET="$1"

if ! [ -d "$TARGET" ]; then
  echo "Target directory $TARGET does not exist. Exiting."
  exit 1
fi

echo -e "Uncompressed size:\n"
printf "HTML  : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.html" -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "JS    : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.js"   -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "CSS   : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.css"  -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "JSON  : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.json" -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "SVG   : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.svg"  -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "TXT   : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.txt"  -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "TOTAL : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -regex ".*\.\(html\|js\|css\|json\|svg\|txt\)" -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')

echo ""
echo -e "Compressed size:\n"
printf "HTML  : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.html.gz" -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "JS    : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.js.gz"   -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "CSS   : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.css.gz"  -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "JSON  : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.json.gz" -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "SVG   : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.svg.gz"  -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "TXT   : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.txt.gz"  -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
printf "TOTAL : "; numfmt --to=iec-i --padding=10 --suffix=B $(find ${TARGET} -name "*.gz"      -exec ls -l {} \; | awk '{ Total += $5} END { print Total }')
echo ""
