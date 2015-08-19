//
//  quickstartpluginPane.h
//  quickstartplugin
//
//  Created by Jeffrey Dean Morgan on 8/19/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#import <InstallerPlugins/InstallerPlugins.h>

@interface quickstartpluginPane : InstallerPane
@property (weak) IBOutlet NSButton *quickstartImageView;
@property (weak) IBOutlet NSButton *kitematicImageView;
@property (weak) IBOutlet NSTextField *quickstartLabel;
@property (weak) IBOutlet NSTextField *kitematicLabel;

@end
