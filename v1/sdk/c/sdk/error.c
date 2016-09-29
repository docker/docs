//
//  error.c
//  AgentSdk
//
//  Created by Gaetan de Villele on 4/6/16.
//  Copyright © 2016 docker. All rights reserved.
//

#include "error.h"

const uint32_t PP_OK = 0;

// sdk error codes
const uint32_t PP_ERROR_UNKNOWN = 1;
const uint32_t PP_ERROR_SDK_ALREADY_INITIALIZED = 2;
const uint32_t PP_ERROR_SDK_NOT_INITIALIZED = 3;

// ssh error codes
const uint32_t PP_ERROR_SSH = 200;
const uint32_t PP_ERROR_SSH_INIT = 201;