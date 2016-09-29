//
//  Paths.swift
//  docker
//
//  Created by Emmanuel Briney on 28/12/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Foundation

// Paths class provides functions that return paths for the application storage
// locations (as a convenience)
class Paths {
    
    // groupContainerPath returns the path to the app-group container of Config.dockerAppGroupIdentifier.
    // returns:
    // - path to the app-group container or nil if failed
    static func groupContainerPath() -> String? {

        // get URL to the app-group container
        // In the context of a sandboxed application, this will create the 
        // group container directory (doesn't create it if the app is not sandboxed)
        if let groupContainerPath = NSFileManager.defaultManager().containerURLForSecurityApplicationGroupIdentifier(Config.dockerAppGroupIdentifier)?.path {
        
            ////////////////////////////////////////////////////////////////////////
            ////////////////////////////////////////////////////////////////////////
            /// TODO: BACK TO THE SANDBOX (gdevillele)
            /// this is a hack that is needed only because Docker.app is not a sandboxed app
            /// this block must be removed when the app is back in the sandbox
            ////////////////////////////////////////////////////////////////////////
            ////////////////////////////////////////////////////////////////////////
            // Apple sandbox issue :
            // The documentation (https://developer.apple.com/library/mac/documentation/Cocoa/Reference/Foundation/Classes/NSFileManager_Class/#//apple_ref/occ/instm/NSFileManager/containerURLForSecurityApplicationGroupIdentifier:)
            // says that the directory will be created if not exists but i did some
            // tests and the directory is not created if the sandbox flag is not set!
            // If you create it with code the preferences folder will never been
            // created and NSUserDefault.sync will not save your properties in the
            // Group Containers plist.
            // As a workaround, writing a dummy preference in the shared preferences
            // will trigger the creation of the Group Container directory.
            // So we create a dummy pref and then remove it.
            //
            // test if the shared preferences plist already exists

            if let ud: NSUserDefaults = NSUserDefaults(suiteName: Config.dockerAppGroupIdentifier) {
                ud.setValue(true, forKey: "Initialized")
            }

            let preferencesPath: String = NSString.pathWithComponents([groupContainerPath, "Library", "Preferences"])
            CreateDirectoryIfNecessary(preferencesPath)
            let plistFullPath = NSString.pathWithComponents([preferencesPath, Config.dockerAppGroupIdentifier + ".plist"])
            if NSFileManager.defaultManager().fileExistsAtPath(plistFullPath) == false {
                // if it doesn't exist, we store a dummy pref value in it
                let userDefaults = NSUserDefaults.init(suiteName: Config.dockerAppGroupIdentifier)
                userDefaults?.setValue(true, forKey: "Initialized")
                userDefaults?.synchronize()
            }
            return groupContainerPath
        } else {
            // return an error as we failed to get the container URL for the group
            Logger.log(level: ASLLogger.Level.Fatal, content: "failed to get \(Config.dockerAppGroupIdentifier) app-group container")
            return nil
        }
    }


    ////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////
    /// NOTE: (gdevillele)
    /// this is needed only because Docker.app is not a sandboxed app
    /// When it is back in the sandbox, we can replace the call to this function
    /// by "NSHomeDirectory()"
    ////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////
    static func appContainerPath() -> String? {
        let home = NSHomeDirectory()
        if home.containsString("/Library/Containers/") && home.hasSuffix("/Data") {
            // it means the app is sandboxed, we return the home directory value
            return home
        }
        // the app isn't sandboxed, we handcraft the path to the app container
        guard let bundleID: String = NSBundle.mainBundle().bundleIdentifier else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "failed to find app container path")
            return nil
        }
        let directoryPath: String = NSString.pathWithComponents([home, "Library", "Containers", bundleID, "Data"])
        // try to create the directory if it doesn't exit.
        if let error: String = CreateDirectoryIfNecessary(directoryPath) {
            Logger.log(level: Logger.Level.Fatal, content: "failed to create app container directory (\(error))")
            return nil
        }
        return directoryPath
    }
    ////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////
    
    
    
    // This function will not try to override an existing file
    // returns nil if the directory has been created or if a directory already exists
    // at the given path. Otherwise, it returns an error message.
    static func CreateDirectoryIfNecessary(absolutePath: String) -> String? {
        do {
            var isDir: ObjCBool = ObjCBool(false)
            if NSFileManager.defaultManager().fileExistsAtPath(absolutePath, isDirectory:&isDir) {
                // file exists
                if isDir.boolValue {
                    // file is a directory, we don't need to do anything
                    return nil
                } else {
                    // file exists but it's not a directory
                    return "a regular file already exists at the given path (\(absolutePath))"
                }
            }
            // file doesn't exist. Try to create a directory
            try NSFileManager.defaultManager().createDirectoryAtPath(absolutePath, withIntermediateDirectories: true, attributes: nil)
        } catch let error as NSError {
            return error.localizedDescription
        }
        return nil
    }
    
    // returns whether a file exists and it is a regular file (not a directory)
    static func regularFileExists(absolutePath: String) -> Bool {
        var isDir: ObjCBool = ObjCBool(false)
        if NSFileManager.defaultManager().fileExistsAtPath(absolutePath, isDirectory:&isDir) {
            // a file exists, we test if it is a regular file
            if isDir.boolValue == false {
                // it is not a directory
                return true
            }
        }
        return false
    }
    
    // return whether a directory exists at the given path
    static func directoryExists(absolutePath: String) -> Bool {
        var isDir: ObjCBool = ObjCBool(false)
        if NSFileManager.defaultManager().fileExistsAtPath(absolutePath, isDirectory:&isDir) {
            // a file exists, we test if it is a regular file
            if isDir.boolValue {
                // it is a directory
                return true
            }
        }
        return false
    }
}
