//
//  sdk.c
//  sdk
//
//  Created by Gaetan de Villele on 2/3/16.
//  Copyright © 2016 docker. All rights reserved.
//

#include "sdk.h"

// SDK
#include "connection.h"
#include "error.h"
#include "utils.h"
// C
#include <stdlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <sys/un.h> // for sockaddr_un
// SSH
#include "libssh/libssh.h"


// indicates whether the sdk is initialized
bool initialized = false;


uint32_t pp_init() {
    // If the sdk is already initialized, we return an error.
    if (initialized == true) {
        return PP_ERROR_SDK_ALREADY_INITIALIZED;
    }
    // This is initializing crypto data structures. It can be omitted if the
    // program using libssh is not multithreaded, but we call it to be safe.
    if (ssh_init() != SSH_OK) {
        return PP_ERROR_SSH_INIT;
    }
    initialized = true;
    return PP_OK;
}

uint32_t pp_connect(const char *socket_path, uint32_t *conn_handle) {
    // If the sdk is not initialized, we return an error.
    if (initialized == false) {
        return PP_ERROR_SDK_NOT_INITIALIZED;
    }
    
    // create a new ppconnection
    ppconnection *conn = connection_new();
    if (conn == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    
    // alloc ssh session
    conn->session = ssh_new();
    if (conn->session == NULL) {
        pplog("error: failed to alloc ssh session");
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
    // create socket
    int socketFD = socket(AF_UNIX, SOCK_STREAM, PF_UNSPEC);
    if (socketFD < 0) {
        pplog("error: failed to socket()");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    
    // connect to socket at the given path
    struct sockaddr_un address;
    // fill the struct with 0s
    memset(&address, 0, sizeof(struct sockaddr_un));
    // set the socket family
    address.sun_family = AF_UNIX;
    // set the socket path
    snprintf(address.sun_path, 104, "%s", socket_path);
    if(connect(socketFD, (struct sockaddr*) &address, sizeof(struct sockaddr_un)) != 0)
    {
        pplog("error: socket connect() failed");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    // tell libssh to use the socket
    if (ssh_options_set(conn->session, SSH_OPTIONS_FD, &socketFD) != SSH_OK) {
        pplog("error: failed to set session param socket");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
//    if (ssh_options_set(conn->session, SSH_OPTIONS_HOST, "coucou") != SSH_OK) {
//        pplog("error: failed to set hostname");
//        connection_close(conn);
//        connection_free(conn);
//        return PP_ERROR_SSH; // TODO: gdevillele:
//    }
//    // set the custom known_hosts file path so libssh is not reading the default
//    // one itself
//    if (ssh_options_set(conn->session, SSH_OPTIONS_KNOWNHOSTS, "coucou") != SSH_OK) {
//        pplog("error: failed to set session known_host path: %s", ssh_get_error(conn->session));
//        connection_close(conn);
//        connection_free(conn);
//        return PP_ERROR_SSH; // TODO: gdevillele:
//    }
    
    // since we are passing NULL as username, the default libssh username will be used
    if (ssh_options_set(conn->session, SSH_OPTIONS_USER, NULL) < 0) {
        pplog("failed to set username option");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
    // set verbosity
    int verbosity = SSH_LOG_NOLOG; //SSH_LOG_PROTOCOL SSH_LOG_PACKET SSH_LOG_FUNCTIONS
    if (ssh_options_set(conn->session, SSH_OPTIONS_LOG_VERBOSITY, &verbosity) != SSH_OK) {
        pplog("ssh session set option SSH_OPTIONS_LOG_VERBOSITY failed");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
    // actually connect session
    if(ssh_connect(conn->session) != SSH_OK){
        pplog("failed to connect ssh session: %s", ssh_get_error(conn->session));
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
    // check if session is connected
    if (!ssh_is_connected(conn->session)) {
        pplog("ssh failed to connect...");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
    // TODO: gdevillele: verify server's identity (later)
    
    // set authentication mode ("none" method for now)
    if (ssh_userauth_none(conn->session, NULL) != SSH_AUTH_SUCCESS) {
        pplog("ssh session failed to authenticate");
        connection_close(conn);
        connection_free(conn);
        return PP_ERROR_SSH; // TODO: gdevillele:
    }
    
    uint32_t handle = connection_obtain_handle(conn);
    *conn_handle = handle;
    return PP_OK;
}

uint32_t pp_disconnect(uint32_t conn_handle) {
    // If the sdk is not initialized, we return an error.
    if (initialized == false) {
        return PP_ERROR_SDK_NOT_INITIALIZED;
    }
    // retrieve connection using handle
    ppconnection* conn = connection_by_handle(conn_handle);
    if (conn == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    if (connection_close(conn) != PP_OK) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    if (connection_free(conn) != PP_OK) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    connection_recycle_handle(conn_handle);
    return PP_OK;
}

uint32_t pp_request(uint32_t conn_handle, const char* request, char** response) {
    // If the sdk is not initialized, we return an error.
    if (initialized == false) {
        return PP_ERROR_SDK_NOT_INITIALIZED;
    }
    // retrieve connection using handle
    ppconnection* conn = connection_by_handle(conn_handle);
    if (conn == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    if (conn->session == NULL || !ssh_is_connected(conn->session)) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    if (request == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    if (response == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    if (!ssh_is_connected(conn->session)) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    
    // open ssh channel
    ssh_channel channel = NULL;
    channel = ssh_channel_new(conn->session);
    if (channel == NULL) {
        pplog("failed to create ssh channel for request");
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    // TODO: use the payload to indicate the type of channel in the PP protocol
    //       examples:
    //       - "json-request"
    //       - "json-stream"
    //       - "protobuf-request"
    if (ssh_channel_open_custom(channel, "com.docker.api.pp", NULL) < 0) {
        pplog("failed to open channel: %s", ssh_get_error(conn->session));
        ssh_channel_close(channel); // TODO: check for error
        ssh_channel_free(channel);
        channel = NULL;
        return PP_ERROR_UNKNOWN;
        // TODO: gdevillele: eventually handle retry here
    }
    if (!ssh_channel_is_open(channel)) {
        pplog("channel has not been opened properly");
        ssh_channel_close(channel); // TODO: check for error
        ssh_channel_free(channel);
        channel = NULL;
        return PP_ERROR_UNKNOWN;
    }
    // channel is open
    
    // request
    const uint32_t request_size = (uint32_t)strlen(request);
    uint32_t request_bytes_sent = 0;
    // response
    const uint32_t RESPONSE_READ_BUF_SIZE = 1;
    uint32_t response_size = 0;
    char* response_buffer = NULL;
    char response_read_buf[RESPONSE_READ_BUF_SIZE];
    
    while (request_bytes_sent < request_size) {
        // try to write the remaining bytes of the request to the channel
        int bytes_written = ssh_channel_write(channel, request+request_bytes_sent, request_size-request_bytes_sent);
        if (bytes_written == SSH_ERROR) {
            pplog("failed to write request into ssh channel");
            ssh_channel_close(channel); // TODO: check for error
            ssh_channel_free(channel);
            channel = NULL;
            return PP_ERROR_SSH;
        }
        // pplog("bytes written: %d", bytes_written);
        request_bytes_sent += bytes_written;
    }
    // tell the ssh server the request has been sent in whole by writting EOF
    // in channel. From now on, no bytes can be written in this channel.
    if (ssh_channel_send_eof(channel) == SSH_ERROR) {
        pplog("failed to write EOF in channel");
        ssh_channel_close(channel); // TODO: check for error
        ssh_channel_free(channel);
        channel = NULL;
        return PP_ERROR_SSH;
    }
    // pplog("request sent through channel");
    
    // read from channel until EOF
    while (!ssh_channel_is_eof(channel)) {
        int rc = ssh_channel_read(channel, response_read_buf, RESPONSE_READ_BUF_SIZE, 0);
        if (rc > 0) {
            response_buffer = realloc(response_buffer, response_size + rc);
            memcpy(response_buffer + response_size, response_read_buf, rc);
            response_size += rc;
        } else if (rc == 0) {
            // means EOF
            break;
        } else if (rc == SSH_ERROR) {
            pplog("error while reading from channel");
            ssh_channel_close(channel); // TODO: check for error
            ssh_channel_free(channel);
            channel = NULL;
            return PP_ERROR_SSH;
        }
    }
    
    // add null character at the end of response string
    // TODO: (gdevillele) this will be removed when we will use C buffers instead of C strings
    response_buffer = realloc(response_buffer, response_size + 1);
    memcpy(response_buffer + response_size, "\0", 1);
    response_size += 1;
    
    if (ssh_channel_close(channel) == SSH_ERROR) {
        pplog("failed to close channel");
        ssh_channel_free(channel);
        channel = NULL;
        return PP_ERROR_SSH;
    }
    
    *response = response_buffer;
    
    return PP_OK;
}

void pp_response_free(char* response) {
    free(response);
}

#pragma mark - private -


