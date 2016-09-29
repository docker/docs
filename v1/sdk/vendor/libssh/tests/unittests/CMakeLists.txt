project(unittests C)

include_directories(${OPENSSL_INCLUDE_DIR})

add_cmocka_test(torture_buffer torture_buffer.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_callbacks torture_callbacks.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_init torture_init.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_list torture_list.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_misc torture_misc.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_options torture_options.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_isipaddr torture_isipaddr.c ${TORTURE_LIBRARY})
add_cmocka_test(torture_pki_ed25519 torture_pki_ed25519.c ${TORTURE_LIBRARY})
if (UNIX AND NOT WIN32)
    # requires ssh-keygen
    add_cmocka_test(torture_keyfiles torture_keyfiles.c ${TORTURE_LIBRARY})
    add_cmocka_test(torture_pki torture_pki.c ${TORTURE_LIBRARY})
    # requires /dev/null
    add_cmocka_test(torture_channel torture_channel.c ${TORTURE_LIBRARY})
    # requires pthread
    if (LIBSSH_THREADS)
        add_cmocka_test(torture_rand torture_rand.c ${TORTURE_LIBRARY})
        if (WITH_SERVER)
            add_cmocka_test(torture_server_x11 torture_server_x11.c ${TORTURE_LIBRARY})
        endif (WITH_SERVER)
    endif (LIBSSH_THREADS)
endif (UNIX AND NOT WIN32)
