# compile dependencies and C SDK
all: dep-libssh libpinata

# clean C SDK and dependencies
clean:
	# remove generated libssh project
	rm -rf ./libssh
	rm -rf ./build

### PRIVATE ###

# compile C sdk as a static lib
libpinata:
	#clang ./sdk/main.c ./sdk/sdk.c -c -o ./sdk.o -Ilibssh/build/include -Llibssh/build/lib -lssh
	mkdir -p ./build
	clang ./sdk/sdk.c -c -Ilibssh/build/include -o ./build/sdk.o
	ar rcs ./build/pinata_sdk.a ./build/sdk.o
	rm -f ./build/sdk.o
	cp ./sdk/sdk.h ./build/pinata_sdk.h

# compile libssh into ./libssh/build
dep-libssh:
	mkdir -p libssh
	# use cmake to generate the libssh project in ./libssh
	cd ./libssh && \
	cmake -G"Unix Makefiles" \
	-DOPENSSL_INCLUDE_DIR=../../vendor/openssl/include \
	-DOPENSSL_CRYPTO_LIBRARY=../../vendor/openssl/lib/libcrypto.1.0.0.dylib \
	-DOPENSSL_SSL_LIBRARY=../../vendor/openssl/lib/libssl.1.0.0.dylib \
	-DCMAKE_INSTALL_PREFIX=./build \
	../../vendor/libssh
	# compile it and put output files in ./libssh/build
	cd ./libssh && make install
