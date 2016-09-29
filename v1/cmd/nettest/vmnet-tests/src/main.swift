//
//  Created by Magnus Skjegstad on 17/12/2015.
//  Copyright © 2015 Magnus Skjegstad. All rights reserved.
//

import Foundation
import vmnet

var Timestamp: NSTimeInterval {
    return NSDate().timeIntervalSince1970 * 1000
}

func fail(message : String...) {
    print("Test(s) failed:",message)
    exit(-1)
}

func newtest(name : String) {
    print("TEST:",name)
}

func connect(logger : ASLLog, _ testname : String) -> (Vmnet_if, vmnet_if_config_t, NSUUID)? {
    newtest(testname)
    let uuid = NSUUID.init()
    logger.info("Initialising Vmnet.framework with UUID \(uuid.UUIDString)...")
    let x = Vmnet_if(logger, uuid)
    if !(x.is_running()) {
        logger.err("Vmnet.framework init failed with error \(x.status_msg)")
        return nil
    }
    let config = x.config!
    logger.info("MAC is \(config.mac.asString), UUID is \(uuid.UUIDString)")
    return (x, config, uuid)
}


func make_discover_packet(config : vmnet_if_config_t) -> NSData {
    let xid = arc4random() // DHCP transaction ID
    let dhcp_discover_packet : ether_frame = ether_frame.init(
        dst: MAC((0xff,0xff,0xff,0xff,0xff,0xff)),
        src: config.mac,
        ethertype: 0x0800,
        payload: ip_pkt.init(
            ttl: 64,
            proto: ip_proto.UDP,
            src: IPv4((0x00,0x00,0x00,0x00)),
            dst: IPv4((0xff,0xff,0xff,0xff)),
            payload: udp_pkt.init(
                src_port: 68,
                dst_port: 67,
                payload: dhcp_pkt.init(
                    op: dhcp_op_t.DHCP_OP_REQUEST,
                    xid: xid,
                    secs: 3,
                    flags: dhcp_flags_t.DHCP_FLAG_BROADCAST,
                    ciaddr: IPv4((0,0,0,0)),
                    yiaddr: IPv4((0,0,0,0)),
                    siaddr: IPv4((0,0,0,0)),
                    giaddr: IPv4((0,0,0,0)),
                    chaddr: config.mac,
                    sname: [],
                    file: [],
                    options: [
                        dhcp_option_t.init(code: 53, option_data: [1]),
                        dhcp_option_t.init(code: 61, option_data: [0x01] /* ethernet */ + config.mac.asBytes),
                        dhcp_option_t.init(code: 57, option_data: uint16_to_uint8array(576)),
                        
                        //dhcp_option_t.init(code: 50, option_data: IPv4((192,168,64,123)).asBytes),
                        
                        dhcp_option_t.init(code: 55, option_data: [ // request...
                            1, // Subnet mask
                            3, // Router
                            6, // DNS
                            12, // Host name
                            15, // Domain name
                            28, // Broadcast address
                            42 // NTP server
                            ]),
                        dhcp_option_t.init(code: 12, option_data: [UInt8].init("docker-bootstrap".utf8)),
                        
                        dhcp_option_t.init(code: 255, option_data: [])
                        
                    ])
            )))
    
    let packet_data = dhcp_discover_packet.data
    // TODO This can probably be done without copying the array
    return NSData.init(bytes: packet_data, length: packet_data.count)
    
}

let bundle = "com.docker.nettest"
let logger = ASLLog(bundle, "main")


let send_q = dispatch_queue_create(bundle + ".send_queue", DISPATCH_QUEUE_CONCURRENT)

print("Testing the time it takes for a new network interface to get an IP address. If no other clients are be connected to Vmnet.framework (i.e. whale is not running), the first DHCP DISCOVER packets are always lost.\n\n")

let init_t = Timestamp

// Test regular DHCP
if let (vmnet, config, uuid) = connect(logger, "Testing regular DHCP request w/1sec retransmit") {
    let sem = dispatch_semaphore_create(0)
    let t0 = Timestamp
    // Add listener that signals semaphore on for a DHCP OFFER
    vmnet.listen_async {
        (newpkts : UInt32) -> Void in
        var buf = [UInt8].init(count: Int(config.max_packet_size), repeatedValue: 0)
        var d = NSData.init(bytesNoCopy: &buf, length: buf.count, freeWhenDone: false)
        var len : UInt32 = 0
        let res = vmnet.read_packet(&d, read: &len)
        if res == vmnet_return_t.VMNET_SUCCESS {
            if let pkt = ether_frame.init(bytes: buf) {
                switch pkt.ethertype {
                case 0x0800:
                    let ip = pkt.payload as! ip_pkt
                    let udp = ip.payload as! udp_pkt
                    logger.info("IP/UDP \(ip.src.asString):\(udp.src_port) -> \(ip.dst.asString):\(udp.dst_port)")
                    // TODO this only checks ip and port, should check DHCP data as well
                    if (ip.src == IPv4((192,168,64,1)) && udp.src_port == 67) &&
                        (ip.dst == IPv4((255,255,255,255)) && udp.dst_port == 68) {
                            dispatch_semaphore_signal(sem) // signal
                    }
                case 0x0806:
                    let arp = pkt.payload as! arp_pkt
                    logger.info("ARP OP: \(arp.op.rawValue) SHA: \(arp.sha.asString) SPA: \(arp.spa.asString) THA: \(arp.tha.asString) TPA: \(arp.tpa.asString)")
                default:
                    logger.warning("unknown ethertype \(pkt.ethertype)")
                }
                
            } else {
                logger.warning("could not parse")
            }
        } else {
            logger.warning("read failed: \(vmnet.vmnet_error_msg(res))")
        }
    }
    // Schedule packets
    var packet = make_discover_packet(config)
    var start_time = dispatch_time(dispatch_time_t(DISPATCH_TIME_NOW), 0)
    var packets_to_send = 10
    var t = start_time
    var delta = Int64(NSEC_PER_SEC) / 1
    logger.info("Scheduling \(packets_to_send) packets to be sent with \(delta/1000000) ms intervals")
    for var i = 0; i < packets_to_send; i++ {
        let c = i+1
        dispatch_after(t, send_q) {
            if vmnet.is_running() {
                logger.info("Sending DHCP_DISCOVER #\(c)")
                let s = vmnet.write_packet(&packet)
                if s != vmnet_return_t.VMNET_SUCCESS {
                    logger.info("Test 1: Send failed with error \(s.rawValue)")
                }
            }
        }
        t = dispatch_time(t, delta)
    }
    // Wait for offer 
    dispatch_semaphore_wait(sem, DISPATCH_TIME_FOREVER)
    let t1 = Timestamp
    logger.info("DHCP OFFER received after \(t1-t0) ms")
    
    vmnet.shutdown()
    logger.flush() // wait for remaining output
    
} else {
    fail("Unable to initialize vmnet")
}


// Test regular DHCP
if let (vmnet, config, uuid) = connect(logger, "Testing DHCP DISCOVER w/10ms retransmit") {
    let sem = dispatch_semaphore_create(0)
    let t0 = Timestamp
    // Add listener that signals semaphore on for a DHCP OFFER
    vmnet.listen_async {
        (newpkts : UInt32) -> Void in
        var buf = [UInt8].init(count: Int(config.max_packet_size), repeatedValue: 0)
        var d = NSData.init(bytesNoCopy: &buf, length: buf.count, freeWhenDone: false)
        var len : UInt32 = 0
        let res = vmnet.read_packet(&d, read: &len)
        if res == vmnet_return_t.VMNET_SUCCESS {
            if let pkt = ether_frame.init(bytes: buf) {
                switch pkt.ethertype {
                case 0x0800:
                    let ip = pkt.payload as! ip_pkt
                    let udp = ip.payload as! udp_pkt
                    logger.info("IP/UDP \(ip.src.asString):\(udp.src_port) -> \(ip.dst.asString):\(udp.dst_port)")
                    // TODO this only checks ip and port, should check DHCP data as well
                    if (ip.src == IPv4((192,168,64,1)) && udp.src_port == 67) &&
                        (ip.dst == IPv4((255,255,255,255)) && udp.dst_port == 68) {
                            dispatch_semaphore_signal(sem) // signal
                    }
                case 0x0806:
                    let arp = pkt.payload as! arp_pkt
                    logger.info("ARP OP: \(arp.op.rawValue) SHA: \(arp.sha.asString) SPA: \(arp.spa.asString) THA: \(arp.tha.asString) TPA: \(arp.tpa.asString)")
                default:
                    logger.warning("unknown ethertype \(pkt.ethertype)")
                }
                
            } else {
                logger.warning("could not parse")
            }
        } else {
            logger.warning("read failed: \(vmnet.vmnet_error_msg(res))")
        }
    }
    // Schedule packets
    var packet = make_discover_packet(config)
    var start_time = dispatch_time(dispatch_time_t(DISPATCH_TIME_NOW), 0)
    var packets_to_send = 10
    var t = start_time
    var delta = Int64(NSEC_PER_SEC) / 100
    logger.info("Scheduling \(packets_to_send) packets to be sent with \(delta/1000000) ms intervals")
    for var i = 0; i < packets_to_send; i++ {
        let c = i+1
        dispatch_after(t, send_q) {
            if vmnet.is_running() {
                logger.info("Sending DHCP_DISCOVER #\(c)")
                let s = vmnet.write_packet(&packet)
                if s != vmnet_return_t.VMNET_SUCCESS {
                    logger.info("Test 2: Send failed with error \(s.rawValue)")
                }
            }
        }
        t = dispatch_time(t, delta)
    }
    // Wait for offer
    dispatch_semaphore_wait(sem, DISPATCH_TIME_FOREVER)
    let t1 = Timestamp
    logger.info("DHCP OFFER received after \(t1-t0) ms")
    
    vmnet.shutdown()
    logger.flush() // wait for remaining output
} else {
    fail("Unable to initialize vmnet")
}
