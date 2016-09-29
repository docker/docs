# compile dependencies
all: libssh

# compile libssh into ./libssh/build
libssh:
	mkdir -p ./AgentSdk/libssh
	# use cmake to generate the libssh project in ./libssh
	cd ./AgentSdk/libssh && \
	cmake -G"Unix Makefiles" \
	-DOPENSSL_INCLUDE_DIR=../../../vendor/openssl/include \
	-DOPENSSL_CRYPTO_LIBRARY=../../../vendor/openssl/lib/libcrypto.1.0.0.dylib \
	-DOPENSSL_SSL_LIBRARY=../../../vendor/openssl/lib/libssl.1.0.0.dylib \
	-DCMAKE_INSTALL_PREFIX=./build \
	../../../vendor/libssh
	# compile it and put output files in ./libssh/build
	cd ./AgentSdk/libssh && make install

# clean dependencies
clean:
	# remove generated libssh project
	rm -rf ./AgentSdk/libssh

