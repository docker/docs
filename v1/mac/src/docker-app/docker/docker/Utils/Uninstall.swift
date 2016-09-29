//
//  Uninstall.swift
//  docker
//
//  Created by Gaetan de Villele on 3/18/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation

class Uninstall {
    
    // returns nil or an error message
    static func uninstallWhileKeepingPreferenceKeys(prefKeysToKeep: [Preferences.Key]) -> String? {
        
        // stopping backend processes
        ChildProcesses.stop()
        
        // Ask vmnetd to remove it stuff
        
        // uninstall symlinks (backed up binaries will be renamed)
        var success = Install.uninstallSymlinks()
        if !success {
            // NOTE(aduermael): Just display an error and continue, it's better 
            //to uninstall as much as possible so far, even if one step fails.
            Logger.log(level: ASLLogger.Level.Error, content: "Something went wrong when uninstalling symlinks")
        }
        
        success = Install.uninstallSockets()
        if !success {
            // NOTE(aduermael): Just display an error and continue, it's better
            //to uninstall as much as possible so far, even if one step fails.
            Logger.log(level: ASLLogger.Level.Error, content: "Something went wrong when deleting sockets")
        }
        
        // uninstallVmnetd also removes sockets in /var/tmp
        success = Install.uninstallVmnetd()
        if !success {
            // NOTE(aduermael): Just display an error and continue, it's better
            //to uninstall as much as possible so far, even if one step fails.
            Logger.log(level: ASLLogger.Level.Error, content: "Something went wrong when uninstalling networking")
        }
        
        // deleting preferences
        guard let prefs: NSUserDefaults = NSUserDefaults(suiteName: Config.dockerAppGroupIdentifier) else {
            return "failed to get access to docker shared preferences"
        }
        
        var prefToKeep = [String : AnyObject]()
        // loop over
        for key in prefs.dictionaryRepresentation().keys {
            if !prefKeysToKeep.contains({ $0.rawValue == key }) {
                prefs.removeObjectForKey(key)
            }
            else {
                prefToKeep[key] = prefs.objectForKey(key)
            }
        }
        
        // deleting app containers and group container
        let fileManager = NSFileManager.defaultManager()
        
        guard let dockerAppContainerPath: String = dockerAppContainerPath() else {
            return "failed to get path to Docker application container"
        }
        
        if fileManager.fileExistsAtPath(dockerAppContainerPath) {
            do {
                try fileManager.removeItemAtPath(dockerAppContainerPath)
            } catch {
                return "failed to delete Docker application container"
            }
        }
        
        guard let dockerHelperAppContainerPath: String = dockerHelperContainerPath() else {
            return "failed to get path to DockerHelper application container"
        }
        if fileManager.fileExistsAtPath(dockerHelperAppContainerPath) {
            do {
                try fileManager.removeItemAtPath(dockerHelperAppContainerPath)
            } catch {
                return "failed to delete DockerHelper application container"
            }
        }
        
        guard let groupContainerPath = Paths.groupContainerPath() else {
            return "failed to get path to DockerHelper application container"
        }
        if fileManager.fileExistsAtPath(groupContainerPath) {
            do {
                try fileManager.removeItemAtPath(groupContainerPath)
            } catch {
                return "failed to delete Docker application group container"
            }
        }
        // uninstall complete

        // recreate prefs for keys we want to keep
        if prefToKeep.count > 0 {
            guard Paths.groupContainerPath() != nil else {
                return "failed to get groupContainerPath"
            }
            guard let newPrefs: NSUserDefaults = NSUserDefaults(suiteName: Config.dockerAppGroupIdentifier) else {
                return "failed to get access to docker shared preferences"
            }
            for pref in prefToKeep {
                newPrefs.setValue(pref.1, forKey: pref.0)
            }
            newPrefs.synchronize()
        }
        
        return nil
    }
    
    ////////////////////////////////////////////////////////////////////////////
    
    // returns nil on error
    private static func dockerAppContainerPath() -> String? {
        let home: String = NSHomeDirectory()
        if home.containsString("/Library/Containers/") && home.hasSuffix("/Data") {
            // it means the app is sandboxed, we return the home directory value
            // after removing the last part "/Data"
            return home.substringToIndex(home.endIndex.advancedBy(-5)) // removes the 5 last characters
        }
        // the app isn't sandboxed, we handcraft the path to the app container
        guard let bundleID: String = NSBundle.mainBundle().bundleIdentifier else {
            return nil
        }
        return NSString.pathWithComponents([home, "Library", "Containers", bundleID])
    }
    
    // return nil on error
    private static func dockerHelperContainerPath() -> String? {
        guard var dockerAppContainerPath: String = dockerAppContainerPath() else {
            return nil
        }
        while dockerAppContainerPath.hasSuffix("/Containers") == false {
            guard let indexRange: Range<String.Index> = dockerAppContainerPath.rangeOfString("/", options: NSStringCompareOptions.BackwardsSearch) else {
                return nil
            }
            // remove last path component
            dockerAppContainerPath = dockerAppContainerPath.substringToIndex(indexRange.startIndex)
        }
        // add docker helper bundle identifier
        dockerAppContainerPath = NSString.pathWithComponents([dockerAppContainerPath, Config.dockerHelperBundleIdentifier])
        return dockerAppContainerPath
    }
}
