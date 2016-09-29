//
//  Config.swift
//  docker-installer
//
//  Created by Michel Courtine on 12/5/15.
//  Copyright © 2015 Docker Inc. All rights reserved.
//

// Config contains constant values defining some aspects of the application's
// behavior
class Config {
    
    // MARK: Minimum requirements for the application to launch
    
    static let OSXMinimumMajor = 10
    static let OSXMinimumMinor = 10
    static let OSXMinimumPatch = 3
    
    /// minimum amount of memory allocated to the VM
    static let minimumVMMemory: UInt64 = 1024 * 1024 * 1024 // (1GB)
    
    // default amount of memory allocated to the VM
    static let defaultVMMemory: UInt64 = 1024 * 1024 * 1024 * 2 // (2GB)
    
    /// minimum amount of physical memory required for the host
    static let minimumMemoryRequirement: UInt64 = 1024 * 1024 * 1024 * 4 // (4GB)
    
    
    // MARK: bundle and group identifiers
    
    /// identifier of the docker application group (used to create the path in ~/Library/Group Containers)
    static let dockerAppGroupIdentifier: String = "group.com.docker"
    
    /// application bundle identifier of the DockerHelper application
    static let dockerHelperBundleIdentifier: String = "com.docker.helper"
}
