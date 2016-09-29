#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <err.h>
#include <ctype.h>
#include <sys/types.h>
#include <sys/uio.h>
#include <sys/un.h>
#include <sys/socket.h>
#include <pthread.h>
#include <sys/stat.h>
#include <asl.h>
#include <pwd.h>
#include <assert.h>
#include <netinet/in.h>
#include <inttypes.h>
#include <dispatch/dispatch.h>
#include <vmnet/vmnet.h>
#include <ftw.h>

#include "proxy.h"
#include "protocol.h"
#include "utils.h"
#include "apple_utils.h"


#define DOCKER_VMNETD_BINARY "/Library/PrivilegedHelperTools/com.docker.vmnetd"
#define DOCKER_VMNETD_PLIST "/Library/LaunchDaemons/com.docker.vmnetd.plist"
#define DOCKER_VMNETD_RX_BATCH_LEN 64

#define GROUP_CONTAINERS "Library/Group Containers/group.com.docker"

const static char * const docker_binaries[] =
  { "docker", "docker-compose", "notary",
    "docker-machine" };

const static char * const deleted_docker_binaries[] =
  { "docker-configure", "docker-diagnose", "pinata" };


/* These sockets are not installed by vmnetd, but this list can be used
 when receiving uninstall_sockets command. Taking advantage of vmnet root
 access rights to uninstall them without issues. */
const static char * const sockets_to_uninstall[] =
  { "/var/tmp/com.docker.db.socket", "/var/tmp/com.docker.lofs.socket",
  "/var/tmp/com.docker.osxfs.volume.socket", "/var/tmp/com.docker.slirp.port.socket", "/var/tmp/docker.sock",
  "/var/tmp/com.docker.slirp.socket", "/tmp/fs.socket", "/tmp/u.socket", "/tmp/u2.socket",
  "/tmp/com.docker.9p.filesystem.socket", "/tmp/com.docker.proxy.underlying.socket", "/tmp/diagnose.socket",
  "/var/tmp/com.docker.vmnet.port.socket", "/var/tmp/com.docker.port.socket", "/var/tmp/com.docker.vsock" };

/* This is the proxy code shared by both the launchd frontend (main_launchd.c)
   and the regular Unix command-line version (main.c).
 */

/* From @avsm's ocaml-vmnet package
 * Copyright (C) 2014 Anil Madhavapeddy <anil@recoil.org>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

static int unlink_callback(const char *fpath, const struct stat *sb, int typeflag, struct FTW *ftwbuf)
{
    // NOTE(aduermael): Ok, that is silly... But I can't compile if these
    // arguments are not used (and I don't want to turn off that warning)
    if (sb == NULL && typeflag == 0 && ftwbuf == NULL && false) {}

      int rv = remove(fpath);
      if (rv)
        aslLog(ASL_LEVEL_ERR, "couldn't delete %s", fpath);
      else
        aslLog(ASL_LEVEL_NOTICE, "deleted %s", fpath);

      return rv;
}

static int rmrf(const char *path)
{
    struct stat buf;
    if (stat(path, &buf) != 0) {
      if (errno != ENOENT) {
        aslLog(ASL_LEVEL_CRIT, "Failed to stat %s : %s", path, strerror (errno));
      }
    } else {
      // if is regular file
      if (S_ISDIR(buf.st_mode)) {
        return nftw(path, unlink_callback, 64, FTW_DEPTH | FTW_PHYS);
      } else {
        if (unlink(path) != 0) aslLog(ASL_LEVEL_ERR, "couldn't delete %s", path);
        else aslLog(ASL_LEVEL_NOTICE, "deleted %s", path);
        return 0;
      }
    }
    return 0;
}

 static void hexdump(unsigned char *buffer, size_t len){
   char ascii[17];
   size_t i = 0;
   ascii[16] = '\000';
   while (i < len) {
     unsigned char c = *(buffer + i);
     printf("%02x ", c);
     ascii[i++ % 16] = isprint(c)?(signed char)c:'.';
     if ((i % 2) == 0) printf(" ");
     if ((i % 16) == 0) printf(" %s\r\n", ascii);
   };
   printf("\r\n");
 }

 struct vmnet_state {
   interface_ref iref;
   struct vif_info vif;
   int client_fd;
   int event_counter;
   dispatch_queue_t task_sq;
};

static char *string_of_vmnet_return_t(vmnet_return_t status){
  switch (status) {
  case VMNET_SUCCESS:          return "VMNET_SUCCESS";
  case VMNET_FAILURE:          return "VMNET_FAILURE";
  case VMNET_MEM_FAILURE:      return "VMNET_MEM_FAILURE";
  case VMNET_INVALID_ARGUMENT: return "VMNET_INVALID_ARGUMENT";
  case VMNET_SETUP_INCOMPLETE: return "VMNET_SETUP_INCOMPLETE";
  case VMNET_INVALID_ACCESS:   return "VMNET_INVALID_ACCESS";
  case VMNET_PACKET_TOO_BIG:   return "VMNET_PACKET_TOO_BIG";
  case VMNET_BUFFER_EXHAUSTED: return "VMNET_BUFFER_EXHAUSTED";
  case VMNET_TOO_MANY_PACKETS: return "VMNET_TOO_MANY_PACKETS";
  default: return NULL;
  }
}

static struct vmnet_state*
alloc_vmnet_state(interface_ref i, dispatch_queue_t task_queue)
{
  struct vmnet_state *vms = malloc(sizeof(struct vmnet_state));
  bzero(vms, sizeof(struct vmnet_state));
  if (!vms) abort();
  vms->client_fd = -1;
  vms->iref = i;
  vms->task_sq = task_queue;
  return vms;
}

static struct vmnet_state *open_vmnet(const char *uuid_string){
  __block struct vmnet_state *vms = NULL;
  xpc_object_t interface_desc = xpc_dictionary_create(NULL, NULL, 0);
  xpc_dictionary_set_uint64(interface_desc, vmnet_operation_mode_key, VMNET_SHARED_MODE);
  uuid_t uuid;
  if (uuid_string) {
    uuid_parse(uuid_string, uuid);
  } else {
    uuid_generate_random(uuid);
  }
  xpc_dictionary_set_uuid(interface_desc, vmnet_interface_id_key, uuid);
  __block interface_ref iface = NULL;
  __block vmnet_return_t iface_status = 0;
  char *iface_status_string = NULL;
  dispatch_semaphore_t iface_created = dispatch_semaphore_create(0);
  dispatch_queue_t task_q = dispatch_queue_create("com.docker.proxy.task_sq", DISPATCH_QUEUE_SERIAL);

  iface = vmnet_start_interface(interface_desc, task_q,
    ^(vmnet_return_t status, xpc_object_t interface_param) {
      iface_status = status;
      if (status != VMNET_SUCCESS || !interface_param) {
         if (interface_param != NULL) {
            aslLog(ASL_LEVEL_ERR, "com.docker.vmnetd: interface_param = %s",
                    xpc_copy_description(interface_param));
         } else {
            aslLog(ASL_LEVEL_ERR, "com.docker.vmnetd: interface_param = NULL");
         }
         dispatch_semaphore_signal(iface_created);
         return;
      }
      //printf("mac desc: %s\n", xpc_copy_description(xpc_dictionary_get_value(interface_param, vmnet_mac_address_key)));
      vms = alloc_vmnet_state(iface, task_q);
      const char *mac = xpc_dictionary_get_string(interface_param, vmnet_mac_address_key);
      aslLog(ASL_LEVEL_NOTICE, "com.docker.vmnetd reports MAC address: %s", mac);

      if (sscanf(mac, "%hhx:%hhx:%hhx:%hhx:%hhx:%hhx", &vms->vif.mac[0], &vms->vif.mac[1], &vms->vif.mac[2], &vms->vif.mac[3], &vms->vif.mac[4], &vms->vif.mac[5]) != 6) {
        aslLog(ASL_LEVEL_ERR, "Failed to parse MAC address from vmnet: %s", mac);
        exit(1);
      }
      vms->vif.mtu = xpc_dictionary_get_uint64(interface_param, vmnet_mtu_key);
      vms->vif.max_packet_size = xpc_dictionary_get_uint64(interface_param, vmnet_max_packet_size_key);
      dispatch_semaphore_signal(iface_created);
    });
  dispatch_semaphore_wait(iface_created, DISPATCH_TIME_FOREVER);

  iface_status_string = string_of_vmnet_return_t(iface_status);
  if (iface == NULL || iface_status != VMNET_SUCCESS) {
    if (iface_status_string != NULL) {
      aslLog(ASL_LEVEL_ERR, "Failed to initialise com.docker.vmnetd: status = %s", iface_status_string);
    } else {
      aslLog(ASL_LEVEL_ERR, "Failed to initialise com.docker.vmnetd: status = NULL");
    }
    exit(1);
  }

  // atexit handler to stop vmnet framework interface
  atexit_b(^{
          aslLog(ASL_LEVEL_NOTICE, "Stopping vmnet interface");
          dispatch_semaphore_t iface_stopped = dispatch_semaphore_create(0);
          vmnet_return_t res = vmnet_stop_interface(iface, task_q, ^(vmnet_return_t status) {
              if (status == VMNET_SUCCESS) {
                aslLog(ASL_LEVEL_NOTICE, "Vmnet interface stopped");
              } else {
                aslLog(ASL_LEVEL_ERR, "Failed to stop Vmnet interface");
              }
              dispatch_semaphore_signal(iface_stopped);
          });
          if (res != VMNET_SUCCESS) {
              aslLog(ASL_LEVEL_ERR, "Call to vmnet_stop_interface failed");
          } else {
              dispatch_semaphore_wait(iface_stopped, DISPATCH_TIME_FOREVER);
              aslLog(ASL_LEVEL_NOTICE, "Vmnet shutdown done");
          }
      });

  /* Prepare buffers for vmnet_read */
  __block struct { // Wrap in struct, as arrays can not be used with __block
		struct iovec iov_r[DOCKER_VMNETD_RX_BATCH_LEN];
		struct iovec iov_w[DOCKER_VMNETD_RX_BATCH_LEN][2]; // For calling writev()
		uint8_t headers[DOCKER_VMNETD_RX_BATCH_LEN][2];
		struct vmpktdesc v[DOCKER_VMNETD_RX_BATCH_LEN];
	} pktbufs;

  for (int i = 0; i < DOCKER_VMNETD_RX_BATCH_LEN; i++) {
    void *buf = malloc(vms->vif.max_packet_size);
    if (!buf) {
			aslLog(ASL_LEVEL_ERR,
					"Failed to allocate %lu bytes as a read buffer.", vms->vif.max_packet_size);
			exit(1);
    }
    pktbufs.iov_r[i].iov_base = buf;
    pktbufs.iov_r[i].iov_len = vms->vif.max_packet_size;
    pktbufs.v[i].vm_pkt_size = vms->vif.max_packet_size;
    pktbufs.v[i].vm_pkt_iov = &(pktbufs.iov_r[i]);
    pktbufs.v[i].vm_pkt_iovcnt = 1;
    pktbufs.v[i].vm_flags = 0;
  }

	aslLog(ASL_LEVEL_NOTICE,
									"rx batching enabled (len=%d)", DOCKER_VMNETD_RX_BATCH_LEN);

	/* Register event callback */
	vmnet_interface_set_event_callback(iface, VMNET_INTERFACE_PACKETS_AVAILABLE, vms->task_sq,
			^(__unused interface_event_t event_id, xpc_object_t event) {
				uint64_t estimated = xpc_dictionary_get_uint64(event, vmnet_estimated_packets_available_key);
				int pktcnt = (int)estimated;
				if (estimated > DOCKER_VMNETD_RX_BATCH_LEN) {
					pktcnt = DOCKER_VMNETD_RX_BATCH_LEN;
				}

				vmnet_return_t res = vmnet_read(iface, pktbufs.v, &pktcnt);

				if (res != VMNET_SUCCESS) {
					aslLog(ASL_LEVEL_ERR,
						"vmnet_read(maxlen=%lu) failed; res = %d, batch len=%d", vms->vif.max_packet_size, res, DOCKER_VMNETD_RX_BATCH_LEN);
					exit(1);
				}

				for (int i = 0; i < pktcnt; i++) {
					size_t received = pktbufs.v[i].vm_pkt_size;
					pktbufs.v[i].vm_pkt_size = vms->vif.max_packet_size; // reset pktsize for next iteration
					void *buf = pktbufs.v[i].vm_pkt_iov->iov_base;
					if (received > 0) {
						pktbufs.headers[i][0] = (uint8_t) ((received >> 0) & 0xff);
						pktbufs.headers[i][1] = (uint8_t) ((received >> 8) & 0xff);
						if (verbose_flag){
							aslLog(ASL_LEVEL_DEBUG, "Writing packet of length %lu to client", received);
							hexdump(buf, received);
						}

						// Update writev buffers
						pktbufs.iov_w[i][0].iov_base = &(pktbufs.headers[i]); // store header bytes in first vector
						pktbufs.iov_w[i][0].iov_len = 2;
						pktbufs.iov_w[i][1].iov_base = buf;
						pktbufs.iov_w[i][1].iov_len = received;

					} else {
						aslLog(ASL_LEVEL_NOTICE,
									"vmnet_read returned empty packet");
					}
				}
				// Write all packets in one go
				if (really_writev(vms->client_fd, (void *)pktbufs.iov_w, pktcnt * 2) == -1) {
						aslLog(ASL_LEVEL_ERR,
									"batch call to vmnet_write failed");
                        exit(1);
                }
			});

  if (verbose_flag) {
    aslLog(ASL_LEVEL_DEBUG,
      "Opened vmnet with mac = %02x:%02x:%02x:%02x:%02x:%02x; mtu = %lu; maximum_packet_size = %lu",
      vms->vif.mac[0], vms->vif.mac[1], vms->vif.mac[2], vms->vif.mac[3], vms->vif.mac[4], vms->vif.mac[5],
      vms->vif.mtu, vms->vif.max_packet_size);
  };

  return vms;
}

static void vmnet_write_packet(struct vmnet_state *vms, void *buffer, size_t len){
  if (len < 14) {
		/* Ignore packets that are too short to contain an Ethernet frame to avoid crash in Vmnet.framework (issue #456) */
		aslLog(ASL_LEVEL_ERR,
			"vmnet_write_packet ignored packet with len < 14");
		return;
  }
  interface_ref iface = vms->iref;
  struct iovec iov;
  iov.iov_base = buffer;
  iov.iov_len = len;
  struct vmpktdesc v;
  v.vm_pkt_size = len;
  v.vm_pkt_iov = &iov;
  v.vm_pkt_iovcnt = 1;
  v.vm_flags = 0; /* TODO no clue what this is */
  int pktcnt = 1;
  vmnet_return_t res = vmnet_write(iface, &v, &pktcnt);
  if (res == VMNET_SUCCESS) {
    return;
  }
  aslLog(ASL_LEVEL_ERR,
    "vmnet_write_packet failed res = %d", res);
  exit(1);
}

static void *proxy_writes_forever(void *ptr) {
  struct vmnet_state *vms = (struct vmnet_state *)ptr;
  void *buffer = malloc(vms->vif.max_packet_size);
  uint8_t header[2];
  size_t length;
  if (!buffer) {
    aslLog( ASL_LEVEL_ERR,
      "Failed to allocate %lu bytes as a write buffer.", vms->vif.max_packet_size);
    exit(1);
  }
  while (1) {
    if (really_read(vms->client_fd, header, 2) == -1) {
      aslLog(ASL_LEVEL_ERR, "Failed to read header");
      exit(1);
    }
    length = (size_t)(header[0] | (header[1] << 8));
    if (length > vms->vif.max_packet_size){
      aslLog(ASL_LEVEL_ERR,
        "Received over-large packet from Unix domain socket; len = %lu", length);
      exit(1);
    }
    if (really_read(vms->client_fd, buffer, length) == -1) {
      aslLog(ASL_LEVEL_ERR, "Failed to read packet data");
      exit(1);
    }
    if (verbose_flag){
      aslLog(ASL_LEVEL_DEBUG,
        "Writing packet of length %lu to com.docker.vmnetd", length);
      hexdump(buffer, length);
    }
    vmnet_write_packet(vms, buffer, length);
  }
  return NULL;
}

static void proxy_forever(struct vmnet_state *vms){
  pthread_t writer;
  int err;
  if((err = pthread_create(&writer, NULL, proxy_writes_forever, (void*)vms)) != 0) {
    aslLog(ASL_LEVEL_ERR,
      "pthread_created failed with error %d", err);
    exit(1);
  }
  if((err = pthread_join(writer, NULL)) != 0) {
    aslLog(ASL_LEVEL_ERR,
      "pthread_join failed with error %d", err);
    exit(1);
  }
}

/* install_symlinks does a bunch of things:
    - delete old symlinks listed in deleted_docker_binaries
    - create symlinks in /usr/local/bin for binaries in docker_binaries,
      suffixing pre-existing binaries with ".backup"
    - create symlink from /var/run/docker.sock to user's docker api endpoint socket (s60),
      suffixing pre-existing non-symlink/non-socket files with ".backup"
*/

static void really_install_symlinks(int, const char*) __attribute__((noreturn));

static void really_install_symlinks(int client_fd, const char *container_path){
  uint8_t r = 1;
  uid_t uid;
  gid_t gid;
  char *buffer = NULL;

  if (getpeereid(client_fd, &uid, &gid) != 0) {
    aslLog(ASL_LEVEL_CRIT, "Failed to discover the uid/gid of the peer");
    goto out2;
  }
  aslLog(ASL_LEVEL_CRIT, "client with uid %d requests com.docker.vmnetd install symlinks", uid);
  /* the example from getpwent(3) */
  long bufsize;
  if ((bufsize = sysconf(_SC_GETPW_R_SIZE_MAX)) == -1)
    abort();
  if ((buffer = malloc((size_t)bufsize)) == NULL){
    aslLog(ASL_LEVEL_CRIT, "Failed to allocate buffer of size %ld", bufsize);
    goto out2;
  }
  struct passwd pwd, *pwd_result = NULL;
  if (getpwuid_r(uid, &pwd, buffer, (size_t)bufsize, &pwd_result) != 0 || !pwd_result) {
    aslLog(ASL_LEVEL_CRIT, "Failed to discover the homedir of the peer");
    goto out2;
  }
  if (mkdir("/usr/local", 0755) == -1) {
    /* If it already exists this is ok */
    if (errno != EEXIST) {
      aslLog(ASL_LEVEL_CRIT, "Failed to mkdir /usr/local");
      goto out2;
    }
  } else {
    if (chown("/usr/local", uid, gid) == -1) {
      aslLog(ASL_LEVEL_CRIT, "Failed to chown created /usr/local to %d.%d: %s", uid, gid, strerror (errno));
      goto out2;
    }
  }
  if (mkdir("/usr/local/bin", 0755) == -1) {
    /* If it already exists this is ok */
    if (errno != EEXIST) {
      aslLog(ASL_LEVEL_CRIT, "Failed to mkdir /usr/local/bin");
      goto out2;
    }
  } else {
    if (chown("/usr/local/bin", uid, gid) == -1) {
      aslLog(ASL_LEVEL_CRIT, "Failed to chown created /usr/local/bin to %d.%d: %s", uid, gid, strerror (errno));
      goto out2;
    }
  }
  char group_path[PATH_MAX];
  char usr_local_bin_path[PATH_MAX];
  char usr_local_bin_backup_path[PATH_MAX];
  struct stat buf;

  /* Try to find old binary symlinks that should have been deleted and delete them.
  These binaries are listed in deleted_docker_binaries */
  for (unsigned long i = 0; i < sizeof(deleted_docker_binaries) / sizeof(deleted_docker_binaries[0]); i++ ) {
    if (snprintf(usr_local_bin_path, PATH_MAX, "/usr/local/bin/%s", deleted_docker_binaries[i]) >= PATH_MAX) {
      aslLog(ASL_LEVEL_CRIT, "Failed to construct path to /usr/local/bin/%s", deleted_docker_binaries[i]);
      goto out2;
    }
    if (lstat(usr_local_bin_path, &buf) == 0) {
      if ((buf.st_mode & S_IFMT) == S_IFLNK) {
        if (unlink(usr_local_bin_path) == -1) {
          aslLog(ASL_LEVEL_CRIT, "Failed to delete old docker symlink %s: %s", usr_local_bin_path, strerror (errno));
          /* Let's be nice and not make this a hard failure */
        } else {
          aslLog(ASL_LEVEL_CRIT, "Deleted old docker symlink %s", usr_local_bin_path);
        }
      }
    }
  }


  for (unsigned long i = 0; i < sizeof(docker_binaries) / sizeof(docker_binaries[0]); i++ ) {
    /* set group_path: path to binary in the container folder */
    if (snprintf(group_path, PATH_MAX, "%s/%s/bin/%s", pwd.pw_dir, GROUP_CONTAINERS, docker_binaries[i]) >= PATH_MAX) {
      aslLog(ASL_LEVEL_CRIT, "Failed to construct path to %s", docker_binaries[i]);
      goto out2;
    }
    /* set usr_local_bin_path: path to binary in /usr/local/bin */
    if (snprintf(usr_local_bin_path, PATH_MAX, "/usr/local/bin/%s", docker_binaries[i]) >= PATH_MAX) {
      aslLog(ASL_LEVEL_CRIT, "Failed to construct path to /usr/local/bin/%s", docker_binaries[i]);
      goto out2;
    }
    /* set usr_local_bin_backup_path: path to backup existent binaries */
    if (snprintf(usr_local_bin_backup_path, PATH_MAX, "/usr/local/bin/%s.backup", docker_binaries[i]) >= PATH_MAX) {
      aslLog(ASL_LEVEL_CRIT, "Failed to construct path to /usr/local/bin/%s.backup", docker_binaries[i]);
      goto out2;
    }
    if (lstat(usr_local_bin_path, &buf) != 0) {
      if (errno != ENOENT) {
        aslLog(ASL_LEVEL_CRIT, "Failed to stat %s symlink: %s", usr_local_bin_path, strerror (errno));
        goto out2;
      }
    } else {
      if ((buf.st_mode & S_IFMT) == S_IFLNK) {
        /* Delete it and recreate */
        aslLog(ASL_LEVEL_CRIT, "File %s is a symlink", usr_local_bin_path);
        if (unlink(usr_local_bin_path) != 0) {
          aslLog(ASL_LEVEL_CRIT, "Failed to unlink %s symlink: %s", usr_local_bin_path, strerror (errno));
          goto out2;
        }
      } else {
        /* It's a binary so rename it */
        unlink(usr_local_bin_backup_path);
        if (rename(usr_local_bin_path, usr_local_bin_backup_path) != 0) {
          aslLog(ASL_LEVEL_CRIT, "Failed to rename %s to %s: %s", usr_local_bin_path, usr_local_bin_backup_path, strerror (errno));
          goto out2;
        }
      }
    }
    if (symlink(&group_path[0], &usr_local_bin_path[0]) == -1) {
      aslLog(ASL_LEVEL_CRIT, "Failed to create symlink from %s to %s: %s", &usr_local_bin_path[0], &group_path[0], strerror (errno));
      goto out2;
    }
    if (lchown(&usr_local_bin_path[0], uid, gid) == -1) {
      aslLog(ASL_LEVEL_CRIT, "Failed to chown symlink %s to %d.%d: %s", &usr_local_bin_path[0], uid, gid, strerror (errno));
      goto out2;
    }

    /* create symlink from /var/run/docker.sock to to user's docker api endpoint socket (s60)
      suffixing pre-existing non-symlink/non-socket files with ".backup" */
    const char *var_run_docker_sock_path = "/var/run/docker.sock";
    const char *var_run_docker_sock_backup_path = "/var/run/docker.sock.backup";

    // <container folder>/s60 -> formerly /var/tmp/docker.sock"

    const char *var_tmp_docker_sock_path = NULL;
    if (container_path) {
      aslLog(ASL_LEVEL_NOTICE, "Client supplied container path %s", container_path);

      /* Is the directory owned by the user we're talking to? */
      if (stat(container_path, &buf) == -1) {
        aslLog(ASL_LEVEL_CRIT, "Failed to stat %s: %p", container_path);
        goto out2;
      }
      if (!S_ISDIR(buf.st_mode)) {
        aslLog(ASL_LEVEL_CRIT, "Path %s is not a directory", container_path);
        goto out2;
      }
      if (buf.st_uid != uid) {
        aslLog(ASL_LEVEL_CRIT, "Path %s is owned by uid %d, but we're talking to uid %d", container_path, buf.st_uid, uid);
        goto out2;
      }
      unsigned long len = strlen(container_path) + strlen("/s60") + 1;
      char *tmp = (char*)malloc(len);
      if (!tmp) {
        aslLog(ASL_LEVEL_CRIT, "Failed to allocate space for the container path");
        goto out2;
      }
      strlcpy(tmp, container_path, len);
      strlcat(tmp, "/s60", len);
      var_tmp_docker_sock_path = tmp;
    } else {
      var_tmp_docker_sock_path = apple_utils_get_file_path_in_container_folder(uid, "s60");
    }
    if (lstat(var_run_docker_sock_path, &buf) == -1) {
      if (errno != ENOENT) {
        aslLog(ASL_LEVEL_CRIT, "Failed to stat /var/run/docker.sock symlink: %s", strerror (errno));
        goto out2;
      }
    } else {
      if (S_ISLNK(buf.st_mode) || S_ISSOCK(buf.st_mode)) {
        /* Delete it and recreate */
        aslLog(ASL_LEVEL_CRIT, "File %s is a symlink or a socket", var_run_docker_sock_path);
        if (unlink(var_run_docker_sock_path) != 0) {
          aslLog(ASL_LEVEL_CRIT, "Failed to unlink %s symlink: %s", var_run_docker_sock_path, strerror (errno));
          goto out2;
        }
      } else {
        /* Some kind of binary? Rename it */
        if (rename(var_run_docker_sock_path, var_run_docker_sock_backup_path) != 0) {
          aslLog(ASL_LEVEL_CRIT, "Failed to rename %s to %s: %s", var_run_docker_sock_path, var_run_docker_sock_backup_path, strerror (errno));
          goto out2;
        }
      }
    }
    aslLog(ASL_LEVEL_NOTICE, "Symlinking %s to %s", var_run_docker_sock_path, var_tmp_docker_sock_path);

    if (symlink(var_tmp_docker_sock_path, var_run_docker_sock_path) == -1) {
      aslLog(ASL_LEVEL_CRIT, "Failed to create %s symlink: %s", var_run_docker_sock_path, strerror (errno));
    }
  }

  r = 0;
out2:
  (void)really_write(client_fd, &r, 1);
  free(buffer);
  exit(r);
}

void fork_one_proxy_on(int fd) {
  /* We want to 'fire and forget' our child processes */
  signal(SIGCHLD, SIG_IGN);

  pid_t child;
  int client_fd;
  struct sockaddr_un client_address;
  socklen_t address_len;

  if ((client_fd = accept(fd, (struct sockaddr*) &client_address, &address_len)) == -1){
    aslLog(ASL_LEVEL_ERR, "Failed to accept connection: %s", strerror (errno));
    return; /* not fatal */
  }
  child = fork();
  if (child != 0) {
    /* The parent will clean up and return, leaving the child to do the work. */
    if (verbose_flag) {
      aslLog(ASL_LEVEL_DEBUG, "Client has connected and handed to pid: %d", child);
    }
    close(client_fd);
    return;
  }
  struct init_message *me = create_init_message(); /* leaked */
  struct init_message you;
  struct ethernet_args eth;
  enum command command;
  if (read_init_message(client_fd, &you) == -1) {
    exit(1);
  }
  char *txt = print_init_message(&you);
  aslLog(ASL_LEVEL_NOTICE, "Client reports %s", txt);
  free(txt);
  switch (you.version){
    case 0:
      command = ethernet;
      break;
    case 1: /* The oldest version will be used for probe and uninstall */
    case CURRENT_VERSION: /* Hyperkit will always have the current version */
      if (write_init_message(client_fd, me) == -1) {
        exit(1);
      }
      if (read_command(client_fd, &command) == -1) {
        exit(1);
      }
      break;
    default:
      aslLog(ASL_LEVEL_CRIT, "Client from the past or future has connected: %d", you.version);
      /* Shouldn't happen because we will trigger a reinstall when the versions
         don't match exactly. */
      exit(1);
  }
  uint8_t r = 1;

  char *buffer = NULL;
  switch (command) {
    case ethernet:
    {
      if (read_ethernet_args(client_fd, &eth) == -1) {
        exit(1);
      }
      struct vmnet_state *vms = open_vmnet(eth.uuid_string);
      vms->client_fd = client_fd;
      /* Send the MAC and MTU to the client */
      if (write_vif_info(client_fd, &vms->vif) == -1) {
        exit(1);
      }
      proxy_forever(vms);
      break;
    /* uninstall removes vmnetd from launchd */
    }
    case uninstall:
    {
      aslLog(ASL_LEVEL_CRIT, "client requests com.docker.vmnetd uninstall");
      /* Join a new process group temporarily so launchd doesn't kill us */
      pid_t pgid = setsid();
      if (pgid == -1) {
        aslLog(ASL_LEVEL_CRIT, "failed to start new process group: %s", strerror (errno));
        goto out;
      }
      if (system("/bin/launchctl unload " DOCKER_VMNETD_PLIST) != 0) {
        aslLog(ASL_LEVEL_CRIT, "failed to unload com.docker.vmnetd plist %s: %s", DOCKER_VMNETD_PLIST, strerror (errno));
        goto out;
      }
      if (unlink(DOCKER_VMNETD_PLIST) != 0) {
        aslLog(ASL_LEVEL_CRIT, "failed to unlink com.docker.vmnetd plist %s: %s", DOCKER_VMNETD_PLIST, strerror (errno));
        goto out;
      }
      if (unlink(DOCKER_VMNETD_BINARY) != 0) {
        aslLog(ASL_LEVEL_CRIT, "failed to unlink com.docker.vmnetd binary %s: %s", DOCKER_VMNETD_BINARY, strerror (errno));
        goto out;
      }
      aslLog(ASL_LEVEL_CRIT, "com.docker.vmnetd uninstall complete");
      r = 0;
out:
      (void)really_write(client_fd, &r, 1);

      /* vmnet is now uninstalled, let's remove its socket */
      char *unix_socket_path = "/var/tmp/com.docker.vmnetd.socket";
      struct stat buf;
      if (stat(unix_socket_path, &buf) != 0){
        if (errno != ENOENT) {
          aslLog(ASL_LEVEL_CRIT, "Failed to stat %s : %s", unix_socket_path, strerror (errno));
          goto out4;
        }
      } else {
        if (unlink(unix_socket_path) == -1) {
          if (errno != ENOENT) {
            aslLog(ASL_LEVEL_CRIT, "Failed to delete %s: %s", unix_socket_path, strerror (errno));
            goto out4;
          }
        } else {
          aslLog(ASL_LEVEL_CRIT, "Deleted %s", unix_socket_path);
        }
      }

      exit(r);
    }


    case install_symlinks:
    {
      really_install_symlinks(client_fd, NULL);
    }
    case install_symlinks2:
    {
      struct install_symlinks2 sym;
      if (read_install_symlinks2(client_fd, &sym) == -1) {
        exit(1);
      }
      /* The char[1024] may not be NULL terminated: let's copy it and terminate */
      char *tmp = malloc(sizeof(sym.container) + 1);
      if (!tmp) {
        aslLog(ASL_LEVEL_CRIT, "Failed to allocate space for the container path");
        exit(1);
      }
      memcpy(tmp, &sym.container, sizeof(sym.container));
      *(tmp + sizeof(sym.container)) = '\000';

      really_install_symlinks(client_fd, tmp);
    }

  /* uninstall_symlinks uninstalls everything that's been
    installed by install_symlinks */
  case uninstall_symlinks:
  {
      uid_t uid;
      gid_t gid;
      if (getpeereid(client_fd, &uid, &gid) != 0) {
        aslLog(ASL_LEVEL_CRIT, "Failed to discover the uid/gid of the peer");
        goto out3;
      }
      aslLog(ASL_LEVEL_CRIT, "client with uid %d requests com.docker.vmnetd uninstall symlinks", uid);
      /* the example from getpwent(3) */
      long bufsize;
      if ((bufsize = sysconf(_SC_GETPW_R_SIZE_MAX)) == -1)
        abort();
      if ((buffer = malloc((size_t)bufsize)) == NULL){
        aslLog(ASL_LEVEL_CRIT, "Failed to allocate buffer of size %ld", bufsize);
        goto out3;
      }
      struct passwd pwd, *pwd_result = NULL;
      if (getpwuid_r(uid, &pwd, buffer, (size_t)bufsize, &pwd_result) != 0 || !pwd_result) {
        aslLog(ASL_LEVEL_CRIT, "Failed to discover the homedir of the peer");
        goto out3;
      }

      char usr_local_bin_path[PATH_MAX];
      char usr_local_bin_backup_path[PATH_MAX];
      struct stat buf;

      /* Try to find old binary symlinks that should have been deleted and delete them.
      These binaries are listed in deleted_docker_binaries */
      for (unsigned long i = 0; i < sizeof(deleted_docker_binaries) / sizeof(deleted_docker_binaries[0]); i++ ) {
        if (snprintf(usr_local_bin_path, PATH_MAX, "/usr/local/bin/%s", deleted_docker_binaries[i]) >= PATH_MAX) {
          aslLog(ASL_LEVEL_CRIT, "Failed to construct path to /usr/local/bin/%s", deleted_docker_binaries[i]);
          goto out3;
        }
        if (lstat(usr_local_bin_path, &buf) == 0) {
          if ((buf.st_mode & S_IFMT) == S_IFLNK) {
            if (unlink(usr_local_bin_path) == -1) {
              aslLog(ASL_LEVEL_CRIT, "Failed to delete old docker symlink %s: %s", usr_local_bin_path, strerror (errno));
              /* Let's be nice and not make this a hard failure */
            } else {
              aslLog(ASL_LEVEL_CRIT, "Deleted old docker symlink %s", usr_local_bin_path);
            }
          }
        }
      }


      for (unsigned long i = 0; i < sizeof(docker_binaries) / sizeof(docker_binaries[0]); i++ ) {

        /* set usr_local_bin_path: path to symlinks droped in /usr/local/bin */
        if (snprintf(usr_local_bin_path, PATH_MAX, "/usr/local/bin/%s", docker_binaries[i]) >= PATH_MAX) {
          aslLog(ASL_LEVEL_CRIT, "Failed to construct path to /usr/local/bin/%s", docker_binaries[i]);
          goto out3;
        }
        /* set usr_local_bin_backup_path: path where existent binaries have been backed up */
        if (snprintf(usr_local_bin_backup_path, PATH_MAX, "/usr/local/bin/%s.backup", docker_binaries[i]) >= PATH_MAX) {
          aslLog(ASL_LEVEL_CRIT, "Failed to construct path to /usr/local/bin/%s.backup", docker_binaries[i]);
          goto out3;
        }

        /* unlink symlinks */
        if (lstat(usr_local_bin_path, &buf) != 0) {
          if (errno != ENOENT) {
            aslLog(ASL_LEVEL_CRIT, "Failed to stat %s symlink: %s", usr_local_bin_path, strerror (errno));
            goto out3;
          }
        } else {
          if ((buf.st_mode & S_IFMT) == S_IFLNK) {
            aslLog(ASL_LEVEL_CRIT, "File %s is a symlink", usr_local_bin_path);
            if (unlink(usr_local_bin_path) != 0) {
              aslLog(ASL_LEVEL_CRIT, "Failed to unlink %s symlink: %s", usr_local_bin_path, strerror (errno));
              goto out3;
            }
          }
        }

        /* restore backup files */
        if (stat(usr_local_bin_backup_path, &buf) != 0) {
          if (errno != ENOENT) {
            aslLog(ASL_LEVEL_CRIT, "Failed to stat %s backup: %s", usr_local_bin_backup_path, strerror (errno));
            goto out3;
          }
        } else {
            if (rename(usr_local_bin_backup_path, usr_local_bin_path) != 0) {
              aslLog(ASL_LEVEL_CRIT, "Failed to rename %s to %s: %s", usr_local_bin_backup_path, usr_local_bin_path, strerror (errno));
              goto out3;
            }
        }
      }

      /* unlink /var/run/docker.sock
        and restore backup*/
      const char *var_run_docker_sock_path = "/var/run/docker.sock";
      const char *var_run_docker_sock_backup_path = "/var/run/docker.sock.backup";
      if (lstat(var_run_docker_sock_path, &buf) == -1) {
        if (errno != ENOENT) {
          aslLog(ASL_LEVEL_CRIT, "Failed to stat /var/run/docker.sock symlink: %s", strerror (errno));
          goto out3;
        }
      } else {
        if (S_ISLNK(buf.st_mode) || S_ISSOCK(buf.st_mode)) {
          aslLog(ASL_LEVEL_CRIT, "File %s is a symlink or a socket", var_run_docker_sock_path);
          if (unlink(var_run_docker_sock_path) != 0) {
            aslLog(ASL_LEVEL_CRIT, "Failed to unlink %s symlink: %s", var_run_docker_sock_path, strerror (errno));
            goto out3;
          }
        }
      }

      /* restore backup */
      if (stat(var_run_docker_sock_backup_path, &buf) != 0) {
        if (errno != ENOENT) {
          aslLog(ASL_LEVEL_CRIT, "Failed to stat %s backup: %s", var_run_docker_sock_backup_path, strerror (errno));
          goto out3;
        }
      } else {
          if (rename(var_run_docker_sock_backup_path, var_run_docker_sock_path) != 0) {
            aslLog(ASL_LEVEL_CRIT, "Failed to rename %s to %s: %s", var_run_docker_sock_backup_path, var_run_docker_sock_path, strerror (errno));
            goto out3;
          }
      }

    r = 0;
out3:
      (void)really_write(client_fd, &r, 1);
      free(buffer);
      exit(r);
    }
    /* uninstall_sockets uninstalls all sockets
    but /var/tmp/com.docker.vmnetd.socket */
    case uninstall_sockets:
    {
      /* Try to find old binary symlinks that should have been deleted and delete them.
      These binaries are listed in deleted_docker_binaries */
      for (unsigned long i = 0; i < sizeof(sockets_to_uninstall) / sizeof(sockets_to_uninstall[0]); i++ ) {
        // NOTE(aduermael): we don't need to check for errors here
        // it generates too much logs, as we're trying to delete everything
        // even sockets that may not be created by the app anymore.
        rmrf(sockets_to_uninstall[i]);
      }
      r = 0;
out4:
      (void)really_write(client_fd, &r, 1);
      exit(r);
    }
    case bind_ipv4:
    {
      struct bind_ipv4 ip;
      if (read_bind_ipv4(client_fd, &ip) != 0) {
        aslLog(ASL_LEVEL_CRIT, "Failed to receive IPv4 bind request");
        exit(1);
      }
      /* We will send this socket to the client: */
      int s = socket(AF_INET, (ip.stream == 0)?SOCK_STREAM:SOCK_DGRAM, 0);
      {
        int enable = 1;
        if (setsockopt(s, SOL_SOCKET, SO_REUSEADDR, &enable, sizeof(int)) < 0) {
          aslLog(ASL_LEVEL_CRIT, "Failed to set SO_REUSEADDR: %m");
          exit(1);
        }
      }
      /* The message will have a single result value: */
      struct msghdr msg;
      struct iovec vec;
      char iobuf[8]; /* send back an errno value as a uint64-t */
      msg.msg_name = NULL;
      msg.msg_namelen = 0;
      vec.iov_base=iobuf;
      vec.iov_len=sizeof(iobuf);
      msg.msg_iov=&vec;
      msg.msg_iovlen=1;
      /* and a control message including the fd: */
      struct cmsghdr *cmsg;
      union {
        /* ensure suitable alignment */
        char buf[100];
        struct cmsghdr align;
      } u;
      assert(CMSG_SPACE(sizeof(int)) < sizeof(u.buf));
      msg.msg_control = u.buf;
      msg.msg_controllen = CMSG_SPACE(sizeof(int));
      msg.msg_flags = 0;
      cmsg = CMSG_FIRSTHDR(&msg);
      cmsg->cmsg_level = SOL_SOCKET;
      cmsg->cmsg_type = SCM_RIGHTS;
      cmsg->cmsg_len = CMSG_LEN(sizeof(s));
      int *fdptr = (int *) (void*) CMSG_DATA(cmsg);
      *fdptr = s;

      struct sockaddr_in sockaddr;
      bzero((char *) &sockaddr, sizeof(sockaddr));
      sockaddr.sin_family = AF_INET;
      sockaddr.sin_addr.s_addr = htonl(ip.ipv4);
      sockaddr.sin_port = htons(ip.port);
      int64_t result = 0;
      if (bind(s, (struct sockaddr *) &sockaddr, sizeof(sockaddr)) == -1) {
        aslLog(ASL_LEVEL_CRIT, "Failed to bind IPv4 %"PRId32" %d: %m", ip.ipv4, ip.port);
        result = errno;
      }
      memcpy(&iobuf[0], &result, sizeof(int64_t));
      if (sendmsg(client_fd, &msg, 0) == -1) {
        aslLog(ASL_LEVEL_CRIT, "Failed to send fd: %m");
        exit(1);
      }
      aslLog(ASL_LEVEL_CRIT, "sendmsg ok");
      exit(0);
    }
  } // end switch
}

void test_open_vmnet(){
  open_vmnet(NULL);
  aslLog(ASL_LEVEL_NOTICE, "Successfully opened com.docker.vmnetd.");
}
