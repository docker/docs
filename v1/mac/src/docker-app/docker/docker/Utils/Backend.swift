//
//  Backend.swift
//  Docker
//
//  Created by Doby Mock on 4/15/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation

// Backend provides a set of functions that can be used to communicate with
// the backend. It also provides listeners for backend notifications.

class Backend: NSObject {

    enum Error {
    case Timeout
    case Process
    case UTF8
    case JSON
    case Invalid
    case Unknown
    case Message(String)
    }

    static func errorMessage(error: Error) -> String {
        switch error {
        case .Timeout: return "request timed out"
        case .Process: return "frontend process failed"
        case .UTF8:    return "UTF-8 decoding failed"
        case .JSON:    return "JSON decoding failed"
        case .Invalid: return "invalid return type"
        case .Unknown: return "unknown error"
        case .Message(let message): return message
        }
    }

    static func sendRequest(json: String, completionQueue: dispatch_queue_t!, retry: Int = 0, completion: ((error: Error?, response: String) -> Void)? = nil) {
        
        // when the backend is not ready, we may want to retry instead of
        // failing. Backend.vmGetSettings & Backend.getSharedDirectories
        // retry 5 times for instance.
        func sendRequestAgain(retry: Int = 0) -> Bool {
            if retry > 0 {
                dispatch_after(dispatch_time(DISPATCH_TIME_NOW, Int64(NSEC_PER_SEC)), dispatch_get_main_queue()) { () -> Void in
                    sendRequest(json, completionQueue: completionQueue, retry: retry - 1, completion: completion)
                }
                return true
            } else {
                return false
            }
        }
        
        let backgroundQueue = dispatch_get_global_queue(QOS_CLASS_BACKGROUND, 0)
        dispatch_async(backgroundQueue, {
            let bundlePath: String = NSBundle.mainBundle().bundlePath
            let executablePath: String = NSString.pathWithComponents([bundlePath, "Contents", "MacOS", "com.docker.frontend"])
            let task = NSTask()
            let out = NSPipe()
            // let err = NSPipe()
            task.launchPath = executablePath
            task.arguments = [json]
            task.standardOutput = out
            // task.standardError = err
            task.launch()
            task.waitUntilExit()
            let exitCode = task.terminationStatus
            let response = String(data: out.fileHandleForReading.readDataToEndOfFile(), encoding: NSUTF8StringEncoding)
            // _ = String(data: err.fileHandleForReading.readDataToEndOfFile(), encoding: NSUTF8StringEncoding)
            out.fileHandleForReading.closeFile()
            // err.fileHandleForReading.closeFile()
            if exitCode == 1 {
                // error
                if !sendRequestAgain(retry) {
                    dispatch_async(completionQueue, { () -> Void in
                        completion?(error: Error.Process, response: "")
                    })
                }
                return
            }
            if exitCode == 2 {
                // timeout
                if !sendRequestAgain(retry) {
                    dispatch_async(completionQueue, { () -> Void in
                        completion?(error: Error.Timeout, response: "")
                    })
                }
                return
            }
            guard let responseValue: String = response else {
                // error
                if !sendRequestAgain(retry) {
                    dispatch_async(completionQueue, { () -> Void in
                        completion?(error: Error.UTF8, response: "")
                    })
                }
                return
            }
            dispatch_async(completionQueue, { () -> Void in
                completion?(error: nil, response: responseValue)
            })
        })
    }
    
    
    static func sendRequestSync(json: String) -> (error: Error?, response: String) {
        let bundlePath: String = NSBundle.mainBundle().bundlePath
        let executablePath: String = NSString.pathWithComponents([bundlePath, "Contents", "MacOS", "com.docker.frontend"])
        let task = NSTask()
        let out = NSPipe()
        // let err = NSPipe()
        task.launchPath = executablePath
        task.arguments = [json]
        task.standardOutput = out
        // task.standardError = err
        task.launch()
        task.waitUntilExit()
        let exitCode = task.terminationStatus
        let response = String(data: out.fileHandleForReading.readDataToEndOfFile(), encoding: NSUTF8StringEncoding)
        // _ = String(data: err.fileHandleForReading.readDataToEndOfFile(), encoding: NSUTF8StringEncoding)
        out.fileHandleForReading.closeFile()
        // err.fileHandleForReading.closeFile()
        if exitCode == 1 {
            // error
            return (error: Error.Process, response: "")
        }
        if exitCode == 2 {
            // timeout
            return (error: Error.Timeout, response: "")
        }
        guard let responseValue: String = response else {
            return (error: Error.UTF8, response: "")
        }
        return (error: nil, response: responseValue)
    }
    
    // restart vm
    static func vmRestart(completion: ((error: Error?) -> Void)? = nil) {
        Backend.sendRequest("{\"action\":\"restartvm\"}", completionQueue: dispatch_get_main_queue(), completion: { (error: Error?, response: String) in
            // check if there was an error while communication with the backend
            if let errorValue = error {
                completion?(error: errorValue)
                return
            }
            // check that the response is a UTF-8 string
            guard let responseData: NSData = response.dataUsingEncoding(NSUTF8StringEncoding) else {
                completion?(error: Error.UTF8)
                return
            }
            // parse response and check whether it is a success or failure
            let responseObject: [String: AnyObject]!
            do {
                responseObject = try NSJSONSerialization.JSONObjectWithData(responseData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
            } catch {
                completion?(error: Error.JSON)
                return
            }
            // handle response
            // get the status value from the response
            guard let responseStatus: String = responseObject["status"] as? String else {
                completion?(error: Error.Invalid)
                return
            }
            // handle case where response contains an error message
            if responseStatus != "ok" {
                if responseStatus == "error" {
                    // backend returned an error
                    guard let errorMessage: String = responseObject["message"] as? String else {
                        completion?(error: Error.Unknown)
                        return
                    }
                    completion?(error: Error.Message(errorMessage))
                    return
                } else {
                    completion?(error: Error.Invalid)
                    return
                }
            }
            completion?(error: nil)
        })
    }
    
    
    
    /// memory in GBs
    /// cups in cpu count
    static func vmGetSettings(completion: ((error: Error?, settings: (memory: UInt, cpus: UInt, daemonjson: String, proxyHttpSystem: String, proxyHttpsSystem: String, proxyExcludeSystem: String, proxyHttpOverride: String, proxyHttpsOverride: String, proxyExcludeOverride: String)?) -> Void)? = nil)  {
        
        Backend.sendRequest("{\"action\":\"getvmsettings\"}", completionQueue: dispatch_get_main_queue(), retry: 5, completion: { (error: Error?, response: String) in
            
            // check if there was an error while communication with the backend
            if let errorValue: Error = error {
                completion?(error: errorValue, settings: nil)
                return
            }
            // check that the response is a UTF-8 string
            guard let responseData: NSData = response.dataUsingEncoding(NSUTF8StringEncoding) else {
                completion?(error: Error.UTF8, settings: nil)
                return
            }
            // parse response and check whether it is a success or failure
            let responseObject: [String: AnyObject]!
            do {
                responseObject = try NSJSONSerialization.JSONObjectWithData(responseData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
            } catch {
                completion?(error: Error.JSON, settings: nil)
                return
            }
            // handle response
            // get the status value from the response
            guard let responseStatus: String = responseObject["status"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            // handle case where response contains an error message
            if responseStatus != "ok" {
                if responseStatus == "error" {
                    // backend returned an error
                    guard let errorMessage: String = responseObject["message"] as? String else {
                        completion?(error: Error.Unknown, settings: nil)
                        return
                    }
                    completion?(error: Error.Message(errorMessage), settings: nil)
                    return
                } else {
                    completion?(error: Error.Invalid, settings: nil)
                    return
                }
            }
            // otherwise there was no error
            guard let memoryValue: UInt = responseObject["memory"] as? UInt else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let cpusValue: UInt = responseObject["cpus"] as? UInt else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let daemonjsonValue: String = responseObject["daemonjson"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let proxyHttpSystemValue: String = responseObject["systemProxyHttp"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let proxyHttpsSystemValue: String = responseObject["systemProxyHttps"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let proxyExcludeSystemValue: String = responseObject["systemProxyExclude"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let proxyHttpOverrideValue: String = responseObject["overrideProxyHttp"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let proxyHttpsOverrideValue: String = responseObject["overrideProxyHttps"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            guard let proxyExcludeOverrideValue: String = responseObject["overrideProxyExclude"] as? String else {
                completion?(error: Error.Invalid, settings: nil)
                return
            }
            let settingsObject: (memory: UInt, cpus: UInt, daemonjson: String, proxyHttpSystem: String, proxyHttpsSystem: String, proxyExcludeSystem: String, proxyHttpOverride: String, proxyHttpsOverride: String, proxyExcludeOverride: String) = (memoryValue, cpusValue, daemonjsonValue, proxyHttpSystemValue, proxyHttpsSystemValue, proxyExcludeSystemValue, proxyHttpOverrideValue, proxyHttpsOverrideValue, proxyExcludeOverrideValue)
            completion?(error: nil, settings: settingsObject)
            return
        })
    }
    
    /// send request to the backend to set VM settings
    /// memory in GBs
    /// cpus in cpu count
    static func vmSetSettings(settings: (memory: UInt?, cpus: UInt?, daemonjson: String?, proxyHttpOverride: String?, proxyHttpsOverride: String?, proxyExcludeOverride: String?), completion: ((error: Error?) -> Void)?) {
        
        var args = [String: AnyObject]()
        
        if let memory = settings.memory {
            args["memory"] = memory
        }
        if let cpus = settings.cpus {
            args["cpus"] = cpus
        }
        if let daemonjson = settings.daemonjson {
            args["daemonjson"] = daemonjson
        }
        if let proxyHttp = settings.proxyHttpOverride {
            args["overrideProxyHttp"] = proxyHttp
        }
        if let proxyHttps = settings.proxyHttpsOverride {
            args["overrideProxyHttps"] = proxyHttps
        }
        if let proxyExclude = settings.proxyExcludeOverride {
            args["overrideProxyExclude"] = proxyExclude
        }
        
        let json: [String: AnyObject] = [
            "action": "setvmsettings",
            "args": args
        ]
        
        do {
            let data = try NSJSONSerialization.dataWithJSONObject(json, options: NSJSONWritingOptions())
            // convert the request data into a UTF-8 string
            guard let string = NSString(data: data, encoding: NSUTF8StringEncoding) as? String else {
                completion?(error: Error.UTF8)
                return
            }
            Backend.sendRequest(string, completionQueue: dispatch_get_main_queue(), completion: { (error: Error?, response: String) in
                // check if there was an error while communication with the backend
                if let errorValue: Error = error {
                    completion?(error: errorValue)
                    return
                }
                // check that the response is a UTF-8 string
                guard let responseData: NSData = response.dataUsingEncoding(NSUTF8StringEncoding) else {
                    completion?(error: Error.UTF8)
                    return
                }
                // parse response and check whether it is a success or failure
                let responseObject: [String: AnyObject]!
                do {
                    responseObject = try NSJSONSerialization.JSONObjectWithData(responseData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
                } catch {
                    completion?(error: Error.JSON)
                    return
                }
                // get the status value from the response
                guard let responseStatus: String = responseObject["status"] as? String else {
                    completion?(error: Error.Invalid)
                    return
                }
                // handle case where response contains an error message
                if responseStatus != "ok" {
                    if responseStatus == "error" {
                        // backend returned an error
                        guard let errorMessage: String = responseObject["message"] as? String else {
                            completion?(error: Error.Unknown)
                            return
                        }
                        completion?(error: Error.Message(errorMessage))
                        return
                    } else {
                        completion?(error: Error.Invalid)
                        return
                    }
                }
                // otherwise it is a success
                completion?(error: nil)
            })
        } catch _ {
            completion?(error: Error.JSON)
            return
        }
    }
    
    enum DockerState: String {
        case Starting = "starting"
        case Running = "running"
        case Unknown = "unknown"
    }
    
    // MARK: file sharing
    
    //
    typealias GetSharedDirectoriesCallback = (error: Error?, directories: [String]) -> Void
    static func getSharedDirectories(completion: GetSharedDirectoriesCallback?) {
        Backend.sendRequest("{\"action\":\"getshareddirectories\"}", completionQueue: dispatch_get_main_queue(), retry: 5, completion: { (error: Error?, response: String) in
            
            // check if there was an error while communication with the backend
            if let errorValue: Error = error {
                completion?(error: errorValue, directories: [String]())
                return
            }
            // check that the response is a UTF-8 string
            guard let responseData: NSData = response.dataUsingEncoding(NSUTF8StringEncoding) else {
                completion?(error: Error.UTF8, directories: [String]())
                return
            }
            // parse response and check whether it is a success or failure
            let responseObject: [String: AnyObject]!
            do {
                responseObject = try NSJSONSerialization.JSONObjectWithData(responseData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
            } catch {
                completion?(error: Error.JSON, directories: [String]())
                return
            }
            // handle response
            // get the status value from the response
            guard let responseStatus: String = responseObject["status"] as? String else {
                completion?(error: Error.Invalid, directories: [String]())
                return
            }
            // handle case where response contains an error message
            if responseStatus != "ok" {
                if responseStatus == "error" {
                    // backend returned an error
                    guard let errorMessage: String = responseObject["message"] as? String else {
                        completion?(error: Error.Unknown, directories: [String]())
                        return
                    }
                    completion?(error: Error.Message(errorMessage), directories: [String]())
                    return
                } else {
                    completion?(error: Error.Invalid, directories: [String]())
                    return
                }
            }
            // otherwise there was no error
            guard let sharedDirectoriesValue: [String] = responseObject["directories"] as? [String] else {
                completion?(error: Error.Invalid, directories: [String]())
                return
            }
            completion?(error: nil, directories: sharedDirectoriesValue)
        })
    }
    
    //
    typealias SetSharedDirectoriesCallback = (error: Error?, directories: [String]) -> Void
    static func setSharedDirectories(directories: [String], completion: SetSharedDirectoriesCallback? ) {
        
        let args: [String: AnyObject] = [
            "directories": directories
        ]
        let json: [String: AnyObject] = [
            "action": "setshareddirectories",
            "args": args
        ]
        
        do {
            let data = try NSJSONSerialization.dataWithJSONObject(json, options: NSJSONWritingOptions())
            // convert the request data into a UTF-8 string
            guard let string = NSString(data: data, encoding: NSUTF8StringEncoding) as? String else {
                completion?(error: Error.UTF8, directories: [String]())
                return
            }
            Backend.sendRequest(string, completionQueue: dispatch_get_main_queue(), completion: { (error: Error?, response: String) in
                // check if there was an error while communication with the backend
                if let errorValue: Error = error {
                    completion?(error: errorValue, directories: [String]())
                    return
                }
                // check that the response is a UTF-8 string
                guard let responseData: NSData = response.dataUsingEncoding(NSUTF8StringEncoding) else {
                    completion?(error: Error.UTF8, directories: [String]())
                    return
                }
                // parse response and check whether it is a success or failure
                let responseObject: [String: AnyObject]!
                do {
                    responseObject = try NSJSONSerialization.JSONObjectWithData(responseData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
                } catch {
                    completion?(error: Error.JSON, directories: [String]())
                    return
                }
                // get the status value from the response
                guard let responseStatus: String = responseObject["status"] as? String else {
                    completion?(error: Error.Invalid, directories: [String]())
                    return
                }
                // handle case where response contains an error message
                if responseStatus != "ok" {
                    if responseStatus == "error" {
                        // backend returned an error
                        guard let errorMessage: String = responseObject["message"] as? String else {
                            completion?(error: Error.Unknown, directories: [String]())
                            return
                        }
                        completion?(error: Error.Message(errorMessage), directories: [String]())
                        return
                    } else {
                        completion?(error: Error.Invalid, directories: [String]())
                        return
                    }
                }
                // otherwise it is a success
                guard let sharedDirectoriesValue: [String] = responseObject["directories"] as? [String] else {
                    completion?(error: Error.Invalid, directories: [String]())
                    return
                }
                completion?(error: nil, directories: sharedDirectoriesValue)
            })
        } catch _ {
            completion?(error: Error.JSON, directories: [String]())
            return
        }
    }
    
    // by default, let's consider Docker is starting
    // and ask backend for an event when the state changes
    // (dockerState is thread safe)
    private static var dockerStateLock = NSLock()
    private static var _dockerState: DockerState = .Unknown
    private static var dockerState: DockerState {
        get {
            dockerStateLock.lock()
            var res = _dockerState
            if _enforceStarting {
                res = .Starting
            }
            dockerStateLock.unlock()
            return res
        }
        set {
            dockerStateLock.lock()
            _dockerState = newValue
            dockerStateLock.unlock()
        }
    }
    
    private static var _enforceStarting: Bool = true
    static func enforceStarting() {
        dockerStateLock.lock()
        _enforceStarting = true
        dockerStateLock.unlock()
        NSNotificationCenter.defaultCenter().postNotificationName(dockerStateNotificationName, object: dockerState.rawValue)
    }
    
    private static var startedToListen = false
    private static var dockerStateAskAgainDelay = 0.5
    static let dockerStateNotificationName: String = "dockerStateChanged"
    
    private static var lock = NSLock()
    
    static func getCurrentState() -> DockerState {
        // if we never started to listen, it means we did not even
        // call the function to launch the backend, so the Swift app is starting
        // or the user is simply entering its admin password.
        // it's confusing to see an animated icon at this moment, so let's
        // consider Docker is running.
        if startedToListen == false {
            return .Running
        }
        return dockerState
    }
    
    static func listenForDockerStateChanges() {
        
        startedToListen = true
        enforceStarting()
        
        let backgroundQueue = dispatch_get_global_queue(QOS_CLASS_BACKGROUND, 0)
        dispatch_async(backgroundQueue, {
    
            while true {
                let args: [String: AnyObject] = [
                    "vmstate": dockerState.rawValue,
                ]
                let json: [String: AnyObject] = [
                    "action": "vmstateevent",
                    "args": args
                ]
                
                var data: NSData
                
                do {
                    data = try NSJSONSerialization.dataWithJSONObject(json, options: NSJSONWritingOptions())
                } catch _ {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "failed to create valid json")
                    sleep(1)
                    continue
                }
                
                guard let string = NSString(data: data, encoding: NSUTF8StringEncoding) as? String else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "failed to convert to UTF-8")
                    sleep(1)
                    continue
                }

                let response = Backend.sendRequestSync(string)

                if let error = response.error {
                    switch error {
                    case .Timeout: continue
                    default:
                        let msg = errorMessage(error)
                        Logger.log(level: ASLLogger.Level.Error, content: msg)
                        // The backend is not ready to answer, retry after delay:

                        // TODO: use dispatch_after instead, sleep can
                        // block other operations in the same queue
                        sleep(1)
                        continue
                    }
                }
                
                guard let respData = response.response.dataUsingEncoding(NSUTF8StringEncoding) else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "vmstateevent result is not valid UTF-8")
                    return
                }
                
                let serverResponse: [String: AnyObject]!
                do {
                    serverResponse = try NSJSONSerialization.JSONObjectWithData(respData, options: .AllowFragments) as? [String: AnyObject]
                } catch {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "vmstateevent invalid json: \(response.response)")
                    return
                }
                
                if let status: String = serverResponse["status"] as? String {
                    if status == "ok" {
                        if let vmstateStr = serverResponse["vmstate"] as? String {
                            guard let vmstate = DockerState(rawValue: vmstateStr) else {
                                Logger.log(level: ASLLogger.Level.Fatal, content: "unknown vmstate")
                                return
                            }
                            dockerState = vmstate
                            Logger.log(level: ASLLogger.Level.Notice, content: "dockerState = \(dockerState)")
                            _enforceStarting = false
                            NSNotificationCenter.defaultCenter().postNotificationName(dockerStateNotificationName, object: dockerState.rawValue)
                            
                        } else {
                            Logger.log(level: ASLLogger.Level.Fatal, content: "vmstateevent response from server")
                            return
                        }
                    }
                }
            }
        })
    }
}
