//
//  Preferences.swift
//  Docker
//
//  Created by Gaetan de Villele on 3/25/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation

// singleton
public class Preferences {
    
    static let preferencesDidChangeNotification: String = "dockerPreferencesDidChangeNotification"
    
    static private var instance: Preferences? = nil
    static var sharedInstance: Preferences {
        get {
            if let result: Preferences = Preferences.instance {
                return result
            }
            if let result: Preferences = Preferences() {
                Preferences.instance = result
                return result
            }
            Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
            exit(EXIT_FAILURE)
        }
    }
    
    // constants containing the string values of the preference keys
    //
    // /!\ IF YOU ADD A NEW KEY IN THIS ENUM, PLEASE ADD IT ALSO IN THE
    // /!\ "allValues" CONST ARRAY
    //
    enum Key: String {
        // preferences window
        case autoStart = "prefAutoStart" // Bool
        case checkForUpdates = "prefCheckForUpdates" // Bool
        case localMachineMigration = "prefLocalMachineMigration" // UInt
        // analytics
        case analyticsUserID = "analyticsUserID" // String
        case analyticsEvents = "analyticsEvents" // [String]
        case analyticsEnabled = "analyticsEnabled" // Bool
        case autoSendCrashReports = "autoSendCrashReports" // Bool
        // diagnostic
        case diagnosticID = "diagnosticID" // String
        // other
        case installWelcomedUser = "installWelcomedUser" // Bool
        case installDisplayedPopover = "installDisplayedPopover" // Bool
        case dockerAppLaunchPath = "DockerAppLaunchPath" // String
        case dockerBuildType = "DockerBuildType" // String
        case appVersionHistory = "appVersionHistory" // [Dictionary<String, AnyObject>]
        // OLD KEYS THAT AREN'T USED ANYMORE
        case memory = "prefMemory" // UInt // NOT USED ANYMORE
        
        
        
        static let allValues = [autoStart, memory, checkForUpdates, localMachineMigration, analyticsUserID, analyticsEvents, installWelcomedUser, installDisplayedPopover, dockerAppLaunchPath, appVersionHistory, analyticsEnabled, autoSendCrashReports]
    }
    
    
    // default values for preferences
    private static let defaultValues: Dictionary<Key, AnyObject> = [
        Key.autoStart: true,
        Key.memory: UInt(Config.defaultVMMemory),
        Key.checkForUpdates: true,
        Key.localMachineMigration: UInt(MachineMigrator.localMachineMigrationStatus.None.rawValue),
        Key.analyticsUserID: "",
        Key.analyticsEvents: [String](),
        Key.analyticsEnabled: true,
        Key.diagnosticID: "",
        Key.autoSendCrashReports: true,
        Key.installWelcomedUser: false,
        Key.installDisplayedPopover: false,
        Key.dockerAppLaunchPath: "",
        Key.dockerBuildType: "",
        Key.appVersionHistory: [Dictionary<String, AnyObject>](),
        ]
    
    
    
    // MARK: attributes
    
    // UserDefaults of the docker application group
    private var userDefaults: NSUserDefaults
    
    
    
    // MARK: initializer
    
    // initializer
    init?() {
        Paths.groupContainerPath()
        // get docker application group userDefaults
        guard let ud: NSUserDefaults = NSUserDefaults(suiteName: Config.dockerAppGroupIdentifier) else {
            return nil
        }
        self.userDefaults = ud
    }
    
    
    
    // MARK: type-enforce storage and retrival
    
    func keyExists(key: Preferences.Key) -> Bool {
        return userDefaults.objectForKey(key.rawValue) != nil
    }
    
    func removeKey(key: Preferences.Key) {
        userDefaults.removeObjectForKey(key.rawValue)
    }
    
    func getAutoStart() -> Bool {
        guard let result: Bool = getBooleanForKey(Key.autoStart) else {
            guard let result: Bool = Preferences.defaultValues[Key.autoStart] as? Bool else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return true
            }
            return result
        }
        return result
    }
    
    func setAutoStart(value: Bool) -> Bool {
        return setBoolean(value, forKey: Key.autoStart)
    }
    
    func getMemory() -> UInt {
        // TODO: gdevillele: use UInt here (store UInt in preferences)
        guard let result: Int = getIntegerForKey(Key.memory) else {
            guard let result: UInt = Preferences.defaultValues[Key.memory] as? UInt else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return 0
            }
            return result
        }
        return UInt(result)
    }
    
    func setMemory(value: UInt) -> Bool {
        // TODO: gdevillele: use UInt here
        return setInteger(Int(value), forKey: Key.memory)
    }
    
    func getCheckForUpdates() -> Bool {
        guard let result: Bool = getBooleanForKey(Key.checkForUpdates) else {
            guard let result: Bool = Preferences.defaultValues[Key.checkForUpdates] as? Bool else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return true
            }
            return result
        }
        return result
    }
    
    func setCheckForUpdates(value: Bool) -> Bool {
        return setBoolean(value, forKey: Key.checkForUpdates)
    }
    
    func getLocalMachineMigration() -> UInt {
        // TODO: gdevillele: use UInt here (store UInt in preferences)
        guard let result: Int = getIntegerForKey(Key.localMachineMigration) else {
            guard let result: UInt = Preferences.defaultValues[Key.localMachineMigration] as? UInt else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return 0
            }
            return result
        }
        return UInt(result)
    }
    
    func setLocalMachineMigration(value: UInt) -> Bool {
        // TODO: gdevillele: use UInt here
        return setInteger(Int(value), forKey: Key.localMachineMigration)
    }
    
    func getAnalyticsUserId() -> String {
        guard let result: String = getStringForKey(Key.analyticsUserID) else {
            guard let result: String = Preferences.defaultValues[Key.analyticsUserID] as? String else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return ""
            }
            return result
        }
        return result
    }
    
    func setAnalyticsUserId(value: String) -> Bool {
        return setString(value, forKey: Key.analyticsUserID)
    }
    
    func getAnalyticsEvents() -> [String] {
        guard let result: [String] = getArrayForKey(Key.analyticsEvents) as? [String] else {
            guard let result: [String] = Preferences.defaultValues[Key.analyticsEvents] as? [String] else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return [String]()
            }
            return result
        }
        return result
    }
    
    func setAnalyticsEvents(value: [String]) -> Bool {
        return setArray(value, forKey: Key.analyticsEvents)
    }
    
    
    func getAnalyticsEnabled() -> Bool {
        guard let result: Bool = getBooleanForKey(Key.analyticsEnabled) else {
            guard let result: Bool = Preferences.defaultValues[Key.analyticsEnabled] as? Bool else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return true
            }
            return result
        }
        return result
    }
    
    func setAnalyticsEnabled(value: Bool, notify: Bool = true) -> Bool {
        return setBoolean(value, forKey: Key.analyticsEnabled, notify: notify)
    }
    
    func getAutoSendCrashReports() -> Bool {
        guard let result: Bool = getBooleanForKey(Key.autoSendCrashReports) else {
            guard let result: Bool = Preferences.defaultValues[Key.autoSendCrashReports] as? Bool else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return true
            }
            return result
        }
        return result
    }
    
    func getDiagnosticId() -> String {
        guard let result: String = getStringForKey(Key.diagnosticID) else {
            guard let result: String = Preferences.defaultValues[Key.diagnosticID] as? String else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return ""
            }
            return result
        }
        return result
    }
    
    func setDiagnosticId(value: String) -> Bool {
        return setString(value, forKey: Key.diagnosticID)
    }
    
    func setAutoSendCrashReports(value: Bool, notify: Bool = true) -> Bool {
        return setBoolean(value, forKey: Key.autoSendCrashReports, notify: notify)
    }
    
    
    func getInstallWelcomedUser() -> Bool {
        guard let result: Bool = getBooleanForKey(Key.installWelcomedUser) else {
            guard let result: Bool = Preferences.defaultValues[Key.installWelcomedUser] as? Bool else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return false
            }
            return result
        }
        return result
    }
    
    func setInstallWelcomedUser(value: Bool) -> Bool {
        return setBoolean(value, forKey: Key.installWelcomedUser)
    }
    
    func getInstallDisplayedPopover() -> Bool {
        guard let result: Bool = getBooleanForKey(Key.installDisplayedPopover) else {
            guard let result: Bool = Preferences.defaultValues[Key.installDisplayedPopover] as? Bool else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return false
            }
            return result
        }
        return result
    }
    
    func setInstallDisplayedPopover(value: Bool) -> Bool {
        return setBoolean(value, forKey: Key.installDisplayedPopover)
    }
    
    func getDockerAppLaunchPath() -> String {
        guard let result: String = getStringForKey(Key.dockerAppLaunchPath) else {
            guard let result: String = Preferences.defaultValues[Key.dockerAppLaunchPath] as? String else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return ""
            }
            return result
        }
        return result
    }
    
    func setDockerAppLaunchPath(value: String) -> Bool {
        return setString(value, forKey: Key.dockerAppLaunchPath)
    }
    
    func getAppVersionHistory() -> [Dictionary<String, AnyObject>] {
        guard let result: [Dictionary<String, AnyObject>] = getArrayForKey(Key.appVersionHistory) as? [Dictionary<String, AnyObject>] else {
            guard let result: [Dictionary<String, AnyObject>] = Preferences.defaultValues[Key.appVersionHistory] as? [Dictionary<String, AnyObject>] else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return [Dictionary<String, AnyObject>]()
            }
            return result
        }
        return result
    }
    
    func setAppVersionHistory(value: [Dictionary<String, AnyObject>]) -> Bool {
        return setArray(value, forKey: Key.appVersionHistory)
    }
    
    func getDockerBuildType() -> String {
        guard let result: String = getStringForKey(Key.dockerBuildType) as String? else {
            guard let result: String = Preferences.defaultValues[Key.dockerBuildType] as? String ?? "" else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
                return ""
            }
            return result
        }
        return result
    }
    
    func setDockerBuildType(value: String) -> Bool {
        return setString(value, forKey: Key.dockerBuildType)
    }
    
    
    // MARK: private value storage and retrieval
    
    private func synchronize(notify: Bool = true) -> Bool {
        // From 10.12 method comments
        /*!
         -synchronize is deprecated and will be marked with the NS_DEPRECATED macro in a future release.
         
         -synchronize blocks the calling thread until all in-progress set operations have completed. This is no longer necessary. Replacements for previous uses of -synchronize depend on what the intent of calling synchronize was. If you synchronized...
         - ...before reading in order to fetch updated values: remove the synchronize call
         - ...after writing in order to notify another program to read: the other program can use KVO to observe the default without needing to notify
         - ...before exiting in a non-app (command line tool, agent, or daemon) process: call CFPreferencesAppSynchronize(kCFPreferencesCurrentApplication)
         - ...for any other reason: remove the synchronize call
         */
        if Utils.getOSXMajorVersion() == 10 && Utils.getOSXMinorVersion() <= 11 {
            // synchronise() can return false when it's not an error. Perhaps this means
            // it didn't synchronise because it didn't need to, because the background
            // synchronisation has already happened.
            if !userDefaults.synchronize() {
                Logger.log(level: ASLLogger.Level.Notice, content: "NSUserDefaults.synchronise() returned false")
            }
        }
        if notify {
            NSNotificationCenter.defaultCenter().postNotificationName(Preferences.preferencesDidChangeNotification, object: nil)
        }
        return true
    }
    
    private func setString(string: String, forKey key: Preferences.Key) -> Bool {
        userDefaults.setObject(string, forKey: key.rawValue)
        return synchronize()
    }
    
    private func getStringForKey(key: Preferences.Key) -> String? {
        return userDefaults.stringForKey(key.rawValue)
    }
    
    private func setInteger(integer: Int, forKey key: Preferences.Key) -> Bool {
        userDefaults.setInteger(integer, forKey: key.rawValue)
        return synchronize()
    }
    
    private func getIntegerForKey(key: Preferences.Key) -> Int? {
        if keyExists(key) {
            return userDefaults.integerForKey(key.rawValue)
        }
        return nil
    }
    
    private func setUnsigned(unsigned: UInt, forKey key: Preferences.Key) -> Bool {
        userDefaults.setObject(unsigned, forKey: key.rawValue)
        return synchronize()
    }
    
    private func getUnsignedForKey(key: Preferences.Key) -> UInt? {
        if keyExists(key) {
            if let obj: AnyObject? = userDefaults.objectForKey(key.rawValue) {
                if let unsigned: UInt = obj as? UInt {
                    return unsigned
                }
            }
        }
        return nil
    }
    
    private func setBoolean(boolean: Bool, forKey key: Preferences.Key, notify: Bool = true) -> Bool {
        userDefaults.setBool(boolean, forKey: key.rawValue)
        return synchronize(notify)
    }
    
    private func getBooleanForKey(key: Preferences.Key) -> Bool? {
        if keyExists(key) {
            return userDefaults.boolForKey(key.rawValue)
        }
        return nil
    }
    
    private func getArrayForKey(key: Preferences.Key) -> [AnyObject]? {
        return userDefaults.arrayForKey(key.rawValue)
    }
    
    private func setArray(array: [AnyObject], forKey key: Preferences.Key) -> Bool {
        userDefaults.setObject(array, forKey: key.rawValue)
        return synchronize()
    }
    
    // setObject sets an object for a given Key
    private func setObject(object: AnyObject?, forKey key: Preferences.Key) -> Bool {
        userDefaults.setObject(object, forKey: key.rawValue)
        return synchronize()
    }
    
    // getValue gets an object associated to a given Key
    private func getObjectForKey(key: Preferences.Key) -> AnyObject? {
        return userDefaults.objectForKey(key.rawValue)
    }
}
