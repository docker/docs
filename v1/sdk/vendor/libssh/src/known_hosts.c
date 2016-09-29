/*
 * keyfiles.c - private and public key handling for authentication.
 *
 * This file is part of the SSH Library
 *
 * Copyright (c) 2003-2009 by Aris Adamantiadis
 * Copyright (c) 2009      by Andreas Schneider <asn@cryptomilk.org>
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

#include "config.h"

#include <ctype.h>
#include <errno.h>
#include <stdio.h>

#include "libssh/priv.h"
#include "libssh/session.h"
#include "libssh/buffer.h"
#include "libssh/misc.h"
#include "libssh/pki.h"
#include "libssh/options.h"
#include "libssh/knownhosts.h"
/*todo: remove this include */
#include "libssh/string.h"

#ifdef HAVE_LIBGCRYPT
#include <gcrypt.h>
#elif defined HAVE_LIBCRYPTO
#include <openssl/pem.h>
#include <openssl/dsa.h>
#include <openssl/err.h>
#include <openssl/rsa.h>
#endif /* HAVE_LIBCRYPTO */

#ifndef _WIN32
# include <netinet/in.h>
# include <arpa/inet.h>
#endif

/**
 * @addtogroup libssh_session
 *
 * @{
 */

static int alldigits(const char *s) {
  while (*s) {
    if (isdigit(*s)) {
      s++;
    } else {
      return 0;
    }
  }

  return 1;
}

/**
 * @internal
 *
 * @brief Free a token array.
 */
static void tokens_free(char **tokens) {
  if (tokens == NULL) {
    return;
  }

  SAFE_FREE(tokens[0]);
  /* It's not needed to free other pointers because tokens generated by
   * space_tokenize fit all in one malloc
   */
  SAFE_FREE(tokens);
}

/**
 * @internal
 *
 * @brief Return one line of known host file.
 *
 * This will return a token array containing (host|ip), keytype and key.
 *
 * @param[out] file     A pointer to the known host file. Could be pointing to
 *                      NULL at start.
 *
 * @param[in]  filename The file name of the known host file.
 *
 * @param[out] found_type A pointer to a string to be set with the found key
 *                        type.
 *
 * @returns             The found_type type of key (ie "dsa","ssh-rsa1"). Don't
 *                      free that value. NULL if no match was found or the file
 *                      was not found.
 */
static char **ssh_get_knownhost_line(FILE **file, const char *filename,
                                     const char **found_type) {
  char buffer[4096] = {0};
  char *ptr;
  char **tokens;

  if(*file == NULL){
    *file = fopen(filename,"r");
    if (*file == NULL) {
      return NULL;
    }
  }

  while (fgets(buffer, sizeof(buffer), *file)) {
    ptr = strchr(buffer, '\n');
    if (ptr) {
      *ptr =  '\0';
    }

    ptr = strchr(buffer,'\r');
    if (ptr) {
      *ptr = '\0';
    }

    if (buffer[0] == '\0' || buffer[0] == '#') {
      continue; /* skip empty lines */
    }

    tokens = space_tokenize(buffer);
    if (tokens == NULL) {
      fclose(*file);
      *file = NULL;

      return NULL;
    }

    if(!tokens[0] || !tokens[1] || !tokens[2]) {
      /* it should have at least 3 tokens */
      tokens_free(tokens);
      continue;
    }

    *found_type = tokens[1];
    if (tokens[3]) {
      /* openssh rsa1 format has 4 tokens on the line. Recognize it
         by the fact that everything is all digits */
      if (tokens[4]) {
        /* that's never valid */
        tokens_free(tokens);
        continue;
      }
      if (alldigits(tokens[1]) && alldigits(tokens[2]) && alldigits(tokens[3])) {
        *found_type = "ssh-rsa1";
      } else {
        /* 3 tokens only, not four */
        tokens_free(tokens);
        continue;
      }
    }

    return tokens;
  }

  fclose(*file);
  *file = NULL;

  /* we did not find anything, end of file*/
  return NULL;
}

/**
 * @brief Check the public key in the known host line matches the public key of
 * the currently connected server.
 *
 * @param[in] session   The SSH session to use.
 *
 * @param[in] tokens    A list of tokens in the known_hosts line.
 *
 * @returns             1 if the key matches, 0 if the key doesn't match and -1
 *                      on error.
 */
static int check_public_key(ssh_session session, char **tokens) {
  ssh_string pubkey = session->current_crypto->server_pubkey;
  ssh_buffer pubkey_buffer;
  char *pubkey_64;

  /* ok we found some public key in known hosts file. now un-base64it */
  if (alldigits(tokens[1])) {
    /* openssh rsa1 format */
    bignum tmpbn;
    ssh_string tmpstring;
    unsigned int len;
    int i;

    pubkey_buffer = ssh_buffer_new();
    if (pubkey_buffer == NULL) {
      return -1;
    }

    tmpstring = ssh_string_from_char("ssh-rsa1");
    if (tmpstring == NULL) {
      ssh_buffer_free(pubkey_buffer);
      return -1;
    }

    if (buffer_add_ssh_string(pubkey_buffer, tmpstring) < 0) {
      ssh_buffer_free(pubkey_buffer);
      ssh_string_free(tmpstring);
      return -1;
    }
    ssh_string_free(tmpstring);

    for (i = 2; i < 4; i++) { /* e, then n */
      tmpbn = NULL;
      bignum_dec2bn(tokens[i], &tmpbn);
      if (tmpbn == NULL) {
        ssh_buffer_free(pubkey_buffer);
        return -1;
      }
      /* for some reason, make_bignum_string does not work
         because of the padding which it does --kv */
      /* tmpstring = make_bignum_string(tmpbn); */
      /* do it manually instead */
      len = bignum_num_bytes(tmpbn);
      tmpstring = malloc(4 + len);
      if (tmpstring == NULL) {
        ssh_buffer_free(pubkey_buffer);
        bignum_free(tmpbn);
        return -1;
      }
      /* TODO: fix the hardcoding */
      tmpstring->size = htonl(len);
#ifdef HAVE_LIBGCRYPT
      bignum_bn2bin(tmpbn, len, ssh_string_data(tmpstring));
#elif defined HAVE_LIBCRYPTO
      bignum_bn2bin(tmpbn, ssh_string_data(tmpstring));
#endif
      bignum_free(tmpbn);
      if (buffer_add_ssh_string(pubkey_buffer, tmpstring) < 0) {
        ssh_buffer_free(pubkey_buffer);
        ssh_string_free(tmpstring);
        bignum_free(tmpbn);
        return -1;
      }
      ssh_string_free(tmpstring);
    }
  } else {
    /* ssh-dss or ssh-rsa */
    pubkey_64 = tokens[2];
    pubkey_buffer = base64_to_bin(pubkey_64);
  }

  if (pubkey_buffer == NULL) {
    ssh_set_error(session, SSH_FATAL,
        "Verifying that server is a known host: base64 error");
    return -1;
  }

  if (buffer_get_rest_len(pubkey_buffer) != ssh_string_len(pubkey)) {
    ssh_buffer_free(pubkey_buffer);
    return 0;
  }

  /* now test that they are identical */
  if (memcmp(buffer_get_rest(pubkey_buffer), ssh_string_data(pubkey),
        buffer_get_rest_len(pubkey_buffer)) != 0) {
    ssh_buffer_free(pubkey_buffer);
    return 0;
  }

  ssh_buffer_free(pubkey_buffer);
  return 1;
}

/**
 * @brief Check if a hostname matches a openssh-style hashed known host.
 *
 * @param[in]  host     The host to check.
 *
 * @param[in]  hashed   The hashed value.
 *
 * @returns             1 if it matches, 0 otherwise.
 */
static int match_hashed_host(const char *host, const char *sourcehash)
{
  /* Openssh hash structure :
   * |1|base64 encoded salt|base64 encoded hash
   * hash is produced that way :
   * hash := HMAC_SHA1(key=salt,data=host)
   */
  unsigned char buffer[256] = {0};
  ssh_buffer salt;
  ssh_buffer hash;
  HMACCTX mac;
  char *source;
  char *b64hash;
  int match;
  unsigned int size;

  if (strncmp(sourcehash, "|1|", 3) != 0) {
    return 0;
  }

  source = strdup(sourcehash + 3);
  if (source == NULL) {
    return 0;
  }

  b64hash = strchr(source, '|');
  if (b64hash == NULL) {
    /* Invalid hash */
    SAFE_FREE(source);

    return 0;
  }

  *b64hash = '\0';
  b64hash++;

  salt = base64_to_bin(source);
  if (salt == NULL) {
    SAFE_FREE(source);

    return 0;
  }

  hash = base64_to_bin(b64hash);
  SAFE_FREE(source);
  if (hash == NULL) {
    ssh_buffer_free(salt);

    return 0;
  }

  mac = hmac_init(buffer_get_rest(salt), buffer_get_rest_len(salt), SSH_HMAC_SHA1);
  if (mac == NULL) {
    ssh_buffer_free(salt);
    ssh_buffer_free(hash);

    return 0;
  }
  size = sizeof(buffer);
  hmac_update(mac, host, strlen(host));
  hmac_final(mac, buffer, &size);

  if (size == buffer_get_rest_len(hash) &&
      memcmp(buffer, buffer_get_rest(hash), size) == 0) {
    match = 1;
  } else {
    match = 0;
  }

  ssh_buffer_free(salt);
  ssh_buffer_free(hash);

  SSH_LOG(SSH_LOG_PACKET,
      "Matching a hashed host: %s match=%d", host, match);

  return match;
}

/* How it's working :
 * 1- we open the known host file and bitch if it doesn't exist
 * 2- we need to examine each line of the file, until going on state SSH_SERVER_KNOWN_OK:
 *  - there's a match. if the key is good, state is SSH_SERVER_KNOWN_OK,
 *    else it's SSH_SERVER_KNOWN_CHANGED (or SSH_SERVER_FOUND_OTHER)
 *  - there's no match : no change
 */

/**
 * @brief Check if the server is known.
 *
 * Checks the user's known host file for a previous connection to the
 * current server.
 *
 * @param[in]  session  The SSH session to use.
 *
 * @returns SSH_SERVER_KNOWN_OK:       The server is known and has not changed.\n
 *          SSH_SERVER_KNOWN_CHANGED:  The server key has changed. Either you
 *                                     are under attack or the administrator
 *                                     changed the key. You HAVE to warn the
 *                                     user about a possible attack.\n
 *          SSH_SERVER_FOUND_OTHER:    The server gave use a key of a type while
 *                                     we had an other type recorded. It is a
 *                                     possible attack.\n
 *          SSH_SERVER_NOT_KNOWN:      The server is unknown. User should
 *                                     confirm the MD5 is correct.\n
 *          SSH_SERVER_FILE_NOT_FOUND: The known host file does not exist. The
 *                                     host is thus unknown. File will be
 *                                     created if host key is accepted.\n
 *          SSH_SERVER_ERROR:          Some error happened.
 *
 * @see ssh_get_pubkey_hash()
 *
 * @bug There is no current way to remove or modify an entry into the known
 *      host table.
 */
int ssh_is_server_known(ssh_session session) {
  FILE *file = NULL;
  char **tokens;
  char *host;
  char *hostport;
  const char *type;
  int match;
  int ret = SSH_SERVER_NOT_KNOWN;

  if (session->opts.knownhosts == NULL) {
    if (ssh_options_apply(session) < 0) {
      ssh_set_error(session, SSH_REQUEST_DENIED,
          "Can't find a known_hosts file");

      return SSH_SERVER_FILE_NOT_FOUND;
    }
  }

  if (session->opts.host == NULL) {
    ssh_set_error(session, SSH_FATAL,
        "Can't verify host in known hosts if the hostname isn't known");

    return SSH_SERVER_ERROR;
  }

  if (session->current_crypto == NULL){
  	ssh_set_error(session, SSH_FATAL,
  			"ssh_is_host_known called without cryptographic context");

  	return SSH_SERVER_ERROR;
  }
  host = ssh_lowercase(session->opts.host);
  hostport = ssh_hostport(host, session->opts.port > 0 ? session->opts.port : 22);
  if (host == NULL || hostport == NULL) {
    ssh_set_error_oom(session);
    SAFE_FREE(host);
    SAFE_FREE(hostport);

    return SSH_SERVER_ERROR;
  }

  do {
    tokens = ssh_get_knownhost_line(&file,
                                    session->opts.knownhosts,
                                    &type);

    /* End of file, return the current state */
    if (tokens == NULL) {
      break;
    }
    match = match_hashed_host(host, tokens[0]);
    if (match == 0){
    	match = match_hostname(hostport, tokens[0], strlen(tokens[0]));
    }
    if (match == 0) {
      match = match_hostname(host, tokens[0], strlen(tokens[0]));
    }
    if (match == 0) {
      match = match_hashed_host(hostport, tokens[0]);
    }
    if (match) {
      /* We got a match. Now check the key type */
      if (strcmp(session->current_crypto->server_pubkey_type, type) != 0) {
          SSH_LOG(SSH_LOG_PACKET,
                  "ssh_is_server_known: server type [%s] doesn't match the "
                  "type [%s] in known_hosts file",
                  session->current_crypto->server_pubkey_type,
                  type);
        /* Different type. We don't override the known_changed error which is
         * more important */
        if (ret != SSH_SERVER_KNOWN_CHANGED)
          ret = SSH_SERVER_FOUND_OTHER;
        tokens_free(tokens);
        continue;
      }
      /* so we know the key type is good. We may get a good key or a bad key. */
      match = check_public_key(session, tokens);
      tokens_free(tokens);

      if (match < 0) {
        ret = SSH_SERVER_ERROR;
        break;
      } else if (match == 1) {
        ret = SSH_SERVER_KNOWN_OK;
        break;
      } else if(match == 0) {
        /* We override the status with the wrong key state */
        ret = SSH_SERVER_KNOWN_CHANGED;
      }
    } else {
      tokens_free(tokens);
    }
  } while (1);

  if ((ret == SSH_SERVER_NOT_KNOWN) &&
      (session->opts.StrictHostKeyChecking == 0)) {
    ssh_write_knownhost(session);
    ret = SSH_SERVER_KNOWN_OK;
  }

  SAFE_FREE(host);
  SAFE_FREE(hostport);
  if (file != NULL) {
    fclose(file);
  }

  /* Return the current state at end of file */
  return ret;
}

/**
 * @brief Write the current server as known in the known hosts file.
 *
 * This will create the known hosts file if it does not exist. You generaly use
 * it when ssh_is_server_known() answered SSH_SERVER_NOT_KNOWN.
 *
 * @param[in]  session  The ssh session to use.
 *
 * @return              SSH_OK on success, SSH_ERROR on error.
 */
int ssh_write_knownhost(ssh_session session) {
    ssh_key key;
    ssh_string pubkey_s;
    char *b64_key;
    char buffer[4096] = {0};
    FILE *file;
    char *dir;
    char *host;
    char *hostport;
    int rc;

    if (session->opts.host == NULL) {
        ssh_set_error(session, SSH_FATAL,
                "Can't write host in known hosts if the hostname isn't known");
        return SSH_ERROR;
    }

    host = ssh_lowercase(session->opts.host);
    /* If using a nonstandard port, save the host in the [host]:port format */
    if (session->opts.port > 0 && session->opts.port != 22) {
        hostport = ssh_hostport(host, session->opts.port);
        SAFE_FREE(host);
        if (hostport == NULL) {
            return SSH_ERROR;
        }
        host = hostport;
        hostport = NULL;
    }

    if (session->opts.knownhosts == NULL) {
        if (ssh_options_apply(session) < 0) {
            ssh_set_error(session, SSH_FATAL, "Can't find a known_hosts file");
            SAFE_FREE(host);
            return SSH_ERROR;
        }
    }

    if (session->current_crypto==NULL) {
        ssh_set_error(session, SSH_FATAL, "No current crypto context");
        SAFE_FREE(host);
        return SSH_ERROR;
    }

    pubkey_s = session->current_crypto->server_pubkey;
    if (pubkey_s == NULL){
        ssh_set_error(session, SSH_FATAL, "No public key present");
        SAFE_FREE(host);
        return SSH_ERROR;
    }

    /* Check if ~/.ssh exists and create it if not */
    dir = ssh_dirname(session->opts.knownhosts);
    if (dir == NULL) {
        ssh_set_error(session, SSH_FATAL, "%s", strerror(errno));
        SAFE_FREE(host);
        return SSH_ERROR;
    }

    if (!ssh_file_readaccess_ok(dir)) {
        if (ssh_mkdir(dir, 0700) < 0) {
            ssh_set_error(session, SSH_FATAL,
                    "Cannot create %s directory.", dir);
            SAFE_FREE(dir);
            SAFE_FREE(host);
            return SSH_ERROR;
        }
    }
    SAFE_FREE(dir);

    file = fopen(session->opts.knownhosts, "a");
    if (file == NULL) {
        ssh_set_error(session, SSH_FATAL,
                "Couldn't open known_hosts file %s for appending: %s",
                session->opts.knownhosts, strerror(errno));
        SAFE_FREE(host);
        return SSH_ERROR;
    }

    rc = ssh_pki_import_pubkey_blob(pubkey_s, &key);
    if (rc < 0) {
        fclose(file);
        SAFE_FREE(host);
        return -1;
    }

    if (strcmp(session->current_crypto->server_pubkey_type, "ssh-rsa1") == 0) {
        /* openssh uses a different format for ssh-rsa1 keys.
           Be compatible --kv */
        rc = ssh_pki_export_pubkey_rsa1(key, host, buffer, sizeof(buffer));
        ssh_key_free(key);
        SAFE_FREE(host);
        if (rc < 0) {
            fclose(file);
            return -1;
        }
    } else {
        rc = ssh_pki_export_pubkey_base64(key, &b64_key);
        if (rc < 0) {
            ssh_key_free(key);
            fclose(file);
            SAFE_FREE(host);
            return -1;
        }

        snprintf(buffer, sizeof(buffer),
                "%s %s %s\n",
                host,
                key->type_c,
                b64_key);

        ssh_key_free(key);
        SAFE_FREE(host);
        SAFE_FREE(b64_key);
    }

    if (fwrite(buffer, strlen(buffer), 1, file) != 1 || ferror(file)) {
        fclose(file);
        return -1;
    }

    fclose(file);
    return 0;
}

#define KNOWNHOSTS_MAXTYPES 10

/**
 * @internal
 * @brief Check which kind of host keys should be preferred for connection
 *        by reading the known_hosts file.
 *
 * @param[in]  session  The SSH session to use.
 *
 * @returns array of supported key types
 *			NULL on error
 */
char **ssh_knownhosts_algorithms(ssh_session session) {
  FILE *file = NULL;
  char **tokens;
  char *host;
  char *hostport;
  const char *type;
  int match;
  char **array;
  int i=0, j;

  if (session->opts.knownhosts == NULL) {
    if (ssh_options_apply(session) < 0) {
      ssh_set_error(session, SSH_REQUEST_DENIED,
          "Can't find a known_hosts file");
      return NULL;
    }
  }

  if (session->opts.host == NULL) {
    return NULL;
  }

  host = ssh_lowercase(session->opts.host);
  hostport = ssh_hostport(host, session->opts.port > 0 ? session->opts.port : 22);
  array = malloc(sizeof(char *) * KNOWNHOSTS_MAXTYPES);

  if (host == NULL || hostport == NULL || array == NULL) {
    ssh_set_error_oom(session);
    SAFE_FREE(host);
    SAFE_FREE(hostport);
    SAFE_FREE(array);
    return NULL;
  }

  do {
    tokens = ssh_get_knownhost_line(&file,
    		session->opts.knownhosts, &type);

    /* End of file, return the current state */
    if (tokens == NULL) {
      break;
    }
    match = match_hashed_host(host, tokens[0]);
    if (match == 0){
    	match = match_hostname(hostport, tokens[0], strlen(tokens[0]));
    }
    if (match == 0) {
      match = match_hostname(host, tokens[0], strlen(tokens[0]));
    }
    if (match == 0) {
      match = match_hashed_host(hostport, tokens[0]);
    }
    if (match) {
      /* We got a match. Now check the key type */
    	SSH_LOG(SSH_LOG_DEBUG, "server %s:%d has %s in known_hosts",
    							host, session->opts.port, type);
    	/* don't copy more than once */
    	for(j=0;j<i && match;++j){
    		if(strcmp(array[j], type)==0)
    			match=0;
    	}
    	if (match){
    		array[i] = strdup(type);
    		i++;
    		if(i>= KNOWNHOSTS_MAXTYPES-1){
    			tokens_free(tokens);
    			break;
    		}
    	}
    }
    tokens_free(tokens);
  } while (1);

  array[i]=NULL;
  SAFE_FREE(host);
  SAFE_FREE(hostport);
  if (file != NULL) {
    fclose(file);
  }

  /* Return the current state at end of file */
  return array;
}

/** @} */

/* vim: set ts=4 sw=4 et cindent: */
