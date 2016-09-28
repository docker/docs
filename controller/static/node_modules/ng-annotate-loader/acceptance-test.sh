
# Simple test
npm run example-simple

cmp examples/simple/annotated-reference.js examples/simple/dist/build.js && echo 'Simple test passed' || exit 123

#Babel test
npm run example-babel

cmp examples/babel/annotated-reference.js examples/babel/dist/build.js && echo 'Babel test passed' || exit 123
cmp examples/babel/sourcemap-reference.js.map examples/babel/dist/build.js.map && echo 'Babel sourcemap test passed' || exit 123