//
//  Utils.swift
//  docker-installer
//
//  Created by Michel Courtine on 12/6/15.
//  Copyright © 2015 David Scott. All rights reserved.
//

import Foundation
import Cocoa

class Utils {
    
    // cache for the following function
    static private var _buildType: String? = nil
    static func getBuildType() -> String {
        // use cache if possible
        if let buildType = _buildType {
            return buildType
        }
        // no cache available
        if let _buildType = NSBundle.mainBundle().objectForInfoDictionaryKey("DockerBuildType") as? String {
            return _buildType
        }
        _buildType = "dev"
        return "dev"
    }
    
    static func isDevBuild() -> Bool {
        return getBuildType() == "dev"
    }
    
    static func isStableBuild() -> Bool {
        return getBuildType() == "stable"
    }
    
    static func isBetaBuild() -> Bool {
        return getBuildType() == "beta"
    }
    
    static func getChannel() -> String {
        if getBuildType() == "stable" {
            return "stable"
        }
        return "beta"
    }
    
    // returns ISO representation of the current OS language.
    // returns an empty string if the language code cannot be obtained.
    static func getCurrentLanguage() -> String? {
        let locale: [String: String] = NSLocale.componentsFromLocaleIdentifier(NSLocale.currentLocale().localeIdentifier)
        return locale[NSLocaleLanguageCode]
    }
    
    static func getOSXMajorVersion() -> Int {
        let osVersion: NSOperatingSystemVersion = NSProcessInfo().operatingSystemVersion
        return osVersion.majorVersion
    }
    
    static func getOSXMinorVersion() -> Int {
        let osVersion: NSOperatingSystemVersion = NSProcessInfo().operatingSystemVersion
        return osVersion.minorVersion
    }
    
    static func getOSXPatchVersion() -> Int {
        let osVersion: NSOperatingSystemVersion = NSProcessInfo().operatingSystemVersion
        return osVersion.patchVersion
    }
    
    static func getOSXVersionString() -> String {
        return NSProcessInfo().operatingSystemVersionString
    }
    
    // MARK: user and permissions
    
    // returns whether the current user is the root user
    static var userIsRoot: Bool {
        get {
            return Darwin.geteuid() == 0
        }
    }
    
    // return whether the current user is an Administrator
    static var userIsAdministrator: Bool {
        get {
            return user_info_is_admin()
        }
    }
    
    // return whether the current user is an administrator or the root user
    static var userIsAdminOrRoot: Bool {
        get {
            return (userIsAdministrator || userIsRoot)
        }
    }
    
    static func getDirectorySizeInByte(path: String) -> Int{
        var dirSize = 0
        if let dirURL = NSURL(string: path) {
            var bool: ObjCBool = false
            if let dirPath = dirURL.path {
                if NSFileManager().fileExistsAtPath(dirPath, isDirectory: &bool) {
                    if bool.boolValue {
                        let fileManager =  NSFileManager.defaultManager()
                        
                        if let filesEnumerator = fileManager.enumeratorAtURL(dirURL, includingPropertiesForKeys: nil, options: [], errorHandler: {
                            (url, error) -> Bool in
                            print(url.path)
                            print(error.localizedDescription)
                            return true
                        }) {
                            while let fileURL = filesEnumerator.nextObject() as? NSURL {
                                if let filePath = fileURL.path {
                                    do {
                                        let attributes = try fileManager.attributesOfItemAtPath(filePath) as NSDictionary
                                        dirSize += attributes.fileSize().hashValue
                                    } catch let error as NSError {
                                        print(error.localizedDescription)
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        return dirSize
    }
    
    // humanReadableByteCount formats bytes into human readable String
    static func humanReadableByteCount(bytes: Double) -> String {
        if (bytes < 1024) {
            return String(bytes) + " B"
        }
        var exp: Int = Int(log(bytes) / log(1024))
        // don't allow units above G
        // (1T will be 1000G)
        if exp > 3 {
            exp = 3
        }
        let pre: String = String("kMGTPE".substringWithRange(NSRange(location: exp-1, length: 1)))
        return String(format: "%.f %@B", bytes / pow(1024, Double(exp)), pre)
    }
    
    // bytesToGigaBytes converts bytes to gigabytes
    static func bytesToGigaBytes(bytes: Double) -> Int {
        let exp = 3
        return Int(bytes / pow(1024, Double(exp)))
    }
    
    // bytesToGigaBytes converts gigabytes to bytes
    static func gigaBytesToBytes(giga: UInt) -> UInt {
        return giga * 1024 * 1024 * 1024
    }
}

// Execute a shell command and return error code
func Execute(command: String) -> Int32 {
    let task = NSTask()
    task.launchPath = "/bin/sh"
    task.arguments = ["-c", command]
    task.launch()
    task.waitUntilExit()
    let exitCode = task.terminationStatus
    return exitCode
}
