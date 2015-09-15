//
//  bugsnag.h
//  installerplugins
//
//  Created by Jeffrey Dean Morgan on 9/12/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#ifndef installerplugins_bugsnag_h
#define installerplugins_bugsnag_h

#import <InstallerPlugins/InstallerPlugins.h>
#import "utils.h"

@interface Bugsnag : NSObject

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane;

@end

@implementation Bugsnag

+ (void) reportError:(NSString *)name forPane:(InstallerPane*)pane withProperties:(NSDictionary *)properties {
    NSString *uuid = [Utils uuid];
    if (!uuid) {
        return;
    }

    NSString *props = @"";
    for (NSString *key in properties) {
        props = [props stringByAppendingFormat:@",\"%@\": \"%@\"", key, [properties objectForKey:key]];
    }

    NSString *osVersion = [@"Mac OS X " stringByAppendingString:[[[[NSProcessInfo processInfo] operatingSystemVersionString] componentsSeparatedByString:@" "] objectAtIndex:1]];

    NSBundle* bundle = [[pane section] bundle];
    NSString* token = [bundle objectForInfoDictionaryKey:@"Bugsnag Token"];
    NSString* installerVersion = [bundle objectForInfoDictionaryKey:@"Installer Version"];
    NSString* payload = [NSString stringWithFormat:@"{\"event\": \"%@\", \"properties\": {\"token\": \"%@\", \"distinct_id\": \"%@\", \"os\": \"darwin\", \"os version\":\"%@\", \"version\": \"%@\" %@}}", name, token, uuid, osVersion, installerVersion, props];

    NSData * data = [payload dataUsingEncoding:NSUTF8StringEncoding];
    NSString* base64Encoded = [data base64EncodedStringWithOptions:0];
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString: [NSString stringWithFormat:@"https://notify.bugsnag.com", base64Encoded]]];

    NSOperationQueue *queue = [[NSOperationQueue alloc] init];
    [NSURLConnection sendAsynchronousRequest:request queue:queue completionHandler:^(NSURLResponse *response, NSData *data, NSError *error) {}];
}

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane {
    [self trackEvent:name forPane:pane withProperties:nil];
}

@end
#endif
