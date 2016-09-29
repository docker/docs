//
//  UpdateManager.swift
//  docker
//
//  Created by Gaetan de Villele on 3/9/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Sparkle

// this class is a singleton
class UpdateManager {
    
    // shared instance
    static private var instance: UpdateManager? = nil
    static func sharedManager() -> UpdateManager {
        // if there is an instance we return it
        if let instance: UpdateManager = UpdateManager.instance {
            return instance
        }
        let newInstance: UpdateManager = UpdateManager()
        UpdateManager.instance = newInstance
        return newInstance
    }
    
    // MARK: static functions
    
    /// application updates are available only on non-dev builds and if the user
    /// is an administrator on the Mac
    static func isAppUpdateAvailable() -> Bool {
        return !Utils.isDevBuild()
    }
    
    // Sparkle updater delegate
    private let updaterDelegate: UpdaterDelegate = UpdaterDelegate()
    // auto check for update interval (in seconds)
    private let updateCheckInterval: NSTimeInterval = 3600 * 24 // 24h
    
    // auto check for update
    var automaticallyChecksForUpdates: Bool {
        get {
            if UpdateManager.isAppUpdateAvailable() {
                return Sparkle.SUUpdater.sharedUpdater().automaticallyChecksForUpdates
            }
            return false
        }
        set {
            if UpdateManager.isAppUpdateAvailable() {
                Sparkle.SUUpdater.sharedUpdater().automaticallyChecksForUpdates = newValue
            }
        }
    }
    
    // called when user clicks on the "check for updates" item in whale menu
    func checkForUpdates(sender: AnyObject) {
        if UpdateManager.isAppUpdateAvailable() {
            Sparkle.SUUpdater.sharedUpdater().checkForUpdates(sender)
        }
    }
    
    // this must be called once the application setup assistant (install process)
    // is finished to tell this manager that it now can automatically check for
    // updates and eventually show a popup.
    // This is done to avoid showing an "update available" popup while the user 
    // is still going through the setup/installation of Docker for Mac.
    func allowAutoCheckForUpdates() {
        self.updaterDelegate.allowedToCheckForUpdates = true
    }
    
    // initializer
    private init() {
        // by default the Sparkle Updater object is not allowed to do any
        // automatic check for updates or anything
        // (default values are set in Info.plist)
        // We override those values here, only if the current build is NOT a dev
        // build.
        if UpdateManager.isAppUpdateAvailable() {
            // set sparkle updater delegate
            Sparkle.SUUpdater.sharedUpdater().delegate = self.updaterDelegate
            // enable auto-checking depending on the preference value. If the
            // preference value does not exist, we enable the auto-checking as it is
            // the default behavior
            let autoCheck: Bool = Preferences.sharedInstance.getCheckForUpdates()
            Sparkle.SUUpdater.sharedUpdater().automaticallyChecksForUpdates = autoCheck
            // set the update check interval
            Sparkle.SUUpdater.sharedUpdater().updateCheckInterval = self.updateCheckInterval
            // set "send system profile" option
            Sparkle.SUUpdater.sharedUpdater().sendsSystemProfile = false
            // set "automatically download updates" option
            Sparkle.SUUpdater.sharedUpdater().automaticallyDownloadsUpdates = false
        }
    }
}

internal class UpdaterDelegate: NSObject, Sparkle.SUUpdaterDelegate {
    
    var allowedToCheckForUpdates: Bool = false
    
    internal func updaterMayCheckForUpdates(updater: SUUpdater!) -> Bool {
        return self.allowedToCheckForUpdates
    }
}
