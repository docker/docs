//
//  Dashboard.swift
//  docker
//
//  Created by Adrien Duermael on 1/18/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

// The dashboard is the main app graphical interface. It provides a visual 
// representation for CLI commands.
class Dashboard: NSWindowController, NSWindowDelegate {
    
    @IBOutlet weak var tmpView: NSView?
    
    let logViewController = LogViewController()
    
    // showWindow is triggered each time the window gets displayed
    override func showWindow(sender: AnyObject?) {
        super.showWindow(sender)
        // set the position of the window to
        // "center horizontally and somewhat above center vertically"
        self.window?.center()
        
        // make app process foreground

        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.transformProcess(true)
        } else {
            Logger.log(Logger.Level.Fatal, "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
    
    override func windowDidLoad() {
        super.windowDidLoad()

        if let tmpv = tmpView {
            logViewController.view.autoresizingMask = [NSAutoresizingMaskOptions.ViewWidthSizable, NSAutoresizingMaskOptions.ViewHeightSizable]
            logViewController.view.frame = NSRect(origin: CGPointZero, size: tmpv.frame.size)
            tmpv.addSubview(logViewController.view)
        }
    }
}

class DashboardWindow: NSWindow {
    // overrides NSWindow's close to transform app process
    override func close() {
        // process back to agent
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.transformProcess(false)
        } else {
            Logger.log(Logger.Level.Fatal, "Unable to cast NSApp.delegate to AppDelegate")
        }
        super.close()
    }
}