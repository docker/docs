//
//  utils.c
//  sdk
//
//  Created by Gaetan de Villele on 4/6/16.
//  Copyright © 2016 docker. All rights reserved.
//

#include "utils.h"
// C
#include <stdio.h>


void pplog(const char* format, ...) {
    va_list args;
    va_start(args, format);
    // vprintf is a printf that understand va_list
    vprintf(format, args);
    printf("\n");
    va_end(args);
}