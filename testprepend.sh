rm -rf _samples/library
svn co https://github.com/docker-library/docs/trunk _samples/library || (echo "Failed library download" && exit -1)

FILES=$(find _samples/library -type f -name 'README.md')
for f in $FILES
do
  echo "Adding empty front-matter to ${f} ..."
  echo --- >> front-matter.txt
  echo title: $f >> front-matter.txt
  echo keywords: library, sample, $f >> front-matter.txt
  echo layout: library >> front-matter.txt
  echo --- >> front-matter.txt
  cat front-matter.txt README.md > index.md
  #sed -i '1i ---\nlayout: library\ntitle: ${f}\nkeywords: library, sample, ${f}\n---' $f
  # take action on each file.
done
