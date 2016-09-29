//
//  cliInstall.swift
//  docker
//
//  Created by Gaetan de Villele on 1/8/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation

//
func installSymlinksToCliBinariesInGroupContainer(binaries: Array<String>) -> String? {

    // get path to <group_container_path>/bin/
    guard let groupContainerPath = Paths.groupContainerPath() else {
        return "failed to get groupContainerPath"
    }
    // construct path of the bin/ directory inside the group container
    let groupContainerBinPath: String = NSString.pathWithComponents([groupContainerPath, "bin"])
    // test if bin/ directory exists and try to create it if it doesn't
    let errorMessage: String? = Paths.CreateDirectoryIfNecessary(groupContainerBinPath)
    if errorMessage != nil {
        // error
        return errorMessage
    }
    
    // get path to the cli binaries located in the bundle
    let bundlePath: String = NSBundle.mainBundle().bundlePath
    let bundleBinDirectoryPath: String = NSString.pathWithComponents([bundlePath, "Contents", "Resources", "bin"])
    
    for binary in binaries {
        let symlinkPath: String = NSString.pathWithComponents([groupContainerBinPath, binary])
        let realPath: String = NSString.pathWithComponents([bundleBinDirectoryPath, binary])
        // check the binary exists in the bundle
        if NSFileManager.defaultManager().fileExistsAtPath(realPath) == false {
            // binary is missing in bundle
            return "executable missing in bundle: \(binary)"
        }
        // check if symlink already exists and delete it if it does.
        // Note that fileExistsAtPath will return false if the symlink exists but the target does not.
        // We need to include the case of dangling symlinks, e.g. if the user downloads the app to
        // a different directory.
        do {
            let _: NSDictionary = try NSFileManager.defaultManager().attributesOfItemAtPath(symlinkPath)
            do {
                try NSFileManager.defaultManager().removeItemAtPath(symlinkPath)
            } catch let error as NSError {
                return "failed to delete existing symlink for \(binary): \(error.localizedDescription)"
            }
        } catch {
            // File mustn't exist
        }
        // actually create the symlink in the group container
        do {
            try NSFileManager.defaultManager().createSymbolicLinkAtPath(symlinkPath, withDestinationPath: realPath)
        } catch let error as NSError {
            return "failed to create symlink for \(binary): \(error.localizedDescription)"
        }
    }
    return nil
}
