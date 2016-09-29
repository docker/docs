//
//  Options.swift
//  docker-installer
//
//  Created by Michel Courtine on 12/5/15.
//  Copyright © 2015 David Scott. All rights reserved.
//

import Foundation

// command-line arguments:
// --no-privileges   do not install vmnetd
// --no-check        do not run environment checks before install

// Options parses exec option flags.
// An Option instance parses flags at initialization.
class Options {
    
    // --unattended argument
    // By default, we consider that a human in using the app.
    // If the app is started during automated tests, the flag "--unattended"
    // must be used.
    static var unattended: Bool {
        get {
            return Process.arguments.contains("--unattended")
        }
    }
    
    // --uninstall argument
    static var uninstall: Bool {
        get {
            return Process.arguments.contains("--uninstall")
        }
    }
    
    // --quit-after-install argument
    static var quitafterinstall: Bool {
        get {
            return Process.arguments.contains("--quit-after-install")
        }
    }
}

