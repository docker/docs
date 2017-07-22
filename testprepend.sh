svn co https://github.com/docker-library/docs/trunk _samples/library || (echo "Failed library download" && exit -1)
FILES=$(find _samples/library -type f -name 'README.md')
for f in $FILES
do
  echo "Adding empty front-matter to $f ..."
  sed -i '1i ---\nlayout: library\ntitle: $f\nkeywords:library, sample, $f\n---' $f
  # take action on each file.
done
