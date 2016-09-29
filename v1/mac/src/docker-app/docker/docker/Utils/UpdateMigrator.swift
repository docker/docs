//
//  UpdateMigrator.swift
//  docker
//
//  Created by Gaetan de Villele on 3/28/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Cocoa


// this class has been introduced with version "1.10.3-beta5"
//
// /!\ ALL MIGRATIONS MUST BE ABLE TO BE EXECUTED SAFELY ON A FRESH INSTALL
class UpdateMigrator {

    //
    static var startVersion: String = ""

    //
    static func askUserToApplyChannelModifications() {
        let currentBuildType: String = Utils.getBuildType()
        let previousBuildType: String = Preferences.sharedInstance.getDockerBuildType()
        let currentChannel = currentBuildType == "stable" ? "stable" : "beta"
        let nextChannel = previousBuildType == "stable" ? "stable" : "beta"
        if let message = WizardMessageGeneric(message: "switchChannelResetTitle".localize(currentChannel), details: "switchChannelResetMessage".localize(nextChannel, nextChannel), icon: "MsgWarningIcon") {
            message.addButton("Cancel".localize()).blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                message?.close()
            }

            message.addButton("switchChannelResetBtnOk".localize()).makeDefault().blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                message?.close()
                //NB: We need to be sure that vmnetd has the same version than the excutable else uninstall will fail
                Install.performPrivilegedInstallation()
                let error: String? = Uninstall.uninstallWhileKeepingPreferenceKeys([Preferences.Key.analyticsUserID])
                if let error: String = error {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "factoryResetFailed".localize(error))
                } else {
                    // write the channel type
                    if let currentBuildType: String = NSBundle.mainBundle().objectForInfoDictionaryKey("DockerBuildType") as? String {
                        Preferences.sharedInstance.setDockerBuildType(currentBuildType)
                    }
                    // restart application
                    let task = NSTask()
                    task.launchPath = "/usr/bin/open"
                    task.arguments = [NSBundle.mainBundle().bundlePath]
                    task.launch()
                    exit(EXIT_SUCCESS)
                }
            }
            Wizard.show(message)
        }
    }

    //
    static func checkChannelChanges() -> Bool {
        let currentBuildType: String = Utils.getBuildType()
        let previousBuildType: String = Preferences.sharedInstance.getDockerBuildType()

        // During automated tests, we throw an error if the build type has changed
        // (which indicates that an old version has not been uninstalled)
        if Options.unattended {
            if !previousBuildType.isEmpty && currentBuildType != previousBuildType {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Uninstall previous version before running tests")
                exit(EXIT_FAILURE)
            }
            return false
        }

        // On old version we cannot get build type
        if previousBuildType.isEmpty {
            return false
        }

        // There is no changes in channel
        if currentBuildType == previousBuildType {
            return false
        }

        // We are not on a stable channel or not switching on the stable channel
        // (master, beta and dev are considered as the same channel)
        if currentBuildType != "stable" && previousBuildType != "stable" {
            return false
        }

        // We are switching from/to stable channel
        askUserToApplyChannelModifications()
        return true
    }

    //
    static func migrate() -> String? {
        Logger.log(level: ASLLogger.Level.Notice, content: "migrator: start migration process")

        guard let currentAppVersion: String = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleShortVersionString") as? String else {
            return "failed to get current version"
        }
        var previousVersions: [Dictionary<String, AnyObject>] = Preferences.sharedInstance.getAppVersionHistory()
        // TODO: filter previous versions
        // - if current version is stable: remove all entries that are not stable migrations
        // - if current version is beta: remove all entries that are not beta migrations
        // - detect dev channel

        // if there is no previous version registered it means it is a fresh install.
        if previousVersions.isEmpty {
            Logger.log(level: ASLLogger.Level.Notice, content: "migrator: no previous version detected. This is a fresh install.")
            // write the current version in the history
            previousVersions.append(["version": currentAppVersion, "type": "install"])
            if Preferences.sharedInstance.setAppVersionHistory(previousVersions) == false {
                let error = "failed to update app version history"
                Logger.log(level: ASLLogger.Level.Error, content: "migrator: failed to store cureent version in history (\(currentAppVersion)) (\(error))")
                return error
            }
        } else {
            // There was at least one previous versions.
            // We find the most recently installed version or migration, and
            // try to perform all the migrations since then.
            Logger.log(level: ASLLogger.Level.Notice, content: "migrator: previous install detected")
            guard let lastVersionRecord: Dictionary<String, AnyObject> = previousVersions.last else {
                return "failed to get the most recent version previously installed"
            }
            guard let lastVersionName: String = lastVersionRecord["version"] as? String else {
                return "failed to get the name of most recent version previously installed"
            }
            Logger.log(level: ASLLogger.Level.Notice, content: "migrator: previous version is \(lastVersionName)")
            startVersion = lastVersionName

            var migrations = migrationsBeta
            if Utils.isStableBuild() {
                migrations = migrationsStable
            }

            // loop over all migrations and construct the list of the migrations
            // we have to perform in the current situation.
            // We take the migration in the right order and store the migrations
            // that are more recent than the last version installed.
            var migrationsToPerform: [Migration] = []
            var lastVersionDetected: Bool = false
            for migration in migrations {
                if lastVersionDetected {
                    migrationsToPerform.append(migration)
                } else {
                    if migration.version == lastVersionName {
                        lastVersionDetected = true
                    }
                }
            }
            if migrationsToPerform.isEmpty == false {
                Logger.log(level: ASLLogger.Level.Notice, content: "migrator: \(migrationsToPerform.count) migrations to perform")
                for migration in migrationsToPerform {
                    Logger.log(level: ASLLogger.Level.Notice, content: "migrator: migrating for version \(migration.version) ...")
                    if let error: String = migration.function() {
                        return "[\(migration.version)] [\(error)]"
                    }
                    // if the performed migration is for the current bundle version
                    // we flag the migration as "install" instead of a regular "migration"
                    let versionHistoryType: String = migration.version == currentAppVersion ? "install" : "migration"
                    previousVersions.append(["version": migration.version, "type": versionHistoryType])
                    if Preferences.sharedInstance.setAppVersionHistory(previousVersions) == false {
                        let error = "failed to update app version history"
                        Logger.log(level: ASLLogger.Level.Error, content: "migrator: migration failed (\(migration.version)) (\(error))")
                        return error
                    }
                    Logger.log(level: ASLLogger.Level.Notice, content: "migrator: migration succeeded for version \(migration.version)")
                }
            } else {
                // no migration to perform
                Logger.log(level: ASLLogger.Level.Notice, content: "migrator: no migration needed, you are good to go!")
            }
            // if the app has been updated, we can display a nice popup to the user
            if currentAppVersion != lastVersionName {
                Logger.log(level: ASLLogger.Level.Notice, content: "migrator: successfully updated from version \(lastVersionName) to \(currentAppVersion)")
                if !Utils.isDevBuild() {
                    let notification: NSUserNotification = NSUserNotification()
                    notification.title = "dockerUpdatedNotificationTitle".localize()
                    notification.informativeText = "dockerUpdatedNotificationInfo".localize(currentAppVersion)
                    notification.hasActionButton = false
                    NSUserNotificationCenter.defaultUserNotificationCenter().deliverNotification(notification)
                }
            }
        }
        Logger.log(level: ASLLogger.Level.Notice, content: "migrator: end of migration process")
        return nil
    }

    enum ChannelType {
        case Dev
        case Beta
        case Stable
    }

    static private func channelType(version: String) -> ChannelType {
        if migrationsStable.contains({ $0.version == version }) {
            return ChannelType.Stable
        }
        if migrationsBeta.contains({ $0.version == version }) {
            return ChannelType.Beta
        }
        return ChannelType.Dev
    }

    // describe a migration
    struct Migration {
        var version: String
        var function: () -> String?
    }

    // MARK: stable migration
    static private let migrationsStable: [Migration] = [
        Migration(version: "1.12.0-a", function: UpdateMigrator.empty_migration), //
    ]

    // MARK: beta migrations

    // list of all the migrations that ever existed
    static private let migrationsBeta: [Migration] = [
        Migration(version: "1.10.3-beta5", function: UpdateMigrator.beta5_migration), // migration from pre-beta5 to beta5
        Migration(version: "1.11.0-beta6", function: UpdateMigrator.empty_migration), // migration from beta5 to beta6
        Migration(version: "1.11.0-beta7", function: UpdateMigrator.empty_migration), // migration from beta6 to beta7
        Migration(version: "1.11.0-beta8", function: UpdateMigrator.empty_migration), // migration from beta7 to beta8
        Migration(version: "1.11.0-beta9", function: UpdateMigrator.beta9_migration), // migration from beta8 to beta9
        Migration(version: "1.11.1-beta10", function: UpdateMigrator.empty_migration), // migration from beta9 to beta10
        Migration(version: "1.11.1-beta11", function: UpdateMigrator.empty_migration), // migration from beta10 to beta11
        Migration(version: "1.11.1-beta12", function: UpdateMigrator.empty_migration), // migration from beta11 to beta12
        Migration(version: "1.11.1-beta13", function: UpdateMigrator.empty_migration), // migration from beta12 to beta13
        Migration(version: "1.11.1-beta13.1", function: UpdateMigrator.beta13_1_migration), // migration from beta13 to beta13.1
        Migration(version: "1.11.1-beta14", function: UpdateMigrator.empty_migration), // migration from beta13 to beta14
        Migration(version: "1.11.2-beta15", function: UpdateMigrator.empty_migration), // migration from beta14 to beta15
        Migration(version: "1.12.0-rc2-beta16", function: UpdateMigrator.empty_migration), // migration from beta15 to beta16
        Migration(version: "1.12.0-rc2-beta17", function: UpdateMigrator.empty_migration), // migration from beta16 to beta17
        Migration(version: "1.12.0-rc3-beta18", function: UpdateMigrator.empty_migration), // migration from beta17 to beta18
        Migration(version: "1.12.0-rc4-beta19", function: UpdateMigrator.empty_migration), // migration from beta18 to beta19
        Migration(version: "1.12.0-rc4-beta20", function: UpdateMigrator.empty_migration), // migration from beta19 to beta20
        Migration(version: "1.12.0-beta21", function: UpdateMigrator.empty_migration), // migration from beta20 to beta21
        Migration(version: "1.12.0-beta22", function: UpdateMigrator.empty_migration), // migration from beta21 to beta22
        Migration(version: "1.12.1-rc1-beta23", function: UpdateMigrator.empty_migration), // migration from beta22 to beta23
        Migration(version: "1.12.1-beta24", function: UpdateMigrator.empty_migration), // migration from beta23 to beta24
        Migration(version: "1.12.1-beta25", function: UpdateMigrator.empty_migration), // migration from beta24 to beta25
        Migration(version: "1.12.1-beta26", function: UpdateMigrator.empty_migration), // migration from beta25 to beta26
        Migration(version: "1.12.2-rc1-beta27", function: UpdateMigrator.empty_migration), // migration from beta26 to beta27
    ]

    // migration to be used if there is no migration to perform
    static private func empty_migration() -> String? {
        return nil // OK
    }

    // updates memory preference from 1GB to 2GB to match actual VM memory
    static private func beta5_migration() -> String? {
        // update the RAM preference from 1GB to 2GB
        if Preferences.sharedInstance.keyExists(Preferences.Key.memory) {
            if Preferences.sharedInstance.getMemory() == 1024 * 1024 * 1024 { // 1GB
                if Preferences.sharedInstance.setMemory(1024 * 1024 * 1024 * 2) == false { // 2GB
                    return "failed to update memory preference"
                }
            }
        }
        return nil // OK
    }

    // removes memory preference from the NSUserDefaults
    static private func beta9_migration() -> String? {
        if Preferences.sharedInstance.keyExists(Preferences.Key.memory) {
            Preferences.sharedInstance.removeKey(Preferences.Key.memory)
            if Preferences.sharedInstance.keyExists(Preferences.Key.memory) {
                return "failed to remove memory key" // ERROR
            }
        }
        return nil // OK
    }

    // flush logs if any (this migration never fails)
    static private func beta13_1_migration() -> String? {
        var appContainerPath: String = ""
        let home: String = NSHomeDirectory()
        if home.containsString("/Library/Containers/") && home.hasSuffix("/Data") {
            // it means the app is sandboxed, we return the home directory value
            appContainerPath = home
        } else {
            // the app isn't sandboxed, we handcraft the path to the app container
            guard let bundleID: String = NSBundle.mainBundle().bundleIdentifier else {
                // return "Failed to access the container directory to flush the logs"
                return nil // OK
            }
            appContainerPath = NSString.pathWithComponents([home, "Library", "Containers", bundleID, "Data"])
        }
        let logDirectoryPath: String = NSString.pathWithComponents([appContainerPath, "com.docker.driver.amd64-linux", "log"])
        if Paths.directoryExists(logDirectoryPath) == false {
            // if no log directory, we consider it a success
            return nil // OK
        }
        do {
            let filesToDelete: [String] = ["acpid.log", "dmesg", "docker.log", "messages", "proxy-vsockd.log", "vsudd.log", "wtmp", "messages.0", "messages.1", "messages.2", "messages.3", "messages.4", "messages.5", "messages.6", "messages.7", "messages.8", "messages.9"]
            let items: [String] = try NSFileManager.defaultManager().contentsOfDirectoryAtPath(logDirectoryPath)
            for item in items {
                if filesToDelete.contains(item) {
                    do {
                        try NSFileManager.defaultManager().removeItemAtPath(NSString.pathWithComponents([logDirectoryPath, item]))
                        // <item> removed
                    } catch {

                    }
                }
            }
        } catch {

        }
        return nil // OK
    }
}
