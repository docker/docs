//
//  connection.c
//  sdk
//
//  Created by Gaetan de Villele on 4/7/16.
//  Copyright © 2016 docker. All rights reserved.
//

#include "connection.h"
// C
#include <stdlib.h>
// SDK
#include "error.h"



#pragma mark - private -

// dynamic array of open connections
ppconnection **connections = NULL;
uint32_t connections_count = 0;



#pragma mark - exposed -

ppconnection* connection_new() {
    ppconnection *new_conn = (ppconnection*) malloc(sizeof(ppconnection));
    if (new_conn != NULL) {
        new_conn->session = NULL;
    }
    return new_conn;
}

// TODO: gdevillele: return error if a new handle cannot be generated (more than 4.2B connections for example)
uint32_t connection_obtain_handle(ppconnection *conn) {
    
    // try to find a spot available in the connection array
    for (uint32_t i = 0; i < connections_count; i++) {
        if (connections[i] == NULL) {
            // spot available
            connections[i] = conn;
            return i;
        }
    }
    
    // if no spot is available in the array, then we increase the size of the array
    connections = realloc(connections, sizeof(void*)*(connections_count+1));
    connections_count++;
    uint32_t handle = connections_count-1;
    connections[handle] = conn;
    return handle;
}

ppconnection* connection_by_handle(uint32_t conn_handle) {
    // test if handle in in connection array
    if (conn_handle > (connections_count-1)) {
        // handle is invalid
        return NULL;
    }
    return connections[conn_handle];
}

uint32_t connection_close(ppconnection *conn) {
    if (conn == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    // if there is no ssh session, we consider it a success
    if (conn->session == NULL) {
        return PP_OK;
    }
    // if session is connected, we disconnect
    if (ssh_is_connected(conn->session)) {
        ssh_disconnect(conn->session);
    }
    // in any case we free the ssh session
    ssh_free(conn->session);
    conn->session = NULL;
    return PP_OK;
}

uint32_t connection_free(ppconnection *conn) {
    if (conn == NULL) {
        return PP_ERROR_UNKNOWN; // TODO: gdevillele:
    }
    free(conn);
    return PP_OK;
}

void connection_recycle_handle(uint32_t conn_handle) {
    connections[conn_handle] = NULL;
}
