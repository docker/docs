#import "BITCrashExceptionApplication.h"

#import <sys/sysctl.h>

#import "BITHockeyManager.h"
#import "BITCrashManager.h"
#import "BITCrashManagerPrivate.h"


@implementation BITCrashExceptionApplication

/*
 * Solution for Scenario 2
 *
 * Catch all exceptions that are being logged to the console and forward them to our
 * custom UncaughtExceptionHandler
 */
- (void)reportException:(NSException *)exception {
  [super reportException: exception];
  
  // Don't invoke the registered UncaughtExceptionHandler if we are currently debugging this app!
  if (![[BITHockeyManager sharedHockeyManager].crashManager isDebuggerAttached] && exception) {
    // We forward this exception to PLCrashReporters UncaughtExceptionHandler
    // If the developer has implemented their own exception handler and that one is
    // invoked before PLCrashReporters exception handler and the developers
    // exception handler is invoking this method it will not finish it's tasks after this
    // call but directly jump into PLCrashReporters exception handler.
    // If we wouldn't do this, this call would lead to an infinite loop.
    
    NSUncaughtExceptionHandler *plcrExceptionHandler = [[BITHockeyManager sharedHockeyManager].crashManager plcrExceptionHandler];
    if (plcrExceptionHandler) {
      plcrExceptionHandler(exception);
    }
  }
}

/*
 * Solution for Scenario 3
 *
 * Exceptions that happen inside an IBAction implementation do not trigger a call to
 * [NSApp reportException:] and it does not trigger a registered UncaughtExceptionHandler
 * Hence we need to catch these ourselves, e.g. by overwriting sendEvent: as done right here
 *
 * On 64bit systems the @try @catch block doesn't even cost any performance.
 */
- (void)sendEvent:(NSEvent *)theEvent {
  @try {
    [super sendEvent:theEvent];
  } @catch (NSException *exception) {
    [self reportException:exception];
  }
}

@end
