//
//  ssh.c
//  AgentSdk
//
//  Created by Gaetan de Villele on 4/6/16.
//  Copyright © 2016 docker. All rights reserved.
//

#include "ssh.h"

#include "error.h"


//uint32_t pp_ssh_connect() {

//    // test if the sdk is already connected
//    if (ssh_is_connected(_session)) {
//        set_last_error("sdk is already connected");
//        return PINATA_ERROR;
//    }
//    
//    // alloc ssh session
//    _session = ssh_new();
//    if (_session == NULL) {
//        set_last_error("failed to alloc ssh session");
//        return PINATA_ERROR;
//    }
//    
//    // set server address
//    if (ssh_options_set(_session, SSH_OPTIONS_HOST, _address) < 0) {
//        char* error_optionsethost = strdyncat("ssh session set option SSH_OPTIONS_HOST failed: ", ssh_get_error(_session));
//        set_last_error(error_optionsethost);
//        free(error_optionsethost);
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    // set server port
//    unsigned int port_copy = _port;
//    if (ssh_options_set(_session, SSH_OPTIONS_PORT, &port_copy) < 0) {
//        set_last_error("ssh session set option SSH_OPTIONS_PORT failed");
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    // set user (user is "pinata" for now)
//    if (ssh_options_set(_session, SSH_OPTIONS_USER, "pinata") < 0) {
//        set_last_error("ssh session set option SSH_OPTIONS_USER failed");
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    // set verbosity
//    int verbosity = SSH_LOG_NOLOG; //SSH_LOG_PROTOCOL; // SSH_LOG_PACKET // SSH_LOG_FUNCTIONS
//    if (ssh_options_set(_session, SSH_OPTIONS_LOG_VERBOSITY, &verbosity) < 0) {
//        set_last_error("ssh session set option SSH_OPTIONS_LOG_VERBOSITY failed");
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    // connect session
//    if(ssh_connect(_session) == SSH_ERROR){
//        char* error_sessionconnect = strdyncat("ssh session failed to connect: ", ssh_get_error(_session));
//        set_last_error(error_sessionconnect);
//        free(error_sessionconnect);
//        ssh_disconnect(_session);
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    // authenticate
//    if (ssh_userauth_none(_session, "pinata") != SSH_AUTH_SUCCESS) {
//        set_last_error("ssh session failed to authenticate");
//        ssh_disconnect(_session);
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    if (!ssh_is_connected(_session)) {
//        set_last_error("ssh failed to connect to agent");
//        ssh_free(_session);
//        _session = NULL;
//        return PINATA_ERROR;
//    }
//    
//    return PP_OK;
//}