//
//  MultipleUsers.swift
//  docker
//
//  Created by Doby Mock on 5/13/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

// MultipleUsers provides functions to make sure several users can use
// Docker on the same Mac.
class MultipleUsers: NSObject {
    
    static func listenForUserSessionSwitch() {
        NSWorkspace.sharedWorkspace().notificationCenter.addObserver(MultipleUsers.self, selector: #selector(MultipleUsers.userSessionDidBecomeActive), name: NSWorkspaceSessionDidBecomeActiveNotification, object: nil)
    }
    
    static func userSessionDidBecomeActive() {
        Logger.log(level: ASLLogger.Level.Notice, content: "userSessionDidBecomeActive -> update symlinks")
        Install.installSymlinks()
    }
    
}