#import <Foundation/Foundation.h>

#import "PLCrashNamespace.h"
#import "PLCrashReport.h"

// Dictionary keys for array elements returned by arrayOfAppUUIDsForCrashReport:
#ifndef kBITBinaryImageKeyUUID
#define kBITBinaryImageKeyUUID @"uuid"
#define kBITBinaryImageKeyArch @"arch"
#define kBITBinaryImageKeyType @"type"
#endif


@interface BITCrashReportTextFormatter : NSObject {
}

+ (NSString *)stringValueForCrashReport:(BITPLCrashReport *)report crashReporterKey:(NSString *)crashReporterKey;
+ (NSArray *)arrayOfAppUUIDsForCrashReport:(BITPLCrashReport *)report;

@end
