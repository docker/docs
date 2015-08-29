//
//  managevmpluginPane.m
//  managevmplugin
//
//  Created by Jeffrey Dean Morgan on 8/19/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#import "managevmpluginPane.h"
#import "mixpanel.h"

@interface managevmpluginPane()
@property BOOL migrating;
@property BOOL successfullyMigrated;
@end

@implementation managevmpluginPane

NSString *vBoxManagePath = @"/Applications/VirtualBox.app/Contents/MacOS/VBoxManage";
NSString *dockerMachinePath = @"/usr/local/bin/docker-machine";

- (BOOL) vmExists:(NSString*)name {
    NSTask* task = [NSTask launchedTaskWithLaunchPath:@"/usr/bin/sudo" arguments:[NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), vBoxManagePath, @"showvminfo", name, nil]];
    [task waitUntilExit];
    return [task terminationStatus] != 1;
}

- (BOOL) canMigrateBoot2DockerVM {
    // VirtualBox and Docker Machine exist
    if (![[NSFileManager defaultManager] fileExistsAtPath:vBoxManagePath] || ![[NSFileManager defaultManager] fileExistsAtPath:dockerMachinePath]) {
        return NO;
    }
    
    // Boot2Docker certs exist
    if (![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.boot2docker/certs/boot2docker-vm/ca.pem", NSHomeDirectory()]] ||
        ![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.boot2docker/certs/boot2docker-vm/cert.pem", NSHomeDirectory()]] ||
        ![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.boot2docker/certs/boot2docker-vm/key.pem", NSHomeDirectory()]]) {
        return NO;
    }
    
    // Boot2Docker ssh keys exist
    if (![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.ssh/id_boot2docker", NSHomeDirectory()]] ||
        ![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.ssh/id_boot2docker.pub", NSHomeDirectory()]]) {
        return NO;
    }
    
    if (![self vmExists:@"boot2docker-vm"] || [self vmExists:@"default"]) {
        return NO;
    }
    return YES;
}

- (void) migrateBoot2DockerVM {
    
    self.migrating = YES;
    
    // Remove existing vm if it exists (obviously user must have deleted the
    NSTask* removeVMTask = [NSTask launchedTaskWithLaunchPath:@"/usr/bin/sudo" arguments:[NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), dockerMachinePath, @"rm", @"-f", @"default", nil]];
    [removeVMTask waitUntilExit];
    
    // Remove the VM dir in case there's anything left over
    NSTask* removeDirTask = [NSTask launchedTaskWithLaunchPath:@"/bin/rm" arguments:[NSArray arrayWithObjects:@"-rf", [NSString stringWithFormat:@"%@/.docker/machine/machines/default", NSHomeDirectory()], nil]];
    [removeDirTask waitUntilExit];
    
    // Do the migration
    NSTask* migrateTask = [[NSTask alloc] init];
    migrateTask.launchPath = @"/usr/bin/sudo";
    migrateTask.arguments = [NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), dockerMachinePath, @"-D", @"create", @"-d", @"virtualbox", @"--virtualbox-memory", @"2048", @"--virtualbox-import-boot2docker-vm", @"boot2docker-vm", @"default", nil];
    
    // Remove certificates, ssh keys from logs
    NSRegularExpression *regex = [NSRegularExpression regularExpressionWithPattern:@"BEGIN.*END" options:NSRegularExpressionDotMatchesLineSeparators error:NULL];
    NSFont *font = [NSFont fontWithName:@"Menlo" size:10.0];
    NSDictionary *attrsDictionary = [NSDictionary dictionaryWithObject:font forKey:NSFontAttributeName];
    NSMutableData* fullData = [[NSMutableData alloc] init];
    
    void (^appendOutput)(NSFileHandle*) = ^(NSFileHandle *file) {
        NSData *data = [file availableData];
        [fullData appendData:data];
        NSMutableString *str = [[NSMutableString alloc] initWithData:fullData encoding:NSUTF8StringEncoding];
        [regex replaceMatchesInString:str options:0 range:NSMakeRange(0, [str length]) withTemplate:@""];
        
        dispatch_async(dispatch_get_main_queue(), ^{
            [self.migrationLogsTextView.textStorage setAttributedString:[[NSAttributedString alloc] initWithString:str attributes:attrsDictionary]];
            [self.migrationLogsTextView scrollRangeToVisible:NSMakeRange([[self.migrationLogsTextView string] length], 0)];
        });
    };

    migrateTask.standardOutput = [NSPipe pipe];
    migrateTask.standardError = [NSPipe pipe];
    
    [[migrateTask.standardOutput fileHandleForReading] setReadabilityHandler:appendOutput];
    [[migrateTask.standardError fileHandleForReading] setReadabilityHandler:appendOutput];
    
    migrateTask.terminationHandler = ^(NSTask* task) {
        dispatch_async(dispatch_get_main_queue(), ^{
            [task.standardOutput fileHandleForReading].readabilityHandler = nil;
            
            self.boot2dockerImage.hidden = YES;
            self.toolboxImage.hidden = YES;
            [self.migrationProgress stopAnimation:self];
            self.migrationProgress.hidden = YES;
            
            self.migrationStatusLabel.hidden = NO;
            self.migrationStatusImage.hidden = NO;
            
            if (task.terminationStatus == 0) {
                self.successfullyMigrated = YES;
                self.migrationStatusLabel.stringValue = @"Your Boot2Docker VM was successfully migrated to a Docker Machine VM named \"default\".";
                self.migrationStatusImage.image = [[NSImage alloc]initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"toolboxcheck" ofType:@"png"]];
                [Mixpanel trackEvent:@"Boot2Docker Migration Succeeded" forPane:self];
                self.nextEnabled = YES;
                [self gotoNextPane];
            } else {
                [Mixpanel trackEvent:@"Boot2Docker Migration Failed" forPane:self];
                self.migrationStatusLabel.hidden = NO;
                self.migrationStatusLabel.stringValue = @"VM Migration failed. Please see the logs below.";
            }
        });
    };
    
    [migrateTask launch];
}

- (id) init {
    self.migrating = NO;
    self.successfullyMigrated = NO;
    self = [super init];
    return self;
}

- (NSString *)title {
    return [[NSBundle bundleForClass:[self class]] localizedStringForKey:@"PaneTitle" value:nil table:nil];
}

- (void) didEnterPane:(InstallerSectionDirection)dir {
    [Mixpanel trackEvent:@"Installing Files Succeeded" forPane:self];
    if (self.successfullyMigrated) {
        self.nextEnabled = YES;
        return;
    }
    
    if (self.migrating) {
        self.nextEnabled = NO;
        return;
    }
    
    [self.migrationProgress startAnimation:self];
    
    if ([self canMigrateBoot2DockerVM]) {
        NSImage *image = [[NSImage alloc] initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"boot2docker" ofType:@"png"]];
        self.boot2dockerImage.image = image;
        self.boot2dockerImage.hidden = NO;
        NSImage *toolboxImage = [[NSImage alloc]initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"toolbox" ofType:@"png"]];
        self.toolboxImage.image = toolboxImage;
        self.toolboxImage.hidden = NO;
        self.arrowImage.hidden = NO;
        self.migrateCheckbox.hidden = NO;
        self.migrateCheckbox.enabled = YES;
        self.migrateExtraLabel.hidden = NO;
        [self.migrationProgress stopAnimation:self];
        self.migrationProgress.hidden = YES;
    } else {
        [self gotoNextPane];
    }
}

- (BOOL) shouldExitPane:(InstallerSectionDirection)dir {
    if (dir != InstallerDirectionForward) {
        return YES;
    }
    
    if (self.migrateCheckbox.enabled && self.migrateCheckbox.state == NSOnState) {
        [Mixpanel trackEvent:@"Boot2Docker Migration Started" forPane:self];

        self.nextEnabled = false;
        self.migrationProgress.hidden = NO;
        [self.migrationProgress startAnimation:self];
        self.migrationStatusLabel.hidden = NO;
        self.arrowImage.hidden = YES;
        self.migrateCheckbox.enabled = NO;
        self.migrationStatusLabel.stringValue = @"Migrating...";
        self.migrationLogsScrollView.hidden = NO;

        [self migrateBoot2DockerVM];
        return NO;
    } else if (self.migrateCheckbox.state == NSOffState) {
        [Mixpanel trackEvent:@"Boot2Docker Migration Skipped" forPane:self];
    }
    return YES;
}


@end
