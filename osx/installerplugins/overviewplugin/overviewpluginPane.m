//
//  overviewpluginPane.m
//  overviewplugin
//
//  Created by Jeffrey Dean Morgan on 8/24/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#import "overviewpluginPane.h"
#import "mixpanel.h"

@implementation overviewpluginPane

- (NSString *)title {
    return [[NSBundle bundleForClass:[self class]] localizedStringForKey:@"PaneTitle" value:nil table:nil];
}

- (void) willEnterPane:(InstallerSectionDirection)dir {
    NSURL *url = [[[self section] bundle] URLForResource:@"overview" withExtension:@"html"];
    NSMutableAttributedString *formattedHTML = [[NSMutableAttributedString alloc] initWithURL:url documentAttributes:nil];
    [formattedHTML addAttribute:NSFontAttributeName value:[NSFont systemFontOfSize:13.0f] range:NSMakeRange(0, formattedHTML.length)];

    [[self.textView textStorage] setAttributedString:formattedHTML];
    
    [Mixpanel trackEvent:@"Installer Started" forPane:self];
}

- (void) willExitPane:(InstallerSectionDirection)dir {
    if (dir != InstallerDirectionForward) {
        return;
    }
    
    if (self.checkbox.state != NSOnState) {
        [Mixpanel trackEvent:@"Disabled Tracking" forPane:self];
        [[[self section] sharedDictionary] setObject:[NSNumber numberWithBool:YES] forKey:@"disableTracking"];
        return;
    }
    
    [[[self section] sharedDictionary] removeObjectForKey:@"disableTracking"];
    [Mixpanel trackEvent:@"Enabled Tracking" forPane:self];
}

@end
