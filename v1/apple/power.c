#include <asl.h>
#include <ctype.h>
#include <stdlib.h>
#include <stdio.h>

#include <mach/mach_port.h>
#include <mach/mach_interface.h>
#include <mach/mach_init.h>

#include <IOKit/pwr_mgt/IOPMLib.h>
#include <IOKit/IOMessage.h>

#include "asl_logger.h"

/* based on https://developer.apple.com/library/mac/qa/qa1340/_index.html */

static io_connect_t  root_port; // a reference to the Root Power Domain IOService

struct connection {
  int reader;
  int writer;
};

void
apple_power_callback(void *data, io_service_t service, natural_t messageType, void *messageArgument )
{
    struct connection *c = (struct connection*)data;
    pid_t me = getpid();
    uint8_t input;
    switch (messageType) {
    case kIOMessageSystemWillSleep:
      write(c->writer, "S", 1);
      /* Wait for an ack before allowing the power change */
      read(c->reader, &input, 1);
      IOAllowPowerChange( root_port, (long)messageArgument );
      break;

    case kIOMessageSystemWillPowerOn:
      write(c->writer, "W", 1);
      break;
    }
}

void listen_for_power_events(int writer, int reader) {
    IONotificationPortRef  notifyPortRef;
    io_object_t            notifierObject;
    /* Allocate but never free the fd */
    void *data = malloc(sizeof(struct connection));
    if (data == NULL) {
      apple_asl_logger_log(ASL_LEVEL_CRIT, "unable to allocate memory for new connection struct");
      abort();
    }

    struct connection *c = (struct connection*)data;
    c->writer = writer;
    c->reader = reader;

    root_port = IORegisterForSystemPower(data, &notifyPortRef, apple_power_callback, &notifierObject );
    if ( root_port == 0 ) {
        printf("IORegisterForSystemPower failed\n");
        return;
    }

    CFRunLoopAddSource( CFRunLoopGetCurrent(),
            IONotificationPortGetRunLoopSource(notifyPortRef), kCFRunLoopCommonModes );

    CFRunLoopRun();
}
