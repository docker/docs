//
//  BackendSettings.swift
//  docker
//
//  Created by Doby Mock on 6/27/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation

// TODO(aduermael): BackendSettings is currently used by different setting panels.
// But we can make it lighter, specific and private to each panel.

// represents a given state of backend settings
class BackendSettings: Equatable {
    
    // keys of daemon.json
    static internal let DaemonInsecureRegistriesKey: String = "insecure-registries"
    static internal let DaemonRegistryMirrorsKey: String = "registry-mirrors"
    
    var memory: Double // in bytes
    var cpus: Double
    var daemon: [String: AnyObject]
    var proxyHttpSystem: String
    var proxyHttpsSystem: String
    var proxyExcludeSystem: String
    var proxyHttpOverride: String
    var proxyHttpsOverride: String
    var proxyExcludeOverride: String
    
    // init with default values
    init() {
        self.cpus = 1
        self.memory = Double(Config.minimumVMMemory)
        self.daemon = [String: AnyObject]()
        self.proxyHttpSystem = String()
        self.proxyHttpsSystem = String()
        self.proxyExcludeSystem = String()
        self.proxyHttpOverride = String()
        self.proxyHttpsOverride = String()
        self.proxyExcludeOverride = String()
    }
    
    init(memory: Double, cpus: Double, daemonjson: String, // swiftlint:disable:this function_parameter_count
         proxyHttpSystem: String,
         proxyHttpsSystem: String,
         proxyExcludeSystem: String,
         proxyHttpOverride: String,
         proxyHttpsOverride: String,
         proxyExcludeOverride: String) {
        self.memory = memory
        self.cpus = cpus
        self.daemon = [String: AnyObject]()
        self.proxyHttpSystem = proxyHttpSystem
        self.proxyHttpsSystem = proxyHttpsSystem
        self.proxyExcludeSystem = proxyExcludeSystem
        self.proxyHttpOverride = proxyHttpOverride
        self.proxyHttpsOverride = proxyHttpsOverride
        self.proxyExcludeOverride = proxyExcludeOverride
        // NOTE(aduermael): in some case we don't want to carry the
        // daemonjson value, so the String may be empty. 
        // It's ok if we can't parse it.
        if let daemonObject = parseJson(daemonjson) {
            self.daemon = daemonObject
        }
    }
    
    init(value: BackendSettings) {
        self.memory = value.memory
        self.cpus = value.cpus
        self.daemon = value.daemon
        self.proxyHttpSystem = value.proxyHttpSystem
        self.proxyHttpsSystem = value.proxyHttpsSystem
        self.proxyExcludeSystem = value.proxyExcludeSystem
        self.proxyHttpOverride = value.proxyHttpOverride
        self.proxyHttpsOverride = value.proxyHttpsOverride
        self.proxyExcludeOverride = value.proxyExcludeOverride
    }
    
    // read-write
    var insecureRegistry: [String]? {
        get {
            // check if there is an insecure-registries entry
            guard let entry: AnyObject = self.daemon[BackendSettings.DaemonInsecureRegistriesKey] else {
                // no insecure-registries key found, it is not an error as it
                // means that there is no insecure-registries. We return an empty
                // array.
                return [String]()
            }
            // there is an insecure-registries entry, we try casting it into [String]
            guard let ir: [String] = entry as? [String] else {
                // could not cast it. It means that it has not the expected type.
                // This is not supposed to happen
                return nil
            }
            return ir
        }
        set {
            self.daemon[BackendSettings.DaemonInsecureRegistriesKey] = newValue
        }
    }
    
    // read-write
    var registryMirror: [String]? {
        get {
            // check if there is an registry-mirrors entry
            guard let entry: AnyObject = self.daemon[BackendSettings.DaemonRegistryMirrorsKey] else {
                // no registry-mirrors key found, it is not an error as it
                // means that there is no registry-mirrors. We return an empty
                // array.
                return [String]()
            }
            // there is an insecure-registries entry, we try casting it into [String]
            guard let rm: [String] = entry as? [String] else {
                // could not cast it. It means that it has not the expected type.
                // This is not supposed to happen
                return nil
            }
            return rm
        }
        set {
            self.daemon[BackendSettings.DaemonRegistryMirrorsKey] = newValue
        }
    }
    
    /// cleaned daemon object
    /// read-only
    var daemonClean: [String: AnyObject]? {
        get {
            // copy self.daemon and clean it by removing all the empty strings
            var daemonClean: [String: AnyObject] = self.daemon
            // clean insecure-registries entry if there is any.
            if let entry: AnyObject = daemonClean[BackendSettings.DaemonInsecureRegistriesKey] {
                // cast entry into [String]
                guard let ir: [String] = entry as? [String] else {
                    return nil
                }
                // only keep the non-empty strings
                let cleanedIr = ir.filter({ (element) -> Bool in
                    return element.characters.count > 0
                })
                // if array is empty after cleaning, we remove it
                if cleanedIr.isEmpty {
                    daemonClean.removeValueForKey(BackendSettings.DaemonInsecureRegistriesKey)
                } else {
                    daemonClean[BackendSettings.DaemonInsecureRegistriesKey] = cleanedIr
                }
            }
            // clean registry-mirrors entry if there is any.
            if let entry: AnyObject = daemonClean[BackendSettings.DaemonRegistryMirrorsKey] {
                // cast entry into [String]
                guard let rm: [String] = entry as? [String] else {
                    return nil
                }
                // only keep the non-empty strings
                let cleanedRm = rm.filter({ (element) -> Bool in
                    return element.characters.count > 0
                })
                // if array is empty after cleaning, we remove it
                if cleanedRm.isEmpty {
                    daemonClean.removeValueForKey(BackendSettings.DaemonRegistryMirrorsKey)
                } else {
                    daemonClean[BackendSettings.DaemonRegistryMirrorsKey] = cleanedRm
                }
            }
            return daemonClean
        }
    }
    
    /// json string representation of daemonClean
    /// read-only
    var daemonCleanJson: String? {
        get {
            // get daemonClean object
            guard let daemonClean: [String: AnyObject] = self.daemonClean else {
                return nil
            }
            // serialize daemonClean into a Json string
            do {
                let jsonData: NSData = try NSJSONSerialization.dataWithJSONObject(daemonClean, options: NSJSONWritingOptions(rawValue: 0))
                guard let jsonString: String = String(data: jsonData, encoding: NSUTF8StringEncoding) else {
                    return nil
                }
                return jsonString
            } catch {
                return nil
            }
        }
    }
    
    // utility functions
    
    /// Converts a json string into a Foundation object
    /// - parameter json: the json string to be parsed
    /// - returns: a foundation object corresponding to the json string, or nil
    ///            if there was any error.
    private func parseJson(json: String) -> [String: AnyObject]? {
        // convert the json string into NSData
        guard let jsonData = json.dataUsingEncoding(NSUTF8StringEncoding) else {
            return nil
        }
        // try to parse json data into a foundation object
        do {
            let jsonObjectOptional: [String: AnyObject]? = try NSJSONSerialization.JSONObjectWithData(jsonData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
            guard let jsonObject: [String: AnyObject] = jsonObjectOptional else {
                return nil
            }
            return jsonObject
        } catch {
            return nil
        }
    }
}
// Equatable protocol implementation for the BackendSettings class
func == (lhs: BackendSettings, rhs: BackendSettings) -> Bool {
    return lhs.memory == rhs.memory &&
        lhs.cpus == rhs.cpus
        && lhs.daemonClean! == rhs.daemonClean! // swiftlint:disable:this force_unwrapping
        && lhs.proxyHttpSystem == rhs.proxyHttpSystem
        && lhs.proxyHttpsSystem == rhs.proxyHttpsSystem
        && lhs.proxyExcludeSystem == rhs.proxyExcludeSystem
        && lhs.proxyHttpOverride == rhs.proxyHttpOverride
        && lhs.proxyHttpsOverride == rhs.proxyHttpsOverride
        && lhs.proxyExcludeOverride == rhs.proxyExcludeOverride
}

public func == (lhs: [String: AnyObject], rhs: [String: AnyObject] ) -> Bool {
    return NSDictionary(dictionary: lhs).isEqualToDictionary(rhs)
}
