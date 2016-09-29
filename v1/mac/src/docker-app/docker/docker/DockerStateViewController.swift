//
//  DockerStateViewController.swift
//  docker
//
//  Created by Doby Mock on 6/1/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class DockerStateViewController: NSViewController {
    
    // green dot: Docker is running, orange dot: docker is starting
    @IBOutlet weak var dockerStateIcon: NSImageView?
    // textfield: displays docker current status
    @IBOutlet weak var dockerStateLabel: NSTextField?
    
    @IBAction func restartVM(sender: NSButton) {
        Backend.vmRestart { (error: Backend.Error?) in
            if let error = error {
                Logger.log(level: ASLLogger.Level.Error, content: "docker restart error: \(Backend.errorMessage(error))")
            }
        }
        Backend.enforceStarting()
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        let instantNotification = NSNotification(name: Backend.dockerStateNotificationName, object: Backend.getCurrentState().rawValue)
        self.dockerStatusChangedNotificationCallback(instantNotification)
        NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(DockerStateViewController.dockerStatusChangedNotificationCallback(_:)), name: Backend.dockerStateNotificationName, object: nil)
    }

    func addGreyBackground() {
        self.view.wantsLayer = true
        self.view.layer?.backgroundColor = CGColorCreateGenericRGB(0, 0, 0, 0.05)
    }
    
    // this is called by the NotificationCenter
    func dockerStatusChangedNotificationCallback(notification: NSNotification) {
        if let dockerState = notification.object as? String {
             if dockerState == "starting" {
                dockerStateIcon?.image = NSImage(named: NSImageNameStatusPartiallyAvailable)
                dockerStateLabel?.stringValue = "dockerStateStarting".localize()
            } else {
                dockerStateIcon?.image = NSImage(named: NSImageNameStatusAvailable)
                dockerStateLabel?.stringValue = "dockerStateRunning".localize()
            }
        }
    }
}


