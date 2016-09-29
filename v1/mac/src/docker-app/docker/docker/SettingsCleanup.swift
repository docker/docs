//
//  SettingsCleanup.swift
//  docker
//
//  Created by Doby Mock on 6/27/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class SettingsCleanup: NSViewController, SettingsPanelProtocol {

    @IBOutlet weak var settingsWindow: Settings?
    
    // resetToFactoryDefaultsPressed erases data (containers, images...etc)
    // but most settings are kept (excepting localMachineMigration)
    @IBAction func resetToFactoryDefaultsPressed(sender: NSButton) {
        
        guard let resetMessage = WizardMessageGeneric(message: "factoryResetWarningTitle".localize(), details: "factoryResetWarningMessage".localize(), icon: "MsgWarningIcon") else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
            return
        }
        
        let okBtn = resetMessage.addButton("factoryResetWarningOkButton".localize())
        resetMessage.addCloseButton("factoryResetWarningCancelButton".localize()).makeDefault()
        
        okBtn.blockAction = {(button: WizardButton?, resetWarningMessage: WizardMessage?) -> () in
            
            // close the settings window during reset
            self.settingsWindow?.window?.close()
            
            // disable status item during reset
            if let delegate: AppDelegate = NSApp.delegate as? AppDelegate {
                delegate.statusItem?.Disable()
            }
            
            // creating reset progress bar popup
            guard let progressMessage = WizardMessageGeneric(message: "factoryResetWarningTitle".localize(), details: "factoryResetProgressMessage".localize(), icon: "MsgCleaningIcon") else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                return
            }
            let progressBar = progressMessage.addProgressBar()
            progressMessage.addCustomComponent(progressBar)
            resetWarningMessage?.pushAfter(progressMessage)
            resetWarningMessage?.close()
            
            // dispatch application reset and restart
            dispatch_async(dispatch_get_global_queue(QOS_CLASS_DEFAULT, 0), { () -> Void in
                
                let error: String? = Uninstall.uninstallWhileKeepingPreferenceKeys([Preferences.Key.analyticsUserID])
                
                dispatch_async(dispatch_get_main_queue(), { () -> Void in
                    
                    if let error: String = error {
                        progressMessage.close()
                        Logger.log(level: ASLLogger.Level.Fatal, content: "factoryResetFailed".localize(error))
                    } else {
                        guard let resetSuccessMessage = WizardMessageGeneric(message: "factoryResetWarningTitle".localize(), details: "factoryResetSuccessMessage".localize(), icon: "MsgGenericIcon") else {
                            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                            return
                        }
                        let restartButton = resetSuccessMessage.addButton("factoryResetRestartButton".localize())
                        restartButton.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                            // close message
                            message?.close()
                            // restart application
                            let task = NSTask()
                            task.launchPath = "/usr/bin/open"
                            task.arguments = [NSBundle.mainBundle().bundlePath]
                            task.launch()
                            exit(EXIT_SUCCESS)
                        }
                        progressMessage.pushAfter(resetSuccessMessage)
                    }
                    progressMessage.close()
                })
            })
            
        }
        Wizard.show(resetMessage)
    }
    
    // uninstall gets rid of everything that's been installed by Docker for Mac
    // Then it displays a popup inviting user to put Docker.app in the trash
    @IBAction func uninstallDocker(_: NSButton) {
        
        guard let uninstallMessage = WizardMessageGeneric(message: "uninstallWarningTitle".localize(), details: "uninstallWarningMessage".localize(), icon: "MsgWarningIcon") else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
            return
        }
        
        let okBtn = uninstallMessage.addButton("uninstallWarningOkButton".localize())
        uninstallMessage.addCloseButton("uninstallWarningCancelButton".localize()).makeDefault()
        
        okBtn.blockAction = {(button: WizardButton?, resetWarningMessage: WizardMessage?) -> () in
            
            // close the settings window during reset
            self.settingsWindow?.window?.close()
            
            // disable status item during reset
            if let delegate: AppDelegate = NSApp.delegate as? AppDelegate {
                delegate.statusItem?.Disable()
            }
            
            // creating reset progress bar popup
            guard let progressMessage = WizardMessageGeneric(message: "uninstallWarningTitle".localize(), details: "uninstallProgressMessage".localize(), icon: "MsgCleaningIcon") else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                return
            }
            let progressBar = progressMessage.addProgressBar()
            progressMessage.addCustomComponent(progressBar)
            resetWarningMessage?.pushAfter(progressMessage)
            resetWarningMessage?.close()
            
            // dispatch application reset and restart
            dispatch_async(dispatch_get_global_queue(QOS_CLASS_DEFAULT, 0), { () -> Void in
                
                let error: String? = Uninstall.uninstallWhileKeepingPreferenceKeys([Preferences.Key.analyticsUserID])
                
                dispatch_async(dispatch_get_main_queue(), { () -> Void in
                    
                    if let error: String = error {
                        progressMessage.close()
                        Logger.log(level: ASLLogger.Level.Fatal, content: "uninstallFailed".localize(error))
                    } else {
                        guard let uninstallSuccessMessage = WizardMessageGeneric(message: "uninstallWarningTitle".localize(), details: "uninstallSuccessMessage".localize(), icon: "MsgTrashIcon") else {
                            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                            return
                        }
                        
                        uninstallSuccessMessage.addImageWithinMessage(Utils.isStableBuild() ? "MoveToTrash" : "MoveToTrashBeta")
                        
                        let quitButton = uninstallSuccessMessage.addButton("uninstallSuccessQuitButton".localize())
                        quitButton.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                            // close message
                            message?.close()
                            exit(EXIT_SUCCESS)
                        }
                        quitButton.makeDefault()
                        progressMessage.pushAfter(uninstallSuccessMessage)
                    }
                    progressMessage.close()
                })
            })
            
        }
        Wizard.show(uninstallMessage)
    }
    

    func shouldApply() -> Bool {
        return false
    }
    
    func apply() {
        // nothing to apply
    }
}
