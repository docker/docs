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
@property BOOL upgrading;
@property BOOL locked;
@property BOOL successful;
@end

@implementation managevmpluginPane

NSString *vBoxManagePath = @"/Applications/VirtualBox.app/Contents/MacOS/VBoxManage";
NSString *dockerMachinePath = @"/usr/local/bin/docker-machine";

- (BOOL) vmExists:(NSString*)name {
    NSTask* task = [NSTask launchedTaskWithLaunchPath:@"/usr/bin/sudo" arguments:[NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), vBoxManagePath, @"showvminfo", name, nil]];
    [task waitUntilExit];
    return [task terminationStatus] != 1;
}

- (NSString *) boot2dockerISOVersionAtPath:(NSString*)path {
    NSTask* task = [[NSTask alloc] init];
    task.arguments =  [NSArray arrayWithObjects:path, nil];
    task.launchPath = @"/usr/bin/file";

    NSPipe * out = [NSPipe pipe];
    [task setStandardOutput:out];
    [task launch];
    [task waitUntilExit];
    
    if (task.terminationStatus != 0) {
        return nil;
    }

    NSFileHandle * read = [out fileHandleForReading];
    NSData * dataRead = [read readDataToEndOfFile];
    NSString * stringRead = [[NSString alloc] initWithData:dataRead encoding:NSUTF8StringEncoding];
    
    NSCharacterSet *delimiters = [NSCharacterSet characterSetWithCharactersInString:@"v'"];
    NSArray *splitString = [stringRead componentsSeparatedByCharactersInSet:delimiters];

    if (splitString.count < 3) {
        return nil;
    }
    
    NSString *version = [[splitString objectAtIndex:2] stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    return version;
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

- (BOOL) canUpgradeBoot2DockerVM {
    // VirtualBox and Docker Machine exist
    if (![[NSFileManager defaultManager] fileExistsAtPath:vBoxManagePath] || ![[NSFileManager defaultManager] fileExistsAtPath:dockerMachinePath]) {
        return NO;
    }
    
    // docker-machine path and VirtualBox VM exist
    if (![self vmExists:@"default"] ||
        ![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.docker/machine/machines/default", NSHomeDirectory()]] ||
        ![[NSFileManager defaultManager] fileExistsAtPath:[NSString stringWithFormat:@"%@/.docker/machine/machines/default/boot2docker.iso", NSHomeDirectory()]]) {
        return NO;
    }
    
    NSString *incomingVersion = [self boot2dockerISOVersionAtPath:[NSString stringWithFormat:@"%@/.docker/machine/cache/boot2docker.iso", NSHomeDirectory()]];
    NSString *existingVersion = [self boot2dockerISOVersionAtPath:[NSString stringWithFormat:@"%@/.docker/machine/machines/default/boot2docker.iso", NSHomeDirectory()]];
    
    return incomingVersion != existingVersion && [incomingVersion compare:existingVersion options:NSNumericSearch] != NSOrderedAscending;
}

- (void) migrateBoot2DockerVM {
    self.locked = YES;
    
    // Remove existing vm if it exists
    NSTask* removeVMTask = [NSTask launchedTaskWithLaunchPath:@"/usr/bin/sudo" arguments:[NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), dockerMachinePath, @"rm", @"-f", @"default", nil]];
    [removeVMTask waitUntilExit];
    
    // Remove the VM dir in case there's anything left over
    NSTask* removeDirTask = [NSTask launchedTaskWithLaunchPath:@"/bin/rm" arguments:[NSArray arrayWithObjects:@"-rf", [NSString stringWithFormat:@"%@/.docker/machine/machines/default", NSHomeDirectory()], nil]];
    [removeDirTask waitUntilExit];
    
    // Do the migration
    NSTask* migrateTask = [[NSTask alloc] init];
    migrateTask.launchPath = @"/usr/bin/sudo";
    migrateTask.arguments = [NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), dockerMachinePath, @"-D", @"create", @"-d", @"virtualbox", @"--virtualbox-memory", @"2048", @"--virtualbox-import-boot2docker-vm", @"boot2docker-vm", @"default", nil];
    
    
    [self runTask:migrateTask onFinish:^void (int status) {
        self.nextEnabled = YES;
        self.statusImage.hidden = YES;
        [self.migrationProgress stopAnimation:self];
        self.migrationProgress.hidden = YES;
        
        self.migrationStatusLabel.hidden = NO;
        self.statusImage.hidden = NO;
        
        if (status == 0) {
            self.successful = YES;
            self.migrationStatusLabel.stringValue = @"Your Boot2Docker VM was successfully migrated to a Docker Machine VM named \"default\".";
            self.statusImage.image = [[NSImage alloc]initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"toolboxcheck" ofType:@"png"]];
            [Mixpanel trackEvent:@"Boot2Docker Migration Succeeded" forPane:self];
            [self gotoNextPane];
        } else {
            self.migrationLogsScrollView.hidden = NO;
            self.submitButton.hidden = NO;
            [Mixpanel trackEvent:@"Boot2Docker Migration Failed" forPane:self];
            self.migrationStatusLabel.hidden = NO;
            self.migrationStatusLabel.stringValue = @"Creating the VM failed. Following the install, try creating a vm manually via docker-machine.";
        }
 
    }];
}

- (void) upgradeBoot2DockerVM {
    self.locked = YES;
    
    // Do the migration
    NSTask* task = [[NSTask alloc] init];
    task.launchPath = @"/usr/bin/sudo";
    task.arguments = [NSArray arrayWithObjects:@"-i", @"-u", NSUserName(), [NSString stringWithFormat:@"%@ stop default && cp %@/.docker/machine/cache/boot2docker.iso %@/.docker/machine/machines/default/boot2docker.iso", dockerMachinePath, NSHomeDirectory(),  NSHomeDirectory()], nil];

    [self runTask:task onFinish:^void (int status) {
        self.nextEnabled = YES;
        self.statusImage.hidden = YES;
        [self.migrationProgress stopAnimation:self];
        self.migrationProgress.hidden = YES;
        
        self.migrationStatusLabel.hidden = NO;
        self.statusImage.hidden = NO;
        
        if (status == 0) {
            self.successful = YES;
            self.migrationStatusLabel.stringValue = @"Your VirtualBox Docker VM named \"default\" was successfully upgraded.";
            self.statusImage.image = [[NSImage alloc]initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"toolboxcheck" ofType:@"png"]];
            [Mixpanel trackEvent:@"VM Upgrade Succeeded" forPane:self];
            [self gotoNextPane];
        } else {
            self.migrationLogsScrollView.hidden = NO;
            self.submitButton.hidden = NO;
            [Mixpanel trackEvent:@"VM Upgrade Failed" forPane:self];
            self.migrationStatusLabel.hidden = NO;
            self.migrationStatusLabel.stringValue = @"Upgrading your VirtualBox Docker VM Failed. Try upgrading again manually via the docker-machine command.";
        }
        
    }];
}

- (void) runTask:(NSTask *)task onFinish:(void (^)(int))finish {
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
    
    task.standardOutput = [NSPipe pipe];
    task.standardError = [NSPipe pipe];
    
    [[task.standardOutput fileHandleForReading] setReadabilityHandler:appendOutput];
    [[task.standardError fileHandleForReading] setReadabilityHandler:appendOutput];
    
    task.terminationHandler = ^(NSTask* task) {
        dispatch_async(dispatch_get_main_queue(), ^{
            [task.standardOutput fileHandleForReading].readabilityHandler = nil;
            finish(task.terminationStatus);
        });
    };
    
    [task launch];
}

- (id) init {
    self.upgrading = NO;
    self.successful = NO;
    self = [super init];
    return self;
}

- (NSString *)title {
    return [[NSBundle bundleForClass:[self class]] localizedStringForKey:@"PaneTitle" value:nil table:nil];
}

- (void) didEnterPane:(InstallerSectionDirection)dir {
    [Mixpanel trackEvent:@"Installing Files Succeeded" forPane:self];
    self.previousEnabled = NO;
    
    NSImage *toolboxImage = [[NSImage alloc]initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"toolbox" ofType:@"png"]];
    NSImage *boot2dockerImage = [[NSImage alloc] initWithContentsOfFile:[[NSBundle bundleForClass:[self class]] pathForResource:@"boot2docker" ofType:@"png"]];

    [self.migrationProgress startAnimation:self];

    if ([self canUpgradeBoot2DockerVM]) {
        self.migrateCheckbox.title = [NSString stringWithFormat:@"Upgrade your Docker Toolbox VM to Docker %@", [self boot2dockerISOVersionAtPath:[NSString stringWithFormat:@"%@/.docker/machine/cache/boot2docker.iso", NSHomeDirectory()]]] ;
        self.migrateExtraLabel.stringValue = @"Your existing Docker Toolbox VM will not be affected. This should take about a minute.";
        self.statusImage.image = toolboxImage;
        self.upgrading = YES;
    } else if ([self canMigrateBoot2DockerVM]) {
        self.statusImage.image = boot2dockerImage;
    } else {
        [self gotoNextPane];
    }
    
    self.statusImage.hidden = NO;
    self.migrateCheckbox.hidden = NO;
    self.migrateCheckbox.enabled = YES;
    self.migrateExtraLabel.hidden = NO;
    self.migrationProgress.hidden = YES;
    [self.migrationProgress stopAnimation:self];
}

- (BOOL) shouldExitPane:(InstallerSectionDirection)dir {
    if (dir != InstallerDirectionForward) {
        return YES;
    }
    
    if (self.migrateCheckbox.enabled && self.migrateCheckbox.state == NSOnState) {
        self.nextEnabled = false;
        self.migrationProgress.hidden = NO;
        self.migrationStatusLabel.hidden = NO;
        self.migrateCheckbox.enabled = NO;
        [self.migrationProgress startAnimation:self];

        if (self.upgrading) {
            [Mixpanel trackEvent:@"VM Upgrade Started" forPane:self];
            self.migrationStatusLabel.stringValue = @"Upgrading...";
            [self upgradeBoot2DockerVM];
        } else {
            [Mixpanel trackEvent:@"Boot2Docker Migration Started" forPane:self];
            self.migrationStatusLabel.stringValue = @"Migrating...";
            [self migrateBoot2DockerVM];
        }

        return NO;
    } else if (self.migrateCheckbox.state == NSOffState) {
        if (self.upgrading) {
            [Mixpanel trackEvent:@"VM Upgrade Skipped" forPane:self];
        } else {
            [Mixpanel trackEvent:@"Boot2Docker Migration Skipped" forPane:self];
        }
        
    }

    return YES;
}

- (IBAction)submitButtonClicked:(id)sender {
    NSURL *url = [[NSURL alloc] initWithString: @"https://github.com/docker/toolbox/issues"];
    [[NSWorkspace sharedWorkspace] openURL:url];
}

@end
