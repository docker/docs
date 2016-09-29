//
//  packets.swift
//  nettest
//
//  Created by Magnus Skjegstad on 28/12/2015.
//  Copyright © 2015 Magnus Skjegstad. All rights reserved.
//

import Foundation

// Packet parsing/creation helpers for testing. Not optimized for speed.
// Currently supports basic ETHERNET, ARP, IP, UDP and DHCP packets.

protocol Packet {
    var data : [UInt8] { get }
}

func uint16_to_uint8array (i : UInt16) -> [UInt8] {
    var j = i.bigEndian // by value
    let t = withUnsafePointer(&j) {
        Array(UnsafeBufferPointer(start: UnsafePointer<UInt8>($0), count: sizeof(UInt16)))
    }
    return [UInt8](t) // copy
}
func uint32_to_uint8array (i : UInt32) -> [UInt8] {
    var j = i.bigEndian // by value
    let t = withUnsafePointer(&j) {
        Array(UnsafeBufferPointer(start: UnsafePointer<UInt8>($0), count: sizeof(UInt32)))
    }
    return [UInt8](t) // copy
}

enum dhcp_op_t : UInt8 {
    case DHCP_OP_REQUEST = 1
    case DHCP_OP_REPLY = 2
}

enum dhcp_flags_t : UInt16 {
    case DHCP_FLAG_BROADCAST = 0b1000000000000000
    case DHCP_FLAG_NONE = 0
}

struct dhcp_option_t {
    let code : UInt8
    let option_data : [UInt8]
}

typealias IPv4_address_t = (UInt8, UInt8, UInt8, UInt8)
struct IPv4 : Equatable {
    let address : IPv4_address_t
    let asBytes : [UInt8]
    let asString : String

    init(_ address : IPv4_address_t) {
        self.address = address
        self.asBytes = [address.0, address.1, address.2, address.3]
        self.asString = String(address.0) + "." + String(address.1) + "." + String(address.2) + "." + String(address.3)
    }
    init?(_ bytes : [UInt8]) {
        if bytes.count != 4 {
            return nil
        }
        self.init((bytes[0], bytes[1], bytes[2], bytes[3]))
    }
    init?(_ address : String) {
        let comps = address.componentsSeparatedByString(".")
        if comps.count != 4 {
            return nil // invalid ip
        }
        var out = [UInt8](count: comps.count, repeatedValue: 0)
        for (i, v) in comps.enumerate() {
            if let n = UInt8.init(v, radix: 10) {
                out[i] = n
            } else {
                return nil
            }
        }
        self.init((out[0], out[1], out[2], out[3]))
    }
}
func ==(lhs: IPv4, rhs: IPv4) -> Bool {
    return lhs.address.0 == rhs.address.0 &&
        lhs.address.1 == rhs.address.1 &&
        lhs.address.2 == rhs.address.2 &&
        lhs.address.3 == rhs.address.3
}

typealias MAC_address_t = (UInt8, UInt8, UInt8, UInt8, UInt8, UInt8)
struct MAC {
    let address : MAC_address_t
    let asBytes : [UInt8]
    let asString : String
    init(_ address : MAC_address_t) {
        self.address = address
        self.asString = String(address.0, radix: 16) + ":" + String(address.1, radix: 16) + ":" + String(address.2, radix:16) + ":" + String(address.3,radix:16) + ":" + String(address.4,radix:16) + ":" + String(address.5,radix:16)
        self.asBytes = [address.0, address.1, address.2, address.3, address.4, address.5]
    }
    init?(_ bytes : [UInt8]) {
        if bytes.count != 6 {
            return nil
        }
        self.init((bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5]))
    }
    init?(_ address : String) {
        let comps = address.componentsSeparatedByString(":")
        if comps.count != 6 {
            return nil // invalid MAC
        }
        var out = [UInt8](count: comps.count, repeatedValue: 0)
        for (i, v) in comps.enumerate() {
            if let n = UInt8.init(v, radix: 16) {
                out[i] = n
            } else {
                return nil
            }
        }
        self.init((out[0], out[1], out[2], out[3], out[4], out[5]))
    }
}

struct dhcp_pkt : Packet {
    let op : dhcp_op_t
    let htype : UInt8 = 1 // Ethernet
    let hlen : UInt8 = 6 // Ethernet MAC
    let hops : UInt8 = 0
    let xid : UInt32 // Transaction ID
    let secs : UInt16
    let flags : dhcp_flags_t
    let ciaddr : IPv4
    let yiaddr : IPv4
    let siaddr : IPv4
    let giaddr : IPv4
    let chaddr : MAC
    let sname : [UInt8] // 64 bytes
    let file : [UInt8] // 128 bytes
    let magic : [UInt8] = [99,130,83,99]
    let options : [dhcp_option_t] // Varies, see format
    var data : [UInt8] {
        get {
            var options_len = 0
            for i in 0..<options.count {
                if options[i].code == 0 || options[i].code == 255 {
                    options_len += 1
                    continue
                }
                options_len += options[i].option_data.count + 2
            }
            
            var p : [UInt8] = []
            p.reserveCapacity(230 + options_len)
            if p.capacity < 300 {
                p.reserveCapacity(300)
            }
            p.append(op.rawValue)
            p.append(htype)
            p.append(hlen)
            p.append(hops)
            p.appendContentsOf(uint32_to_uint8array(xid))
            p.appendContentsOf(uint16_to_uint8array(secs))
            p.appendContentsOf(uint16_to_uint8array(flags.rawValue))
            p.appendContentsOf(ciaddr.asBytes)
            p.appendContentsOf(yiaddr.asBytes)
            p.appendContentsOf(siaddr.asBytes)
            p.appendContentsOf(giaddr.asBytes)
            p.appendContentsOf(chaddr.asBytes)
            p.appendContentsOf([UInt8].init(count: 16 - chaddr.asBytes.count, repeatedValue: 0))
            p.appendContentsOf(sname)
            p.appendContentsOf([UInt8].init(count: 64-sname.count, repeatedValue: 0))
            p.appendContentsOf(file)
            p.appendContentsOf([UInt8].init(count: 128-file.count, repeatedValue: 0))
            p.appendContentsOf(magic)
            
            // add options
            for i in 0..<options.count {
                p.append(options[i].code)
                if options[i].code != 0 && options[i].code != 255 { // skip special codes
                    p.append(UInt8(options[i].option_data.count))
                    p.appendContentsOf(options[i].option_data)
                }
            }
    
            // pad until at least 300 bytes
            if p.count < 300 {
                p.appendContentsOf([UInt8].init(count: 300-p.count, repeatedValue: 0))
            }
            
            return p
        }
    }
    
}

struct bytes_pkt : Packet {
    let bytes : [UInt8]
    var data : [UInt8] {
        return bytes
    }
    init(_ bytes : [UInt8]) {
        self.bytes = bytes
    }
}

struct udp_pkt : Packet {
    let src_port : UInt16
    let dst_port : UInt16
    let payload : Packet
    var data : [UInt8] {
        get {
            var p : [UInt8] = []
            p.reserveCapacity(0xffff)
            let d = payload.data
            p.appendContentsOf(uint16_to_uint8array(src_port))
            p.appendContentsOf(uint16_to_uint8array(dst_port))
            p.appendContentsOf(uint16_to_uint8array(UInt16(d.count) + 8)) // add udp header len
            p.appendContentsOf(uint16_to_uint8array(0x0000)) // checksum not supported
            p.appendContentsOf(d)
            return p
        }
    }
    init(src_port : UInt16, dst_port : UInt16, payload : Packet) {
        self.src_port = src_port
        self.dst_port = dst_port
        self.payload = payload
    }
    init?(bytes: [UInt8]) {
        if bytes.count < 8 {
            return nil // header too short
        }
        self.src_port = UInt16(bytes[0]) << 8 + UInt16(bytes[1])
        self.dst_port = UInt16(bytes[2]) << 8 + UInt16(bytes[3])
        //let len = UInt16(bytes[4]) << 8 + UInt16(bytes[5])
        //let checksum = UInt16(bytes[6]) << 8 + UInt16(bytes[7])
        self.payload = bytes_pkt(Array(bytes[8..<bytes.count]))
    }
}

enum ip_proto : UInt8 {
    case UDP = 17
}

struct ip_pkt : Packet {
    let header : UInt8
    let diffserv : UInt8
    let ident : UInt16
    let flagsfrag : UInt16
    let ttl : UInt8
    let proto : ip_proto
    let src : IPv4
    let dst : IPv4
    let payload : Packet
    var data : [UInt8] {
        get {
            var p : [UInt8] = []
            p.reserveCapacity(0xffff)
            let d = payload.data
            p.append(header)
            p.append(diffserv)
            p.appendContentsOf(uint16_to_uint8array(ident))
            p.appendContentsOf(uint16_to_uint8array(flagsfrag))
            p.append(ttl)
            p.append(proto.rawValue)
            //p.appendContentsOf(uint16_to_uint8array(0x0)) // checksum disabled
            p.appendContentsOf(src.asBytes)
            p.appendContentsOf(dst.asBytes)
            p.insertContentsOf(uint16_to_uint8array(UInt16(p.count + d.count + 2 + 2)), at: 2) // insert length
            let header_checksum = checksum16(&p)
            p.insertContentsOf(uint16_to_uint8array(header_checksum), at: 10) // add checksum
            p.appendContentsOf(d)
            
            return p
        }
    }
    init(ttl : UInt8, proto : ip_proto, src : IPv4, dst : IPv4, payload : Packet) {
        self.header = 0x45
        self.diffserv = 0
        self.ident = 0
        self.flagsfrag = 0
        self.proto = proto
        self.ttl = ttl
        self.src = src
        self.dst = dst
        self.payload = payload
    }
    init?(bytes : [UInt8]) {
        if bytes.count < 20 {
            return nil
        }
        self.header = bytes[0]
        self.diffserv = bytes[1]
        //let len = UInt16(bytes[2]) << 8 + UInt16(bytes[3])
        self.ident = UInt16(bytes[4]) << 8 + UInt16(bytes[5])
        self.flagsfrag = UInt16(bytes[6]) << 8 + UInt16(bytes[7])
        self.ttl = bytes[8]
        if let proto = ip_proto.init(rawValue: bytes[9]) {
            self.proto = proto
        } else {
            return nil // proto unknown
        }
        //let checksum = UInt16(bytes[10]) << 8 + UInt16(bytes[11])
        if let src = IPv4(Array(bytes[12..<16])), dst = IPv4(Array(bytes[16..<20])){
            self.dst = dst
            self.src = src
        } else {
            return nil // unable to parse src/dst
        }
        switch self.proto {
        case .UDP:
            if let p = udp_pkt.init(bytes: Array(bytes[20..<bytes.count])) {
                self.payload = p
            } else {
                return nil // unable to parse udp data
            }
        }
    }
}

struct ether_frame : Packet {
    let dst : MAC
    let src : MAC
    let ethertype : UInt16
    let payload  : Packet
    var data : [UInt8] {
        get {
            var p : [UInt8] = []
            p.appendContentsOf(dst.asBytes)
            p.appendContentsOf(src.asBytes)
            p.appendContentsOf(uint16_to_uint8array(ethertype))
            p.appendContentsOf(payload.data)
            return p
        }
    }
    init(dst : MAC, src : MAC, ethertype : UInt16, payload : Packet) {
        self.dst = dst
        self.src = src
        self.ethertype = ethertype
        self.payload = payload
    }
    init?(bytes : [UInt8]) {
        if let dst = MAC(Array(bytes[0..<6])), src = MAC(Array(bytes[6..<12])){
            self.dst = dst
            self.src = src
        } else {
            print("unable to read mac src/dst")
            return nil
        }
        self.ethertype = UInt16(bytes[12]) << 8 + UInt16(bytes[13])
        switch self.ethertype {
        case 0x0800:
            if let ip = ip_pkt.init(bytes: Array(bytes[14..<bytes.count])) {
                self.payload = ip
            } else {
                print("couldn't parse ip packet")
                return nil // couldn't parse IP packet
            }
        case 0x0806:
            if let arp = arp_pkt.init(bytes: Array(bytes[14..<bytes.count])) {
                self.payload = arp
            } else {
                print("couldn't parse arp packet")
                return nil
            }
        default:
            print("unsupported ethertype",String(ethertype, radix:16))
            return nil // unsupported ethertype
        }
    }
}

enum arp_op_t : UInt16 {
    case ARP_REQUEST = 1
    case ARP_REPLY = 2
    case RARP_REQUEST = 3
    case RARP_REPLY = 4
    case DRARP_REQUEST = 5
    case DRARP_REPLY = 6
    case DRARP_ERROR = 7
    case INARP_REQUEST = 8
    case INARP_REPLY = 9
    // incomplete
}

struct arp_pkt : Packet {
    let hrd : UInt16
    let pro : UInt16
    let hln : UInt8
    let pln : UInt8
    let op : arp_op_t
    let sha : MAC
    let spa : IPv4
    let tha : MAC
    let tpa : IPv4
    
    var data : [UInt8] {
        get {
            var p : [UInt8] = []
            p.reserveCapacity(0xff)
            p.appendContentsOf(uint16_to_uint8array(hrd))
            p.appendContentsOf(uint16_to_uint8array(pro))
            p.append(hln)
            p.append(pln)
            p.appendContentsOf(uint16_to_uint8array(op.rawValue))
            p.appendContentsOf(sha.asBytes)
            p.appendContentsOf(spa.asBytes)
            p.appendContentsOf(tha.asBytes)
            p.appendContentsOf(tpa.asBytes)
            return p
        }
    }
    
    init(op : arp_op_t, sha : MAC, spa : IPv4, tha : MAC, tpa : IPv4) {
        hrd = 6 // IEEE 802.x
        pro = 0x0800 // IPv4
        hln = 6 // MAC len
        pln = 4 // IP len
        self.op = op
        self.sha = sha
        self.spa = spa
        self.tha = tha
        self.tpa = tpa
    }
    
    init?(bytes : [UInt8]) {
        self.hrd = UInt16(bytes[0]) << 8 + UInt16(bytes[1])
        self.pro = UInt16(bytes[2]) << 8 + UInt16(bytes[3])
        self.hln = bytes[4]
        self.pln = bytes[5]
        if let op = arp_op_t(rawValue: UInt16(bytes[6]) << 8 + UInt16(bytes[7])) {
            self.op = op
        } else {
            print("unknown arp op")
            return nil // unknown op
        }
        if let sha = MAC(Array(bytes[8..<14])),
            spa = IPv4(Array(bytes[14..<18])),
            tha = MAC(Array(bytes[18..<24])),
            tpa = IPv4(Array(bytes[24..<28])) {
                self.sha = sha
                self.spa = spa
                self.tha = tha
                self.tpa = tpa
        } else {
            print("unable to parse sha, spa, tha or tpa")
            return nil
        }
    }
    
}