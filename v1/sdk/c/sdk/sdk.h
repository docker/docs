//
//  sdk.h
//  sdk
//
//  Created by Gaetan de Villele on 2/3/16.
//  Copyright © 2016 docker. All rights reserved.
//

//
// This is the main header of the SDK C client library for the Pinata Protocol (aka PP)
//

#ifndef sdk_h
#define sdk_h

#ifdef _WIN32
#define EXPORTSYMBOL __declspec( dllexport )
#else
#define EXPORTSYMBOL
#endif

#include <stdint.h>

#include "error.h"

/// initializes PP client library.
/// @return PP_OK or a PP_ERROR code
EXPORTSYMBOL uint32_t pp_init();

/// connects to a PP server.
/// @param socket_path
///        path of the socket to connect to
/// @param conn_handle
///        pointer on unsigned which will contain the connection handle
/// @return PP_OK or a PP_ERROR code
EXPORTSYMBOL uint32_t pp_connect(const char *socket_path,
                                 //const char *knownhosts_path,
                                 //const char *identity_path,
                                 //const char *ssh_dir_path,
                                 uint32_t *conn_handle);

/// closes a connection to a PP server and free the connection. The connection
/// handle must not be used after calling this function.
/// @param conn_handle
///        handle of the connection to close
/// @return PP_OK or a PP_ERROR code
EXPORTSYMBOL uint32_t pp_disconnect(uint32_t conn_handle);

/// sends a request to the server and get the response.
/// @param conn_handle
///        handle to the connection on which the request must be sent
/// @param request
///        content of the request (null-terminated string)
/// @param response
///        content of the response (pointer on null-terminated string)
/// @return PP_OK or a PP_ERROR code
EXPORTSYMBOL uint32_t pp_request(uint32_t conn_handle, const char* request, char** response);

/// free response
/// @param response
///        a pointer to the response buffer to free
EXPORTSYMBOL void pp_response_free(char* response);

#endif /* sdk_h */
