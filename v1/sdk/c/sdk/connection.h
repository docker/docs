//
//  connection.h
//  sdk
//
//  Created by Gaetan de Villele on 4/7/16.
//  Copyright © 2016 docker. All rights reserved.
//

#ifndef connection_h
#define connection_h

#include "libssh/libssh.h"

/// structure representing a PP connection
typedef struct ppconnection {
    ssh_session session;
} ppconnection;

#pragma mark - Function protoypes -

/// alloc and init new connection
ppconnection* connection_new();

/// obtain handle for connection
uint32_t connection_obtain_handle(ppconnection *conn);

/// obtain connection by handle
ppconnection* connection_by_handle(uint32_t conn_handle);

/// close connection (including the ssh session)
uint32_t connection_close(ppconnection *conn);

/// free connection
uint32_t connection_free(ppconnection *conn);

/// recycle connection handle
void connection_recycle_handle(uint32_t conn_handle);

#endif /* connection_h */
