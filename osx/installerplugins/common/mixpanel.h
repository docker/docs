//
//  mixpanel.h
//  installerplugins
//
//  Created by Jeffrey Dean Morgan on 8/19/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#ifndef installerplugins_mixpanel_h
#define installerplugins_mixpanel_h

#import <InstallerPlugins/InstallerPlugins.h>

@interface Mixpanel : NSObject

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane;

@end

@implementation Mixpanel

+ (NSString *) uuid {
    NSString *cachePath = [NSSearchPathForDirectoriesInDomains(NSCachesDirectory, NSUserDomainMask, YES) objectAtIndex:0];
    NSString *cacheDirPath = [NSString pathWithComponents:[NSArray arrayWithObjects:cachePath, @"io.docker.pkg.toolbox", nil]];
    NSString *cacheFilePath = [NSString pathWithComponents:[NSArray arrayWithObjects:cacheDirPath, @"id", nil]];
    
    NSString *uuid = [NSString stringWithContentsOfFile:cacheFilePath encoding:NSUTF8StringEncoding error:nil];
    if (!uuid || ![uuid length]) {
        uuid = [[NSUUID UUID] UUIDString];
    }
    
    [[NSFileManager defaultManager] createDirectoryAtPath:cacheDirPath withIntermediateDirectories:YES attributes:nil error:nil];
    [uuid writeToFile:cacheFilePath atomically:YES encoding:NSUTF8StringEncoding error:nil];
    
    return uuid;
}

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane {
    BOOL trackingDisabled = [[[[pane section] sharedDictionary] objectForKey:@"disableTracking"] boolValue];
    NSString *uuid = [self uuid];
    
    if (!uuid || trackingDisabled) {
        return;
    }

    NSBundle* bundle = [[pane section] bundle];
    NSString* token = [bundle objectForInfoDictionaryKey:@"Mixpanel Token"];
    NSString* installerVersion = [bundle objectForInfoDictionaryKey:@"Installer Version"];
    NSString* payload = [NSString stringWithFormat:@"{\"event\": \"%@\", \"properties\": {\"token\": \"%@\", \"distinct_id\": \"%@\", \"os\": \"darwin\", \"version\": \"%@\"}}", name, token, uuid, installerVersion];
    
    NSData * data = [payload dataUsingEncoding:NSUTF8StringEncoding];
    NSString* base64Encoded = [data base64EncodedStringWithOptions:0];
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString: [NSString stringWithFormat:@"https://api.mixpanel.com/track/?data=%@", base64Encoded]]];
    
    NSOperationQueue *queue = [[NSOperationQueue alloc] init];
    [NSURLConnection sendAsynchronousRequest:request queue:queue completionHandler:^(NSURLResponse *response, NSData *data, NSError *error) {}];
}

@end
#endif
