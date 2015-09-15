//
//  uuid.h
//  installerplugins
//
//  Created by Jeffrey Dean Morgan on 9/14/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#ifndef installerplugins_uuid_h
#define installerplugins_uuid_h

@interface Utils : NSObject

+ (NSString *) uuid;

@end

@implementation Utils

+ (NSString *) uuid {
    NSString *appPath = [NSSearchPathForDirectoriesInDomains(NSApplicationSupportDirectory, NSUserDomainMask, YES) objectAtIndex:0];
    NSString *appDirPath = [NSString pathWithComponents:[NSArray arrayWithObjects:appPath, @"DockerToolbox", nil]];
    NSString *appFilePath = [NSString pathWithComponents:[NSArray arrayWithObjects:appDirPath, @"id", nil]];
    
    NSString *uuid = [NSString stringWithContentsOfFile:appFilePath encoding:NSUTF8StringEncoding error:nil];
    if (!uuid || ![uuid length]) {
        uuid = [[NSUUID UUID] UUIDString];
        [[NSFileManager defaultManager] createDirectoryAtPath:appDirPath withIntermediateDirectories:YES attributes:nil error:nil];
        [uuid writeToFile:appFilePath atomically:YES encoding:NSUTF8StringEncoding error:nil];
    }
    
    return uuid;
}

@end
#endif
