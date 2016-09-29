//
//  main.c
//  sdk
//
//  Created by Gaetan de Villele on 1/28/16.
//  Copyright © 2016 docker. All rights reserved.
//

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
// SDK
#include "sdk.h"
#include "error.h"
#include "utils.h"

// main function
int main(int argc, const char * argv[]) {
    
//    if (argc != 3) {
//        pplog("wrong use.\n-> first argument is the action ('greeting' or 'goodbye')\n-> second argument is your name");
//        return EXIT_FAILURE;
//    }
//    const char* action = argv[1];
//    const char* name = argv[2];
    
    // initializing the SDK
    if (pp_init() != PP_OK) {
        pplog("😡 FAILED TO INITIALIZE SDK...");
        return EXIT_FAILURE;
    }
    
    // socket path
    const char *homeDir = getenv("HOME");
    if (homeDir == NULL) {
        pplog("failed to find socket in ~/Desktop/test_socket.sock");
        return EXIT_FAILURE;
    }
    
    size_t socketPathLen = 25 + strlen(homeDir) + 1;
    char* socketPath = (char*) malloc(socketPathLen);
    snprintf(socketPath, socketPathLen, "%s/Desktop/test_socket.sock", homeDir);
    // pplog("socket path: %s", socketPath);
    
    // connecting to server
    uint32_t connection_handle;
    if (pp_connect(socketPath, &connection_handle) != PP_OK) {
        pplog("😡 FAILED TO CONNECT...");
        return EXIT_FAILURE;
    }
    //pplog("🍀 we have a connection [%u]", connection_handle);
    
    // constructing the request
    const char* action = "greeting";
    const char* name = "John Doe";
    size_t request_size = 36 + strlen(action) + strlen(name) + 1;
    char* request = (char*) malloc(request_size);
    snprintf(request, request_size, "{\"Action\":\"%s\",\"Content\":{\"Name\":\"%s\"}}", action, name);
    pplog("> sending request: %s", request);
    
    // sending a request
    char* response = NULL;
    if (pp_request(connection_handle, (char*)request, &response) != PP_OK) {
        pplog("😡 FAILED TO SEND REQUEST...");
        return EXIT_FAILURE;
    }
    //pplog("🍀 we sent a request");
    pplog("response from server:\n%s", response);
    pp_response_free(response);

    if (pp_disconnect(connection_handle) != PP_OK) {
        pplog("😡 FAILED TO DISCONNECT...");
        return EXIT_FAILURE;
    }
//    pplog("🍀 we disconnected");
    return EXIT_SUCCESS;
}
