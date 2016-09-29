//
//  vmnetd.swift
//  docker-installer
//
//  Created by David Scott on 07/12/2015.
//  Copyright © 2015 David Scott. All rights reserved.
//

import Foundation

// If vmnetd is not running this protocol version, reinstall it.
// NOTE: must match v1/cmd/com.docker.vmnetd/src/protocol.h
let required_version = 18

public enum VmnetdState {
    case Missing // next step: install
    case Ancient // next step: ask the user to manually uninstall and rerun the installer
    case NeedsUpgrade // next step: automatic upgrade
    case JustRight // leave it alone
}

let vmnetd_magic = "VMN3T"
let vmnetd_version: [UInt8] = [ 1, 0, 0, 0 ] // This version always works for probing

let command_uninstall: [UInt8] = [ 2 ]
let command_symlink_install: [UInt8] = [ 3 ]
let command_symlink_uninstall: [UInt8] = [ 4 ]
let command_uninstall_sockets: [UInt8] = [5]
let command_symlink_install2: [UInt8] = [ 7 ]


func really_write_array(s: Int32, bytes: [UInt8])-> Bool {
    var remaining = bytes.count
    var data: [UInt8] = bytes
    while remaining > 0 {
        let n = Darwin.write(s, data, remaining)
        if n <= 0 {
            print("short write through socket")
            return false
        }
        remaining = remaining - n
        for _ in 1...n {
            data.removeAtIndex(0)
        }
    }
    return true
}

func really_write_string(s: Int32, data: String)-> Bool {
    let bytes = [UInt8](data.utf8)
    return really_write_array(s, bytes: bytes)
}

func really_read_array(s: Int32, count: Int) -> UnsafeMutablePointer<CChar>? {
    let buffer = UnsafeMutablePointer<CChar>.alloc(count+1)
    var remaining = count
    var offset = 0
    while remaining > 0 {
        let n = Darwin.read(s, &buffer[offset], remaining)
        if n <= 0 {
            print("short read through socket")
            return nil
        }
        remaining = remaining - n
        offset = offset + n
    }
    return buffer
}

func really_read_string(s: Int32, count: Int) -> String? {
    let buffer = UnsafeMutablePointer<CChar>.alloc(count+1)
    var remaining = count
    var offset = 0
    while remaining > 0 {
        let n = Darwin.read(s, &buffer[offset], remaining)
        if n <= 0 {
            print("short read through socket")
            return nil
        }
        remaining = remaining - n
        offset = offset + n
    }
    buffer[count] = 0
    return String.fromCString(buffer)
}

func one_connect_attempt(socket_path: String) -> Int32? {
    let s = Darwin.socket(AF_UNIX, SOCK_STREAM, 0)
    var addr = sockaddr_un()
    var buf = [UInt8]()
    buf += socket_path.utf8
    addr.sun_family = sa_family_t(AF_UNIX)
    memcpy(&addr.sun_path, &buf, buf.count)
    let len = socklen_t(sizeof(sockaddr_un))
    let err = withUnsafeMutablePointer(&addr) {
        Darwin.connect(s, UnsafeMutablePointer($0), len)
    }
    if (err != 0) {
        close(s)
        return nil
    }
    return s
}

func connect(socket_path: String) -> Int32? {
    // When a launchd service is registered, the registration is not
    // synchronous; we might try to use it before it's ready. Try up
    // to 5 times with 1s delay between each try.
    var attempts = 0
    while attempts < 5 {
        attempts = attempts + 1
        let s = one_connect_attempt(socket_path)
        if s != nil {
            return s
        }
        Logger.log(level: Logger.Level.Notice, content: "com.docker.vmnetd failed to connect to \(socket_path), sleeping 1s")
        sleep(1)
    }
    return nil
}

public func getVmnetdState(socket_path: String) -> VmnetdState {
    if let s = connect(socket_path) {
        // Transmit v1 init_message
        if !really_write_string(s, data: vmnetd_magic) {
            close(s)
            return VmnetdState.Missing
        }
        // Ancient v0 vmnetd will close the connection as soon as it
        // fails to recognise the new magic.
        if !really_write_array(s, bytes: vmnetd_version) {
            close(s)
            return VmnetdState.Ancient
        }
        if !really_write_string(s, data: git_commit) {
            close(s)
            return VmnetdState.Ancient
        }
        let hello = really_read_string(s, count: 5)
        if hello == nil {
            close(s)
            return VmnetdState.Ancient
        }
        let bytes = really_read_array(s, count: 4)
        var version = 0
        if let buf = bytes {
            version = (Int(buf[0]) << 0)
            version = version | (Int(buf[1]) << 8)
            version = version | (Int(buf[2]) << 16)
            version = version | (Int(buf[3]) << 24)

        }
        if let commit = really_read_string(s, count: 40) {
            close(s)
            
            Logger.log(level: Logger.Level.Notice, content: "com.docker.docker commit: \(git_commit) wants version \(required_version)")
            Logger.log(level: Logger.Level.Notice, content: "com.docker.vmnetd commit: \(commit) has version \(version)")
            if required_version == version || (git_commit == commit) {
                return VmnetdState.JustRight
            }
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "com.docker.vmnetd: failed to read commit in getVmnetdState")
        }
        return VmnetdState.NeedsUpgrade
    } else {
        return VmnetdState.Missing
    }
}

public func handshake(socket_path: String) -> Int32? {
    if let s = connect(socket_path) {
        if !really_write_string(s, data: vmnetd_magic) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to write magic")
            return nil
        }
        if !really_write_array(s, bytes: vmnetd_version) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to write version")
            return nil
        }
        if !really_write_string(s, data: git_commit) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to write commit")
            return nil
        }
        let hello = really_read_string(s, count: 5)
        if hello == nil {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to read magic")
            return nil
        }
        let version = really_read_string(s, count: 4)
        if version == nil {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to read version")
            return nil
        }
        let commit = really_read_string(s, count: 40)
        if commit == nil {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to read commit")
            return nil
        }
        return s
    }
    return nil
}

public func probe_vmnetd(socket_path: String) -> Bool {
    if let s = handshake(socket_path) {
        close(s)
        return true
    }
    return false
}

public func uninstall(socket_path: String) -> Bool {
    if let s = handshake(socket_path) {
        if !really_write_array(s, bytes: command_uninstall) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd uninstall: failed to write uninstall command")
            return false
        }
        let bytes = really_read_array(s, count: 1)
        if let result = bytes {
            if result[0] != 0 {
                Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd uninstall: server reports failure")
                return false
            }
            return true
        }
        Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd uninstall: failed to read uninstall result")
        return false
    } else {
        return false
    }
}

public func install_symlinks(socket_path: String, container_folder: String) -> Bool {
    if let s = handshake(socket_path) {
        if !really_write_array(s, bytes: command_symlink_install2) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: failed to write install symlinks2 command")
            return false
        }
        // Note we pad the string to 1024 bytes (PATH_MAX on OSX)
        if !really_write_string(s, data: container_folder.stringByPaddingToLength(1024, withString: "\000", startingAtIndex: 0)) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd: failed to write container folder")
            return false
        }
        let bytes = really_read_array(s, count: 1)
        if let result = bytes {
            if result[0] != 0 {
                Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: server reports failure")
                return false
            }
            return true
        }
        Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: failed to read symlink result")
        return false
    } else {
        return false
    }
}

public func uninstall_symlinks(socket_path: String) -> Bool {
    if let s = handshake(socket_path) {
        if !really_write_array(s, bytes: command_symlink_uninstall) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: failed to write uninstall symlinks command")
            return false
        }
        let bytes = really_read_array(s, count: 1)
        if let result = bytes {
            if result[0] != 0 {
                Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: server reports failure")
                return false
            }
            return true
        }
        Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: failed to read symlink result")
        return false
    } else {
        return false
    }
}

public func uninstall_sockets(socket_path: String) -> Bool {
    if let s = handshake(socket_path) {
        if !really_write_array(s, bytes: command_uninstall_sockets) {
            Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: failed to write uninstall sockets command")
            return false
        }
        let bytes = really_read_array(s, count: 1)
        if let result = bytes {
            if result[0] != 0 {
                Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: server reports failure")
                return false
            }
            return true
        }
        Logger.log(level: Logger.Level.Error, content: "com.docker.vmnetd symlink: failed to read symlink result")
        return false
    } else {
        return false
    }
}

