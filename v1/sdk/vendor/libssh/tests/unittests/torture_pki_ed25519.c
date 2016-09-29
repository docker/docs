#define LIBSSH_STATIC

#include "torture.h"
#include "pki.c"
#include <sys/stat.h>
#include <fcntl.h>

#define LIBSSH_ED25519_TESTKEY "libssh_testkey.id_ed25519"

const unsigned char HASH[] = "12345678901234567890";
const uint8_t ref_signature[ED25519_SIG_LEN]=
    "\xbb\x8d\x55\x9f\x06\x14\x39\x24\xb4\xe1\x5a\x57\x3d\x9d\xbe\x22"
    "\x1b\xc1\x32\xd5\x55\x16\x00\x64\xce\xb4\xc3\xd2\xe3\x6f\x5e\x8d"
    "\x10\xa3\x18\x93\xdf\xa4\x96\x81\x11\x8e\x1e\x26\x14\x8a\x08\x1b"
    "\x01\x6a\x60\x59\x9c\x4a\x55\xa3\x16\x56\xf6\xc4\x50\x42\x7f\x03";

static void torture_pki_ed25519_sign(void **state){
    ssh_key privkey;
    ssh_signature sig = ssh_signature_new();
    ssh_string blob;
    int rc;
    (void)state;

    rc = ssh_pki_import_privkey_base64(torture_get_testkey(SSH_KEYTYPE_ED25519,0,0), NULL, NULL, NULL, &privkey);
    assert_true(rc == SSH_OK);

    sig->type = SSH_KEYTYPE_ED25519;
    rc = pki_ed25519_sign(privkey, sig, HASH, sizeof(HASH));
    assert_true(rc == SSH_OK);

    blob = pki_signature_to_blob(sig);
    assert_true(blob != NULL);

    assert_int_equal(ssh_string_len(blob), sizeof(ref_signature));
    assert_memory_equal(ssh_string_data(blob), ref_signature, sizeof(ref_signature));
    /* ssh_print_hexa("signature", ssh_string_data(blob), ssh_string_len(blob)); */
    ssh_signature_free(sig);
    ssh_key_free(privkey);
    ssh_string_free(blob);

}

static void torture_pki_ed25519_verify(void **state){
    ssh_key pubkey;
    ssh_signature sig;
    ssh_string blob = ssh_string_new(ED25519_SIG_LEN);
    char *pkey_ptr = strdup(strchr(torture_get_testkey_pub(SSH_KEYTYPE_ED25519,0), ' ') + 1);
    char *ptr;
    int rc;
    (void) state;

    /* remove trailing comment */
    ptr = strchr(pkey_ptr, ' ');
    if(ptr != NULL){
        *ptr = '\0';
    }
    rc = ssh_pki_import_pubkey_base64(pkey_ptr, SSH_KEYTYPE_ED25519, &pubkey);
    assert_true(rc == SSH_OK);

    ssh_string_fill(blob, ref_signature, ED25519_SIG_LEN);
    sig = pki_signature_from_blob(pubkey, blob, SSH_KEYTYPE_ED25519);
    assert_true(sig != NULL);

    rc = pki_ed25519_verify(pubkey, sig, HASH, sizeof(HASH));
    assert_true(rc == SSH_OK);

    ssh_signature_free(sig);
    /* alter signature and expect false result */

    ssh_key_free(pubkey);
    ssh_string_free(blob);
    free(pkey_ptr);
}

static void torture_pki_ed25519_verify_bad(void **state){
    ssh_key pubkey;
    ssh_signature sig;
    ssh_string blob = ssh_string_new(ED25519_SIG_LEN);
    char *pkey_ptr = strdup(strchr(torture_get_testkey_pub(SSH_KEYTYPE_ED25519,0), ' ') + 1);
    char *ptr;
    int rc;
    int i;
    (void) state;

    /* remove trailing comment */
    ptr = strchr(pkey_ptr, ' ');
    if(ptr != NULL){
        *ptr = '\0';
    }
    rc = ssh_pki_import_pubkey_base64(pkey_ptr, SSH_KEYTYPE_ED25519, &pubkey);
    assert_true(rc == SSH_OK);

    /* alter signature and expect false result */

    for (i=0; i < ED25519_SIG_LEN; ++i){
        ssh_string_fill(blob, ref_signature, ED25519_SIG_LEN);
        ((uint8_t *)ssh_string_data(blob))[i] ^= 0xff;
        sig = pki_signature_from_blob(pubkey, blob, SSH_KEYTYPE_ED25519);
        assert_true(sig != NULL);

        rc = pki_ed25519_verify(pubkey, sig, HASH, sizeof(HASH));
        assert_true(rc == SSH_ERROR);
        ssh_signature_free(sig);

    }
    ssh_key_free(pubkey);
    ssh_string_free(blob);
    free(pkey_ptr);
}

int torture_run_tests(void) {
    int rc;
    const UnitTest tests[] = {
        unit_test(torture_pki_ed25519_sign),
        unit_test(torture_pki_ed25519_verify),
        unit_test(torture_pki_ed25519_verify_bad)
    };

    ssh_init();
    rc=run_tests(tests);
    ssh_finalize();
    return rc;
}
