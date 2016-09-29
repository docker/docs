//
//  vmnet.swift
//  nettest
//
//  Created by Magnus Skjegstad on 28/12/2015.
//  Copyright © 2015 Magnus Skjegstad. All rights reserved.
//

import Foundation
import vmnet

public struct vmnet_if_config_t {
    let mac : MAC
    let iface_params : xpc_object_t
    let max_packet_size : uint64
    let mtu : uint64
}

public enum vmnet_if_status_t : Equatable {
    case Uninitialized ()
    case Vmnet_error(vmnet_return_t)
    case Success(vmnet_if_config_t)
}
public func ==(s1 : vmnet_if_status_t, s2: vmnet_if_status_t) -> Bool {
    switch (s1, s2) {
    case (.Success(let a), .Success(let b)):
        return a.mac.asString == b.mac.asString // compare mac to see if same if
    case (.Vmnet_error(let a), .Vmnet_error(let b)):
        return a == b
    case (.Uninitialized(), .Uninitialized()):
        return true
    default:
        return false
    }
}

@objc public class Vmnet_if : NSObject {
    public let vmnet_uuid : NSUUID
    private let logger : ASLLog
    private let iface_descr : xpc_object_t = xpc_dictionary_create(nil, nil, 0)
    private let task_sq : dispatch_queue_t
    
    private var iface_ref : interface_ref?
    private var event_counter : Int = 0
    private var last_event : Int = -1
    
    private var _status : vmnet_if_status_t = vmnet_if_status_t.Uninitialized() // Should only be accessed with dispatch_sync through task_sq
    
    public var status : vmnet_if_status_t {
        get {
            return self.get_status()
        }
    }
    
    public var status_msg : String {
        get {
            return self.get_status_msg(self.get_status())
        }
    }
    
    public var config : vmnet_if_config_t? {
        get {
            switch self.get_status() {
            case .Success(let config):
                return config
            default:
                return nil
            }
            
        }
    }
    
    public func is_running() -> Bool {
        switch self.status {
        case .Success(_):
            return true
        default:
            return false
        }
    }
    
    func shutdown() {
        dispatch_sync(self.task_sq) {
            if case .Success(_) = self._status {
                if let iface = self.iface_ref {
                    self._status = vmnet_if_status_t.Uninitialized()
                    vmnet_stop_interface(iface, self.task_sq) {_ in                        
                        self.logger.info("Vmnet interface shut down (uuid=\(self.vmnet_uuid.UUIDString))")
                    }
                }
            }
        }
    }
    
    deinit {
        shutdown()
    }
    
    init(_ logger : ASLLog, _ uuid : NSUUID?) {
        self.logger = logger
        self.vmnet_uuid = uuid ?? NSUUID.init()
        self.task_sq = dispatch_queue_create(logger.bundle + ".tasks_sq", DISPATCH_QUEUE_SERIAL)
        
        super.init()
        
        let sem : dispatch_semaphore_t = dispatch_semaphore_create(0)
        
        var uuid_b = [UInt8](count: 16, repeatedValue: 0)
        vmnet_uuid.getUUIDBytes(&uuid_b)
        xpc_dictionary_set_uuid(iface_descr, vmnet_interface_id_key, uuid_b)
        xpc_dictionary_set_uint64(iface_descr, vmnet_operation_mode_key, (uint64)(operating_modes_t.VMNET_SHARED_MODE.rawValue))
        
        self.iface_ref = vmnet_start_interface(iface_descr, task_sq) {
            (status : vmnet_return_t, iface_params : xpc_object_t?) -> Void in
            switch status {
            case .VMNET_SUCCESS:
                if let iface_params = iface_params, mac_s = String(UTF8String : xpc_dictionary_get_string(iface_params, vmnet_mac_address_key))
                {
                    if let mac = MAC(mac_s) {
                        let mtu = xpc_dictionary_get_uint64(iface_params, vmnet_mtu_key);
                        let max_packet_size = xpc_dictionary_get_uint64(iface_params, vmnet_max_packet_size_key);
                        self._status = vmnet_if_status_t.Success(
                            vmnet_if_config_t(
                                mac:mac,
                                iface_params:iface_params,
                                max_packet_size:max_packet_size,
                                mtu:mtu)
                        )
                    } else {
                        fallthrough
                    }
                } else {
                    fallthrough // Return Error if we couldn't get all parameters
                }
            default: // Return error type for all other status codes
                self._status = vmnet_if_status_t.Vmnet_error(status)
            }
            
            dispatch_semaphore_signal(sem)
        }
        
        if dispatch_semaphore_wait(sem, dispatch_time(DISPATCH_TIME_NOW, 10000000000)) != 0 {
            self._status = vmnet_if_status_t.Uninitialized()
        }
        
        
    }
    
    func listen_async(packets_available_handler : (UInt32 -> Void)) {
        if let iface_ref = self.iface_ref {
            switch self.get_status() {
            case .Success(_):
                vmnet_interface_set_event_callback(iface_ref,
                    interface_event_t.VMNET_INTERFACE_PACKETS_AVAILABLE,
                    task_sq) {
                        (event_type : interface_event_t, params : xpc_object_t?) -> Void in
                        switch event_type {
                        case interface_event_t.VMNET_INTERFACE_PACKETS_AVAILABLE:
                            self.event_counter++ // no lock here, as we only read/write from same serial dispatch queue
                            //print("packet: ", self.event_counter, event_type.rawValue)
                            dispatch_async(self.task_sq) {
                                if let params = params {
                                    let packets = UInt32(xpc_dictionary_get_uint64(params, vmnet_estimated_packets_available_key))
                                    return packets_available_handler(packets)
                                } else {
                                    return packets_available_handler(1)
                                }
                            }
                        default:
                            self.logger.warning("unknown event type \(event_type.rawValue)")
                        }
                }
            case .Vmnet_error(_):
                break
            case .Uninitialized:
                break
            }
        }
    }
    
    func vmnet_error_msg(e : vmnet_return_t) -> String {
        switch e {
        case .VMNET_SUCCESS:
            return "Success"
        case .VMNET_BUFFER_EXHAUSTED:
            return "Buffer exhausted"
        case .VMNET_FAILURE:
            return "Failure"
        case .VMNET_INVALID_ACCESS:
            return "Invalid access"
        case .VMNET_INVALID_ARGUMENT:
            return "Invalid argument"
        case .VMNET_MEM_FAILURE:
            return "Mem failure"
        case .VMNET_PACKET_TOO_BIG:
            return "Packet too big"
        case .VMNET_SETUP_INCOMPLETE:
            return "Setup incomplete"
        case .VMNET_TOO_MANY_PACKETS:
            return "Too many packets"
        }
    }
    
    func get_status_msg(status : vmnet_if_status_t) -> String {
        switch status {
        case .Vmnet_error(let e):
            return vmnet_error_msg(e)
        case .Uninitialized():
            return "Uninitialized"
        case .Success:
            return "Success"
        }
    }
    
    public func write_packet(inout buffer : NSData) -> vmnet_return_t {
        if let iface_ref = self.iface_ref {
            var iov : iovec = iovec(
                iov_base : UnsafeMutablePointer<Void>(buffer.bytes),
                iov_len : buffer.length)
            var v : vmpktdesc = vmpktdesc(
                vm_pkt_size : buffer.length,
                vm_pkt_iov : &iov,
                vm_pkt_iovcnt : 1,
                vm_flags : 0
            )
            var pktcnt : Int32 = 1
            return vmnet_write(iface_ref, &v, &pktcnt)
        } else {
            return vmnet_return_t.VMNET_FAILURE
        }
    }
    
    public func read_packet(inout buffer : NSData, inout read : UInt32) -> vmnet_return_t {
        if let iface_ref = self.iface_ref {
            var iov : iovec = iovec(
                iov_base : UnsafeMutablePointer<Void>(buffer.bytes),
                iov_len : buffer.length)
            var v : vmpktdesc = vmpktdesc(
                vm_pkt_size : buffer.length,
                vm_pkt_iov : &iov,
                vm_pkt_iovcnt : 1,
                vm_flags : 0
            )
            var pktcnt : Int32 = 1
            let res = vmnet_read(iface_ref, &v, &pktcnt)
            if res == vmnet_return_t.VMNET_SUCCESS {
                read = UInt32(v.vm_pkt_size)
            } else {
                read = 0
            }
            return res
        } else {
            return vmnet_return_t.VMNET_FAILURE
        }
        
    }
    
    private func get_status() -> vmnet_if_status_t {
        var v : vmnet_if_status_t?
        dispatch_sync(self.task_sq) {
            v = self._status
        }
        return v! // if v is nil, this is an internal error
    }
    
}
