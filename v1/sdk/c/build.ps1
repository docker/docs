rm -rf libssh
mkdir -p libssh
cd ./libssh
C:\Users\michaK\go\src\github.com\docker\pinata\win\bin\nuget.exe install openssl
cmake -DWITH_ZLIB=OFF -DOPENSSL_INCLUDE_DIR=C:\Users\michaK\go\src\github.com\docker\pinata\v1\sdk\c\libssh\openssl.v140.windesktop.msvcstl.dyn.rt-dyn.x86.1.0.2.0\build\native\include -DOPENSSL_ROOT_DIR=C:\Users\michaK\go\src\github.com\docker\pinata\v1\sdk\c\libssh\openssl.v140.windesktop.msvcstl.dyn.rt-dyn.x86.1.0.2.0\lib\native\v140\windesktop\msvcstl\dyn\rt-dyn\x86\release ../../vendor/libssh
cd ..