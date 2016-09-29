//
//  Logger.swift
//  docker
//
//  Created by Adrian Duermael on 12/9/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Foundation
import Cocoa

enum LoggerError: ErrorType {
    case UnableToOpenLogfile(log : String)
    case UnableToCreateDirectory(dir : String, msg : String)
}

class ASLLogger {
    
    enum Level: Int32 {
        case Panic = 1
        case Fatal = 2
        case Error = 3
        case Warning = 4
        case Notice = 5
        case Info = 6
        case Debug = 7
        
        /* asl api
        #define ASL_LEVEL_EMERG   0
        #define ASL_LEVEL_ALERT   1
        #define ASL_LEVEL_CRIT    2
        #define ASL_LEVEL_ERR     3
        #define ASL_LEVEL_WARNING 4
        #define ASL_LEVEL_NOTICE  5
        #define ASL_LEVEL_INFO    6
        #define ASL_LEVEL_DEBUG   7 */
    }
    
    // MARK: Static variables
    
    static private let bundle = "com.docker.app"
    static private let queue = dispatch_queue_create(bundle + ".sync_queue", DISPATCH_QUEUE_SERIAL) //maintain output order in serial queue
    
    // MARK: Static functions
    
    // Helper method to replace macro ASL_FILTER_MASK_UPTO
    static private func asl_filter_mask_upto(level: Level) -> Int32 {
        return (1 << (level.rawValue + 1)) - 1
    }
    
    // MARK: Class variables

    // sender should be the name of the app: "Docker"
    private let sender: String
    // facility should be the a process identifier like "com.docker.db"
    private let facility: String
    private let asl: asl_object_t
    private var fd: Int32 = -1
    
    init(_ sender: String, _ facility: String, _ logPath: String) throws {
        self.facility = facility
        self.sender = sender
        self.asl = asl_open(self.sender, self.facility, 0)
    }

    // Log swift strings
    func log(level: Level, _ text: String) {
        // Set timestamp before queing
        let ts = NSDate().timeIntervalSince1970
        let ts_nsec: Int64 = Int64(ts * 1_000_000_000)
        let ts_sec = Int64(ts)
        let ts_nsec_frac = ts_nsec - (ts_sec * 1_000_000_000)
        let s = String(ts_nsec) + " " + self.facility + " " + text + "\n"
        
        dispatch_async(ASLLogger.queue) {
            String(ts_nsec_frac).withCString { (s_nsec) in
                String(ts_sec).withCString { (s_sec) in
                    text.withCString { (content) in
                        let msg: asl_object_t = asl_new(UInt32(ASL_TYPE_MSG))
                        asl_set(msg, ASL_KEY_SENDER, self.sender)
                        asl_set(msg, ASL_KEY_FACILITY, self.facility)
                        asl_set(msg, ASL_KEY_TIME_NSEC, s_nsec)
                        asl_set(msg, ASL_KEY_TIME, s_sec)
                        asl_vlog(self.asl, msg, level.rawValue, "%s", getVaList([content]))
                        let stderr = NSFileHandle.fileHandleWithStandardError()

                        if let x = s.dataUsingEncoding(NSUTF8StringEncoding) {
                            stderr.writeData(x)
                        } else {
                            let safe = String(ts_nsec) + " " + self.facility + " (omitting non-UTF8 text)\n"
                            /* surely this can't fail by construction */
                            if let x = safe.dataUsingEncoding(NSUTF8StringEncoding) {
                                stderr.writeData(x)
                            }
                        }
                    }
                }
            }
        }
    }

    // shutdown closes log files and releases asl clients
    deinit {
        dispatch_sync(ASLLogger.queue) { // wait for last entry to be logged
            return
        }
        asl_release(self.asl)
        if self.fd != -1 {
            close(self.fd)
        }
    }
}

class Logger {
    // Logger won't log anything if enabled == false (true by default)
    static var enabled: Bool = true
    static private let globalLock = NSLock()
    static private var logClients: [String: ASLLogger] = [String: ASLLogger]()
    static private let sender = "Docker"
    // default facility for the Swift app
    static private let facility = "com.docker.docker"
    
    typealias Level = ASLLogger.Level
    
    static func log(facility: String = Logger.facility, level: ASLLogger.Level, content: String) {
        guard let appContainerPath = Paths.appContainerPath() else {
            return
        }
        let logDirParentPath = NSString.pathWithComponents([appContainerPath, "logs"])
        globalLock.lock()
        if enabled {
            if let l = Logger.logClients[facility] {
                l.log(level, content)
            } else {
                if let l = try? ASLLogger(sender, facility, logDirParentPath) {
                    Logger.logClients[facility] = l
                    l.log(level, content)
                } else {
                    print("unable to initialise logger for %s", facility) // TODO: Fail better
                }
            }
        }
        globalLock.unlock() // TODO The lock is only necessary to access logClients, could release earlier
        
        // once all the "text logging" is done, we ultimately test if the log was a fatal log
        if level == ASLLogger.Level.Fatal {
            Logger.terminate(content)
        }
    }
    
    // displays a "fatal error" popup containing the error message and call
    // NSApp.terminate() when the user closes this popup
    static private func terminate(msg: String) {
        if Process.arguments.contains("--unattended") == false {
            if let errorMessage = WizardMessageGeneric(message: "Fatal Error", details: msg, icon: "FatalIcon") {
                
                let diagnoseAndFeedbackBtn = errorMessage.addButton("fatalErrorDiagnoseAndFeedbackBtn".localize())
                diagnoseAndFeedbackBtn.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                    
                    guard let diagnoseMessage = WizardMessageDiagnose(message: "fatalErrorDiagnoseAndFeedbackBtn".localize()) else {
                        return
                    }
                    
                    diagnoseMessage.addCloseButton("fatalErrorBackBtn".localize())
                    let exitBtn = diagnoseMessage.addButton("fatalErrorExitBtn".localize()).makeDefault()
                    exitBtn.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                        NSApp.terminate(nil)
                    }
                    errorMessage.pushBefore(diagnoseMessage)
                }
                
                let exitBtn = errorMessage.addCloseButton("fatalErrorExitBtn".localize())
                exitBtn.makeDefault()
                Wizard.show(errorMessage)
            }
        }
        NSApp.terminate(nil)
    }
}
