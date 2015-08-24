//
//  overviewpluginPane.h
//  overviewplugin
//
//  Created by Jeffrey Dean Morgan on 8/24/15.
//  Copyright (c) 2015 Docker Inc. All rights reserved.
//

#import <InstallerPlugins/InstallerPlugins.h>

@interface overviewpluginPane : InstallerPane
@property (unsafe_unretained) IBOutlet NSTextView *textView;
@property (weak) IBOutlet NSButton *checkbox;
@end
