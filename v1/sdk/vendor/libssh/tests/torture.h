/*
 * torture.c - torture library for testing libssh
 *
 * This file is part of the SSH Library
 *
 * Copyright (c) 2008-2009 by Andreas Schneider <asn@cryptomilk.org>
 *
 * The SSH Library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation; either version 2.1 of the License, or (at your
 * option) any later version.
 *
 * The SSH Library is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public
 * License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with the SSH Library; see the file COPYING.  If not, write to
 * the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston,
 * MA 02111-1307, USA.
 */

#ifndef _TORTURE_H
#define _TORTURE_H

#ifndef _GNU_SOURCE
#define _GNU_SOURCE
#endif

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>

#include "libssh/priv.h"
#include "libssh/server.h"
#include "libssh/sftp.h"

#include <cmocka.h>

#ifndef assert_return_code
/* hack for older versions of cmocka */
#define assert_return_code(code, errno) \
    assert_true(code >= 0)
#endif /* assert_return_code */

#define TORTURE_TESTKEY_PASSWORD "libssh-rocks"

/* Used by main to communicate with parse_opt. */
struct argument_s {
  const char *pattern;
  int verbose;
};

struct torture_sftp {
    ssh_session ssh;
    sftp_session sftp;
    char *testdir;
};

void torture_cmdline_parse(int argc, char **argv, struct argument_s *arguments);

int torture_rmdirs(const char *path);
int torture_isdir(const char *path);

/*
 * Returns the verbosity level asked by user
 */
int torture_libssh_verbosity(void);

ssh_session torture_ssh_session(const char *host,
                                const unsigned int *port,
                                const char *user,
                                const char *password);

ssh_bind torture_ssh_bind(const char *addr,
                          const unsigned int port,
                          enum ssh_keytypes_e key_type,
                          const char *private_key_file);

struct torture_sftp *torture_sftp_session(ssh_session session);
void torture_sftp_close(struct torture_sftp *t);

const char *torture_get_testkey(enum ssh_keytypes_e type,
                                int ecdsa_bits,
                                int with_passphrase);
const char *torture_get_testkey_pub(enum ssh_keytypes_e type, int ecdsa_bits);
const char *torture_get_testkey_passphrase(void);

void torture_write_file(const char *filename, const char *data);

#define torture_filter_tests(tests) _torture_filter_tests(tests, sizeof(tests) / sizeof(tests)[0])
void _torture_filter_tests(UnitTest *tests, size_t ntests);

/*
 * This function must be defined in every unit test file.
 */
int torture_run_tests(void);

#endif /* _TORTURE_H */
