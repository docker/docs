rm -rf _samples/library
svn co https://github.com/docker-library/docs/trunk _samples/library || (echo "Failed library download" && exit -1)

FILES=$(find _samples/library -type f -name 'README.md')
for f in $FILES
do
  curdir=$(dirname "${f}")
  justcurdir="${curdir##*/}"
  echo "Adding front-matter to ${f} ..."
  echo --- >> $(dirname "${f}")/front-matter.txt
  echo title: "${justcurdir}" >> $(dirname "${f}")/front-matter.txt
  echo keywords: library, sample, ${justcurdir} >> $(dirname "${f}")/front-matter.txt
  echo layout: library >> $(dirname "${f}")/front-matter.txt
  echo --- >> $(dirname "${f}")/front-matter.txt
  cat $(dirname "${f}")/front-matter.txt $(dirname "${f}")/README.md > $(dirname "${f}")/index.md
  rm -rf $(dirname "${f}")/front-matter.txt
  #sed -i '1i ---\nlayout: library\ntitle: ${f}\nkeywords: library, sample, ${f}\n---' $f
  # take action on each file.
done
