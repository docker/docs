//
//  improvetoolboxpluginPane.m
//  improvetoolboxplugin
//
//  Created by Jeffrey Dean Morgan on 8/19/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#import "improvetoolboxpluginPane.h"
#import "mixpanel.h"

@implementation improvetoolboxpluginPane

- (NSString *)title {
    return [[NSBundle bundleForClass:[self class]] localizedStringForKey:@"PaneTitle" value:nil table:nil];
}

- (void) willExitPane:(InstallerSectionDirection)dir {
    if (dir != InstallerDirectionForward) {
        return;
    }
    
    if (self.checkbox.state != NSOnState) {
        [[[self section] sharedDictionary] removeObjectForKey:@"uuid"];
        return;
    }
    
    NSString *uuid = [[[self section] sharedDictionary] objectForKey:@"uuid"];
    if (!uuid) {
        uuid = [[NSUUID UUID] UUIDString];
        [[[self section] sharedDictionary] setObject:uuid forKey:@"uuid"];
    }
    
    [Mixpanel trackEvent:@"Enabled Tracking" forPane:self];
}

@end
