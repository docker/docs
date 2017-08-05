rm -rf _samples/library
svn co https://github.com/docker-library/docs/trunk _samples/library || (echo "Failed library download" && exit -1)

find _samples/library -maxdepth 9  -type l -delete

FILES=$(find _samples/library -type f -name 'README.md')
for f in $FILES
do
  curdir=$(dirname "${f}")
  justcurdir="${curdir##*/}"
  echo "Adding front-matter to ${f} ..."
  echo --- >> ${curdir}/front-matter.txt
  echo title: "${justcurdir}" >> ${curdir}/front-matter.txt
  echo keywords: library, sample, ${justcurdir} >> ${curdir}/front-matter.txt
  echo repo: "${justcurdir}" >> ${curdir}/front-matter.txt
  echo layout: docs >> ${curdir}/front-matter.txt
  echo permalink: /samples/library/${justcurdir}/ >> ${curdir}/front-matter.txt
  echo --- >> ${curdir}/front-matter.txt
  if [ -e ${curdir}/README-short.txt ]
  then
    shortrm=$(<${curdir}/README-short.txt)
    echo >> ${curdir}/front-matter.txt
    echo ${shortrm} >> ${curdir}/front-matter.txt
    echo >> ${curdir}/front-matter.txt
  fi
  if [ -e ${curdir}/github-repo ]
  then
    gitrepo=$(<${curdir}/github-repo)
    echo >> ${curdir}/front-matter.txt
    echo GitHub repo: \["${gitrepo}"\]\("${gitrepo}"\)\{: target="_blank"\} >> ${curdir}/front-matter.txt
    echo >> ${curdir}/front-matter.txt
  fi
  cat ${curdir}/front-matter.txt _samples/boilerplate.txt > ${curdir}/header.txt
  cat ${curdir}/header.txt ${curdir}/README.md > ${curdir}/index.md
  rm -rf ${curdir}/front-matter.txt
  rm -rf ${curdir}/header.txt
  rm -rf _samples/library/index.md
  #sed -i '1i ---\nlayout: library\ntitle: ${f}\nkeywords: library, sample, ${f}\n---' $f
  # take action on each file.
done
