//
//  log.swift
//  nettest
//
//  Created by Magnus Skjegstad on 02/01/2016.
//  Copyright © 2016 Magnus Skjegstad. All rights reserved.
//

import Foundation

// Wrap useful ASL functions. Requires bridging header with "asl.h" imported

enum ASLLogLevel : Int32 {
    case EMERG = 0
    case ALERT = 1
    case CRIT = 2
    case ERR = 3
    case WARNING = 4
    case NOTICE = 5
    case INFO = 6
    case DEBUG = 7
}

class ASLLog {
    let bundle : String
    private let asl : asl_object_t
    private let log_msg : asl_object_t
    private let component : String
    private let queue : dispatch_queue_t // maintain output order in serial queue
    
    init(_ bundle : String, _ component : String) {
        self.queue =  dispatch_queue_create(bundle + ".sync_queue", DISPATCH_QUEUE_SERIAL)
        self.bundle = bundle
        self.component = component
        asl = asl_open(bundle, component, UInt32(ASL_OPT_STDERR))
        log_msg = asl_new(UInt32(ASL_TYPE_MSG))
        asl_set(log_msg, ASL_KEY_SENDER, bundle)
    }
    
    // Note that args must be CStrings, so a String must e.g. be encoded as:
    // s.cStringUsingEncoding(NSUTF8StringEncoding)
    func log(level : ASLLogLevel, _ format : String, _ args: CVarArgType...) {
        // TODO We can't catch failures as we don't wait for the result
        dispatch_async(self.queue) {
            withVaList(args) {
                asl_vlog(self.asl, self.log_msg, level.rawValue, format, $0)
            }
        }
    }
    
    func info(msg : String) {
        self.log(ASLLogLevel.INFO, msg)
    }
    
    func err(msg : String) {
        self.log(ASLLogLevel.ERR, msg)
    }
    
    func warning(msg : String) {
        self.log(ASLLogLevel.WARNING, msg)
    }
    
    func debug(msg : String) {
        self.log(ASLLogLevel.DEBUG, msg)
    }
    
    func flush() {
        // Add empty block to queue and wait for it to execute
        dispatch_sync(self.queue) {
            return // wait for last entry in queue to finish
        }
    }
    
    deinit {
        flush()
        asl_release(asl)
    }
}