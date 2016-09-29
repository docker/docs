//
//  error.h
//  AgentSdk
//
//  Created by Gaetan de Villele on 4/6/16.
//  Copyright © 2016 docker. All rights reserved.
//

#ifndef error_h
#define error_h

// C
#include <stdint.h>

extern const uint32_t PP_OK;

// sdk error codes
extern const uint32_t PP_ERROR_UNKNOWN;
extern const uint32_t PP_ERROR_SDK_ALREADY_INITIALIZED;
extern const uint32_t PP_ERROR_SDK_NOT_INITIALIZED;

// ssh error codes
extern const uint32_t PP_ERROR_SSH;
extern const uint32_t PP_ERROR_SSH_INIT;

#endif /* error_h */
