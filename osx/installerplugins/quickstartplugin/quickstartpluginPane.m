//
//  quickstartpluginPane.m
//  quickstartplugin
//
//  Created by Jeffrey Dean Morgan on 8/19/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#import "quickstartpluginPane.h"
#import "mixpanel.h"

@implementation quickstartpluginPane

- (NSString *)title {
	return [[NSBundle bundleForClass:[self class]] localizedStringForKey:@"PaneTitle" value:nil table:nil];
}

- (void) willEnterPane:(InstallerSectionDirection)dir {
    self.previousEnabled = NO;
    
    BOOL kitematicInstalled = [[NSFileManager defaultManager] fileExistsAtPath:@"/Applications/Docker/Kitematic (Beta).app"];
    BOOL quickstartInstalled = [[NSFileManager defaultManager] fileExistsAtPath:@"/Applications/Docker/Docker Quickstart Terminal.app"];
    
    self.kitematicImageView.image = [[NSImage alloc] initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"kitematic" ofType:@"png"]];
    self.quickstartImageView.image = [[NSImage alloc] initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"quickstart" ofType:@"png"]];
    
    if (!kitematicInstalled) {
        self.kitematicImageView.enabled = NO;
        self.kitematicLabel.enabled = NO;
    }

    if (!quickstartInstalled) {
        self.quickstartImageView.enabled = NO;
        self.quickstartLabel.enabled = NO;
    }

    if (dir == InstallerDirectionForward && !kitematicInstalled && !quickstartInstalled) {
        [self gotoNextPane];
    }
}

- (void) willExitPane:(InstallerSectionDirection)dir {
    if (dir == InstallerDirectionForward) {
        [Mixpanel trackEvent:@"Installer Finished" forPane:self withProperties:[[NSDictionary alloc] initWithObjectsAndKeys:@"Continue Button", @"action", nil]];
    }
}

- (IBAction)quickstartTerminalClicked:(id)sender {
    [Mixpanel trackEvent:@"Installer Finished" forPane:self withProperties:[[NSDictionary alloc] initWithObjectsAndKeys:@"Quickstart Terminal", @"action", nil]];
    NSTask *task = [[NSTask alloc] init];
    task.launchPath = @"/usr/bin/open";
    task.arguments = @[@"/Applications/Docker/Docker Quickstart Terminal.app"];
    [task launch];
}

- (IBAction)kitematicClicked:(id)sender {
    [Mixpanel trackEvent:@"Installer Finished" forPane:self withProperties:[[NSDictionary alloc] initWithObjectsAndKeys:@"Kitematic", @"action", nil]];
    [[NSWorkspace sharedWorkspace] launchApplication:@"Kitematic (Beta)"];
}

@end
