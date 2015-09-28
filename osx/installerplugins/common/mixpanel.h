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
#import "utils.h"

@interface Mixpanel : NSObject

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane;

@end

@implementation Mixpanel

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane withProperties:(NSDictionary *)properties {
    BOOL trackingDisabled = [[[[pane section] sharedDictionary] objectForKey:@"disableTracking"] boolValue];
    if (trackingDisabled) {
        return;
    }

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
    NSString* token = [bundle objectForInfoDictionaryKey:@"Mixpanel Token"];
    NSString* installerVersion = [bundle objectForInfoDictionaryKey:@"Installer Version"];
    NSString* payload = [NSString stringWithFormat:@"{\"event\": \"%@\", \"properties\": {\"token\": \"%@\", \"distinct_id\": \"%@\", \"os\": \"darwin\", \"os version\":\"%@\", \"version\": \"%@\" %@}}", name, token, uuid, osVersion, installerVersion, props];

    @try {
        NSData * data = [payload dataUsingEncoding:NSUTF8StringEncoding];
        NSString* base64Encoded = [data base64EncodedStringWithOptions:0];
        NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString: [NSString stringWithFormat:@"https://api.mixpanel.com/track/?data=%@", base64Encoded]]];
        NSOperationQueue *queue = [[NSOperationQueue alloc] init];
        [NSURLConnection sendAsynchronousRequest:request queue:queue completionHandler:^(NSURLResponse *response, NSData *data, NSError *error) {}];
    }
    @catch (NSException *exception) {
        NSLog(@"%@", @"Failed to send data.");
    }
}

+ (void) trackEvent:(NSString *)name forPane:(InstallerPane*)pane {
    [self trackEvent:name forPane:pane withProperties:nil];
}

@end
#endif
