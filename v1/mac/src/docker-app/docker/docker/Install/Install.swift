//
//  install.swift
//  docker
//
//  Created by Emmanuel Briney on 14/12/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Cocoa
import Foundation
import ServiceManagement

class Install {
    
    // unprivileged installation (sandbox-compliant part of the installation)
    static func performUnprivilegedInstallation() -> String? {
        
        let machineMigrator = MachineMigrator()
        if let errorMessage: String = machineMigrator.migrateLocalMachine() {
            return errorMessage
        }
        
        // install symlinks in group container
        if let errorMessage: String = installSymlinksToCliBinariesInGroupContainer(["docker", "docker-compose", "docker-machine", "notary", "docker-diagnose"]) {
            return errorMessage
        }
        
        return nil
    }
    
    // performPrivilegedInstallation installs components that require
    // privileged access. There's only 2 things so far:
    // - vmnetd (registered in launchd as root daemon)
    // - symlinks in /usr/local/bin (vmnetd process can take care of
    // installing them)
    static func performPrivilegedInstallation() {
        
        let vmnetdBundleIdentifier = "com.docker.vmnetd"
        
        // Does the job already exist?
        let socketPath: String = "/var/tmp/" + vmnetdBundleIdentifier + ".socket"
        if NSFileManager.defaultManager().fileExistsAtPath(socketPath){
            // Connect to the socket: this allows us to distinguish between a working
            // install and a partially broken one (e.g. if the socket was left behind after
            // a failed uninstall).
            switch getVmnetdState(socketPath) {
            case VmnetdState.JustRight:
                Logger.log(level: Logger.Level.Notice, content: "probe of \(socketPath) successful: not reinstalling component")
                return
            case VmnetdState.Ancient:
                Logger.log(level: Logger.Level.Warning, content: "com.docker.vmnetd is too old for automatic upgrade.")
            case VmnetdState.NeedsUpgrade:
                Logger.log(level: Logger.Level.Notice, content: "com.docker.vmnetd is old but new enough to support automatic upgrade")
                // NOTE(aduermael) I don't know what "automatic upgrade" means
                // in that comment above. I've been reading the code in vmnetd
                // proxy.c, and so far uninstall only seems to be removing
                // vmnetd from launchd, it's not updating itself.
                if uninstall(socketPath) {
                    Logger.log(level: Logger.Level.Notice, content: "com.docker.vmnetd automatic uninstall successful")
                } else {
                    // Automatic upgrade of com.docker.vmnetd failed.
                    Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd automatic upgrade failed.")
                }
            case VmnetdState.Missing:
                Logger.log(level: Logger.Level.Notice, content: "com.docker.vmnetd component is missing and will be installed")
            }
        }
        
        // if we reach this point, it means vmnetd has to be installed
        // 3 cases:
        // - user is root -> don't need to prompt for admin password
        // - user is admin or not -> prompt to for admin credentials
        if Utils.userIsRoot {
            var err: Unmanaged<CFError>?
            var result = false
            
            // remove vmnetd (old versions can't be updated)
            // NOTE: SMJobRemove is deprecated, we can actually comment this
            // code if we don't have users with old vmnetd instances (that
            // can't be updated).
            result = SMJobRemove(kSMDomainSystemLaunchd, vmnetdBundleIdentifier, nil, true, &err)
            if !result {
                Logger.log(level: ASLLogger.Level.Error, content: "Failed to uninstall old networking components.")
            }
            
            // register vmnetd in launchd with root privileges. vmnetd creates
            // its socket here: /var/tmp/com.docker.vmnetd.socket
            result = SMJobBless(kSMDomainSystemLaunchd, vmnetdBundleIdentifier, nil, &err)
            if !result {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Failed to install networking components.")
            }
        } else { // user is not root (he can be admin or not)
            
            // If user is not root and the session is unattended, we log a fatal
            // error because we won't be able to prompt for admin password.
            if Options.unattended {
                Logger.log(level: ASLLogger.Level.Fatal, content: "installAsRootNonAdminUserError".localize())
                return
            }
            
            var authItem = AuthorizationItem(name: kSMRightBlessPrivilegedHelper, valueLength: 0, value: nil, flags: 0)
            var authRights: AuthorizationRights = AuthorizationRights(count: 1, items: &authItem)
            let authFlags: AuthorizationFlags = [
                AuthorizationFlags.Defaults,
                AuthorizationFlags.InteractionAllowed,
                AuthorizationFlags.ExtendRights,
                AuthorizationFlags.PreAuthorize
            ]
            
            var authRef: AuthorizationRef = nil
            var exit = false
            
            while authRef == nil {
                guard let msg = WizardMessageGeneric(message: "privilegedAccessWarningTitle".localize(), details: "privilegedAccessWarningMessage".localize(), icon: "MsgWarningIcon") else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                    return
                }
                
                let exitBtn = msg.addButton("Exit")
                exitBtn.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                    exit = true
                    message?.close()
                }
                
                let okBtn = msg.addCloseButton("Ok".localize())
                okBtn.makeDefault()
                
                analytics.track(AnalyticEvent.InstallAskForRootPassword)
                Wizard.show(msg)
                
                if exit {
                    NSApp.terminate(nil)
                    return
                }
                
                let status: OSStatus = AuthorizationCreate(&authRights, nil, authFlags, &authRef)
                if status != errAuthorizationSuccess {
                    // Authentication Failed
                    authRef = nil
                    continue
                }
                
                var err: Unmanaged<CFError>?
                var result = false
                
                // remove vmnetd (old versions can't be updated)
                // NOTE: SMJobRemove is deprecated, we can actually comment this
                // code if we don't have users with old vmnetd instances (that
                // can't be updated).
                result = SMJobRemove(kSMDomainSystemLaunchd, vmnetdBundleIdentifier, authRef, true, &err)
                if !result {
                    Logger.log(level: ASLLogger.Level.Error, content: "Failed to uninstall old networking components.")
                    authRef = nil
                    continue
                }
                
                // register vmnetd in launchd with root privileges. vmnetd creates
                // its socket here: /var/tmp/com.docker.vmnetd.socket
                result = SMJobBless(kSMDomainSystemLaunchd, vmnetdBundleIdentifier, authRef, &err)
                if !result {
                    Logger.log(level: ASLLogger.Level.Error, content: "Failed to install networking components.")
                    authRef = nil
                }
            }
        }
        
        Logger.log(level: Logger.Level.Notice, content: String(format: "privileged installation of %@ successful", arguments: [vmnetdBundleIdentifier]))
        // Probe the socket to trigger the privileged component to complete its install
        if !probe_vmnetd(socketPath) {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Communication with networking components failed.")
        }
    }
    
    // installSymlinks installs symlinks in /usr/local/bin that are pointing
    // to the ones in app container folder. The ones in container folder
    // can be updated without privileges, even with sandboxing.
    // installSymlinks asks vmnetd process to do the job, since it already
    // has root privileges, it avoids prompting the user just to fix symlinks.
    static func installSymlinks() {
        guard let appContainerPath = Paths.appContainerPath() else {
            return
        }
        install_symlinks("/var/tmp/com.docker.vmnetd.socket", container_folder: appContainerPath)
        
        // test if system symlinks are properly installed
        let testResult: (success: Bool, errorMessage: String?) = testCliSystemSymlinks(["docker", "docker-compose", "notary"])
        if let errorMessage = testResult.errorMessage {
            Logger.log(level: Logger.Level.Error, content: "Symlinks error: \(errorMessage)")
        } else if !testResult.success {
            Logger.log(level: Logger.Level.Error, content:"Symlinks are not correctly installed.")
        } else {
            Logger.log(level: Logger.Level.Notice, content:"Symlinks are valid.")
        }
    }
    
    // uninstallSymlinks uninstalls symlinks installed using installSymlinks()
    static func uninstallSymlinks() -> Bool {
        return uninstall_symlinks("/var/tmp/com.docker.vmnetd.socket")
    }
    
    // uninstallSockets uninstalls all sockets created (outside user's Container
    // folder) by the application.
    // Only /var/tmp/com.docker.vmnetd.socket is kept to keep communication
    // possible with vmnetd.
    static func uninstallSockets() -> Bool {
        return uninstall_sockets("/var/tmp/com.docker.vmnetd.socket")
    }
    
    static func uninstallVmnetd() -> Bool {
        return uninstall("/var/tmp/com.docker.vmnetd.socket")
    }
    
    // test if symlinks located in /usr/local/bin actually point to the symlinks
    // located in the group container.
    // returns a boolean indicating if the test succeeded or not and a string
    // which can contain an error message.
    static func testCliSystemSymlinks(binaries: Array<String>) -> (Bool, String?) {
        
        // get path to <group_container_path>/bin/
        guard let groupContainerUrl: NSURL = NSFileManager.defaultManager().containerURLForSecurityApplicationGroupIdentifier(Config.dockerAppGroupIdentifier) else {
            return (false, "failed to get group container URL for identifier \(Config.dockerAppGroupIdentifier)")
        }
        
        guard let groupContainerPath: String = groupContainerUrl.path else {
            return (false, "failed to get path from group container URL")
        }
        
        // construct path of the bin/ directory inside the group container
        let groupContainerBinDirectoryPath: String = NSString.pathWithComponents([groupContainerPath, "bin"])
        
        // path to /usr/local/bin
        let systemBinDirectoryPath: String = "/usr/local/bin"
        
        var systemSymlinkPath: String = ""
        var intendedSymlinkTarget: String = ""
        for binaryName in binaries {
            // path to system symlink
            systemSymlinkPath = NSString.pathWithComponents([systemBinDirectoryPath, binaryName])
            // check that it exists
            if NSFileManager.defaultManager().fileExistsAtPath(systemSymlinkPath) == false {
                return (false, nil)
            }
            // intended path to
            intendedSymlinkTarget = NSString.pathWithComponents([groupContainerBinDirectoryPath, binaryName])
            // check that the system symlink targets the group container symlink
            do {
                if try NSFileManager.defaultManager().destinationOfSymbolicLinkAtPath(systemSymlinkPath) != intendedSymlinkTarget {
                    return (false, nil)
                }
            } catch let error as NSError {
                return (false, error.localizedDescription)
            }
        }
        
        return (true, nil)
    }
    
    // tells if VM files are currently excluded from Time Machine backups
    static func isVmExcludedFromTimeMachineBackups() -> Bool {
        guard let appContainerPath = Paths.appContainerPath() else {
            return false
        }
        // first, test if the com.docker.driver.amd64-linux directory exists
        let filePath: String = NSString.pathWithComponents([appContainerPath, "com.docker.driver.amd64-linux"])
        guard NSFileManager.defaultManager().fileExistsAtPath(filePath) else {
            // file doesn't exist, return failure
            return false
        }
        // second, check if the file is already excluded from Time Machine backups
        let fileUrl: NSURL = NSURL.fileURLWithPath(filePath)
        var darwinFalse: DarwinBoolean = DarwinBoolean(false)
        let darwinFalsePointer: UnsafeMutablePointer<DarwinBoolean> = withUnsafeMutablePointer(&darwinFalse) { $0 }
        let fileIsExcluded: Bool = CSBackupIsItemExcluded(fileUrl, darwinFalsePointer)
        return fileIsExcluded
    }
    
    // exclude VM files from Time Machine backups
    static func excludeVmContentFromTimeMachineBackups(wantToExclude: Bool) -> Bool {
        guard let appContainerPath = Paths.appContainerPath() else {
            return false
        }
        // first, test if the com.docker.driver.amd64-linux directory exists
        let filePath: String = NSString.pathWithComponents([appContainerPath, "com.docker.driver.amd64-linux"])
        guard NSFileManager.defaultManager().fileExistsAtPath(filePath) else {
            // file doesn't exist, return failure
            return false
        }
        
        // second, check if the file is already excluded from Time Machine backups
        let fileUrl: NSURL = NSURL.fileURLWithPath(filePath)
        var darwinFalse: DarwinBoolean = DarwinBoolean(false)
        let darwinFalsePointer: UnsafeMutablePointer<DarwinBoolean> = withUnsafeMutablePointer(&darwinFalse) { $0 }
        let fileIsExcluded: Bool = CSBackupIsItemExcluded(fileUrl, darwinFalsePointer)
        
        if wantToExclude {
            // we want to exclude the vm
            if fileIsExcluded == false {
                // try to exclude the vm
                if CSBackupSetItemExcluded(fileUrl, true, false) != noErr {
                    return false // failure
                }
            }
        } else {
            // we want to un-exclude the vm
            if fileIsExcluded {
                // try to un-exclude the vm
                if CSBackupSetItemExcluded(fileUrl, false, false) != noErr {
                    return false // failure
                }
            }
        }
        return true // succcess
    }
    
}

