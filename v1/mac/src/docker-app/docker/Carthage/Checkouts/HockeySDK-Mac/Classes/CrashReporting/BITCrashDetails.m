#import "BITCrashDetails.h"
#import "BITCrashDetailsPrivate.h"

@implementation BITCrashDetails

- (instancetype)initWithIncidentIdentifier:(NSString *)incidentIdentifier
                               reporterKey:(NSString *)reporterKey
                                    signal:(NSString *)signal
                             exceptionName:(NSString *)exceptionName
                           exceptionReason:(NSString *)exceptionReason
                              appStartTime:(NSDate *)appStartTime
                                 crashTime:(NSDate *)crashTime
                                 osVersion:(NSString *)osVersion
                                   osBuild:(NSString *)osBuild
                                appVersion:(NSString *)appVersion
                                  appBuild:(NSString *)appBuild
{
  if ((self = [super init])) {
    _incidentIdentifier = incidentIdentifier;
    _reporterKey = reporterKey;
    _signal = signal;
    _exceptionName = exceptionName;
    _exceptionReason = exceptionReason;
    _appStartTime = appStartTime;
    _crashTime = crashTime;
    _osVersion = osVersion;
    _osBuild = osBuild;
    _appVersion = appVersion;
    _appBuild = appBuild;
  }
  return self;
}

@end
