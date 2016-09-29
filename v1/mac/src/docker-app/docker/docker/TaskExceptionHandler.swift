//
//  ExceptionHandler.swift
//  docker
//
//  Created by Emmanuel Briney on 22/12/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Foundation
import Cocoa

//Utility class to catch exception (need objective C)
public class ExceptionHandler: NSObject {
    var taskManager: TaskManager?

    public func captureException(exception: NSException) {
        self.captureException(exception)
    }
    
    public func captureUncaughtException(exception: NSException) {
        self.captureException(exception)
    }
    
    public func captureException(exception: NSException, additionalExtra: [String: AnyObject], additionalTags: [String: String]) {
        //let _ = "\(exception.name): \(exception.reason!)"
        //Logger.Record("Exception", channel: "stderr", value: message)
        //Logger.Record("Exception", channel: "stderr", value: String(format: "type:%@", exception.name))
        //Logger.Record("Exception", channel: "stderr", value: String(format: "value:%@", exception.reason!))
        
        let callStack = exception.callStackSymbols
        
        var stacktrace = [[String:String]]()
        
        if !callStack.isEmpty
        {
            for call in callStack
            {
                stacktrace.append(["function": call ])
                //Logger.Record("Exception", channel: "stderr", value: String(format: "function:%@", call))
            }
        }
        
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.quit()
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
    
    public func captureException(exception: NSException, method: String? = __FUNCTION__, file: String? = __FILE__, line: Int = __LINE__) {
        //let _ = "\(exception.name): \(exception.reason!)"
        //Logger.Record("Exception", channel: "stderr", value: message)
        //Logger.Record("Exception", channel: "stderr", value: String(format: "type:%@", exception.name))
        //Logger.Record("Exception", channel: "stderr", value: String(format: "value:%@", exception.reason!))
        
        
        var stacktrace = [[String:AnyObject]]()
        
        if let method = method, file = file {
            if line > 0 {
                var frame = [String: AnyObject]()
                frame = ["filename" : file, "function" : method, "lineno" : line]
                stacktrace = [frame]
                //Logger.Record("Exception", channel: "stderr", value: String(format: "filename:%@ - function:%@ - lineno:%@", file!, method!, line))
            }
        }
        
        let callStack = exception.callStackSymbols
        
        for call in callStack {
            stacktrace.append(["function": call ])
            //Logger.Record("Exception", channel: "stderr", value: String(format: "function:%@", call))
        }
        
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.quit()
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
    
    public func captureSignal(s: Int32) {
        //Logger.Record("Exception", channel: "stderr", value: String(format: "Signal %d catched", s))
        if (s != 9) {
            if let appDelegate = NSApp.delegate as? AppDelegate {
                appDelegate.quit()
            } else {
                Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
            }
        } else {
            // custom code path for SIGKILL
            taskManager?.Interrupt(s)
            WatchDog.Stop()
            //Logger.Shutdown()
            exit(s)
        }
    }
    
    public func setupExceptionHandler() {
        UncaughtExceptionHandler.registerHandler(self)
        NSSetUncaughtExceptionHandler(exceptionHandlerPtr)
        signal(SIGABRT, signalHandler)
        signal(SIGILL, signalHandler)
        signal(SIGSEGV, signalHandler)
        signal(SIGFPE, signalHandler)
        signal(SIGBUS, signalHandler)
        signal(SIGPIPE, signalHandler)
        signal(SIGTERM, signalHandler)
        signal(SIGKILL, signalHandler)
        signal(SIGINT, signalHandler)
    }
}
