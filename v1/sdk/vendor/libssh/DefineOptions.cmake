option(WITH_GSSAPI "Build with GSSAPI support" ON)
option(WITH_ZLIB "Build with ZLIB support" ON)
option(WITH_SSH1 "Build with SSH1 support" OFF)
option(WITH_SFTP "Build with SFTP support" ON)
option(WITH_SERVER "Build with SSH server support" ON)
option(WITH_STATIC_LIB "Build with a static library" OFF)
option(WITH_DEBUG_CRYPTO "Build with cryto debug output" OFF)
option(WITH_DEBUG_CALLTRACE "Build with calltrace debug output" ON)
option(WITH_GCRYPT "Compile against libgcrypt" OFF)
option(WITH_PCAP "Compile with Pcap generation support" ON)
option(WITH_INTERNAL_DOC "Compile doxygen internal documentation" OFF)
option(WITH_TESTING "Build with unit tests" OFF)
option(WITH_CLIENT_TESTING "Build with client tests; requires a running sshd" OFF)
option(WITH_BENCHMARKS "Build benchmarks tools" OFF)
option(WITH_EXAMPLES "Build examples" ON)
option(WITH_NACL "Build with libnacl (curve25519" ON)
if (WITH_ZLIB)
    set(WITH_LIBZ ON)
else (WITH_ZLIB)
    set(WITH_LIBZ OFF)
endif (WITH_ZLIB)

if(WITH_BENCHMARKS)
  set(WITH_TESTING ON)
endif(WITH_BENCHMARKS)

if (WITH_TESTING)
  set(WITH_STATIC_LIB ON)
endif (WITH_TESTING)

if (WITH_NACL)
  set(WITH_NACL ON)
endif (WITH_NACL)