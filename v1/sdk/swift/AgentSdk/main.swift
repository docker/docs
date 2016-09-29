//
//  main.swift
//  AgentSdk
//
//  Created by Gaetan de Villele on 2/11/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Pinata


// MARK: *** EXPOSED ***

// true on success
public func ppInit() -> String? {
    let status: UInt32 = pp_init()
    if status != PP_OK && status != PP_ERROR_SDK_ALREADY_INITIALIZED {
        return "failed to init SDK"
    }
    return nil
}

public func ppConnect(socketPath: String) -> (error: String?, conn: UInt32) {
    // convert arguments
    guard let cSocketPath:[CChar] = socketPath.cStringUsingEncoding(NSUTF8StringEncoding) else {
        print("[\(#function)]: failed to convert socketPath")
        return ("failed to convert socketPath", 0)
    }
    var connRef: UInt32 = 0
    let status: UInt32 = pp_connect(cSocketPath, &connRef)
    if status != PP_OK {
        return ("internal error", 0)
    }
    return (nil, connRef)
}

public func ppDisconnect(conn: UInt32) -> String? {
    let status = pp_disconnect(conn)
    if status != PP_OK {
        return "failed to disconnect"
    }
    return nil
}

public func ppRequest(conn: UInt32, request: String) -> (error: String?, response: String) {
    // convert arguments
    guard let cRequest:[CChar] = request.cStringUsingEncoding(NSUTF8StringEncoding) else {
        print("[\(#function)]: failed to convert request")
        return ("failed to convert request", "")
    }
    var c_response: UnsafeMutablePointer<CChar> = nil
    let status = pp_request(conn, cRequest, &c_response)
    if status != PP_OK {
        return ("failed to request", "")
    }
    let response:String? = String.fromCString(c_response)
    // free c string
    pp_response_free(c_response)
    guard let responseValue = response else {
        return ("failed to read response", "")
    }
    return (nil, responseValue)
}



//public func setAddressAndPort(address:String, port:UInt16) -> String? {
//    let c_address:[CChar]? = address.cStringUsingEncoding(NSUTF8StringEncoding)
//    if c_address == nil {
//        return "couldn't set address"
//    }
//    if pinata_set_address_and_port(c_address!, port) == PINATA_ERROR {
//        let error:String? = String.fromCString(pinata_get_error())
//        if error != nil {
//            return error
//        } else {
//            return "couldn't get error message from sdk"
//        }
//    }
//    return nil
//}
//
////
//public func disconnect() {
//    pinata_disconnect()
//}
//
////
//public func privateAccessAuthenticateWithToken(token:String, AndUDID udid:String) -> String? {
//    // generate the request content
//    let request:String = "privateAccessAuthenticate \(udid) \(token)"
//    // try to send the request and get a response from the agent
//    let result = sendRequest(request)
//    if result.error != nil {
//        // there was an error
//        return result.error
//    }
//    if result.response.substringToIndex(result.response.startIndex.advancedBy(2)) == "OK" {
//        // success
//        return nil
//    }
//    if result.response.substringToIndex(result.response.startIndex.advancedBy(5)) == "ERROR" {
//        if result.response.characters.count > 5 {
//            // error message
//            let errorMessage: String = result.response.substringFromIndex(result.response.startIndex.advancedBy(6))
//            return errorMessage
//        } else {
//            return "unknown error"
//        }
//    }
//    return "unknown response"
//}
//
//// app management
//// container management
//// pipeline management
//
//// MARK: *** PRIVATE ***
//
//private func sendRequest(request:String) -> (response:String, error:String?) {
//    let c_request:[CChar]? = request.cStringUsingEncoding(NSUTF8StringEncoding)
//    guard let c_request_unwrapped = c_request else {
//        return ("", "request is empty")
//    }
//    let c_response:UnsafeMutablePointer<CChar> = pinata_send_request(c_request_unwrapped)
//    let response:String? = String.fromCString(c_response)
//    // String.fromCString creates a copy, so it's safe to free the
//    // response here.
//    pinata_free_response(c_response)
//    if response == nil {
//        // there is an error to retrieve
//        let error:String? = String.fromCString(pinata_get_error())
//        if error != nil {
//            return ("", error)
//        } else {
//            return ("", "couldn't get error message from sdk")
//        }
//    }
//    return (response!, nil)
//}
