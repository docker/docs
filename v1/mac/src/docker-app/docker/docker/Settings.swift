//
//  Settings.swift
//  docker
//
//  Created by Gaetan de Villele on 12/14/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Cocoa


protocol SettingsPanelProtocol {
    /// If something has to be applied, the panel view controller should
    /// return true.
    func shouldApply() -> Bool
    /// Asks panel view controller to apply changes
    func apply()
}

class SettingsPanelRequest {
    var value: PanelRequest = .General
    var animated: Bool = false
    
    init(_ value: PanelRequest, animated: Bool = false) {
        self.value = value
        self.animated = animated
    }
    
    enum PanelRequest {
        case General
        case Advanced
        case FileSharing
        case Cleanup
        case Storage
        case Privacy
    }
}

// Settings class represents the Settings window (UI).
//
// NSTableViewDelegate and NSTableViewDataSource protocols are implemented for
// the multiple NSTableView objects used in the settings panel.
class Settings: NSWindowController, NSWindowDelegate {
    
    
    // VIEWS / PANELS
    
    @IBOutlet weak var generalViewController: NSViewController?
    @IBOutlet weak var advancedViewController: NSViewController?
    @IBOutlet weak var filesharingViewController: NSViewController?
    @IBOutlet weak var cleanupViewController: NSViewController?
    @IBOutlet weak var storageViewController: NSViewController?
    @IBOutlet weak var privacyViewController: NSViewController?
    
    let dockerStateViewController = DockerStateViewController(nibName: "DockerStateView", bundle: nil)
    
    // TOOLBAR
    
    @IBOutlet weak var toolbar: NSToolbar?
    
    @IBOutlet weak var generalButton: NSToolbarItem?
    @IBOutlet weak var advancedButton: NSToolbarItem?
    @IBOutlet weak var filesharingButton: NSToolbarItem?
    @IBOutlet weak var cleanupButton: NSToolbarItem?
    @IBOutlet weak var storageButton: NSToolbarItem?
    @IBOutlet weak var privacyButton: NSToolbarItem?
    
    // animateNextPanel is always true after calling toolbarButtonAction
    // but it can be set to false right before to avoid animation
    private var animateNextPanel: Bool = true
    
    private var currentSettingsPanelProtocol: SettingsPanelProtocol?
    private var currentToolBarItem: NSToolbarItem?
    
    private func restoreToolBarItem() {
        if let currentToolBarItem = currentToolBarItem {
            self.toolbar?.selectedItemIdentifier = currentToolBarItem.itemIdentifier
        }
    }
    
    @IBAction func toolbarButtonAction(toolbarItem: NSToolbarItem) {
        
        // makes sure all changes are applied before opening new panel
        if !applyAndContinue() {
            restoreToolBarItem()
            return
        }
        
        guard let stateView = dockerStateViewController?.view else {
            restoreToolBarItem()
            return
        }
        
        var panelViewController: NSViewController?
        
        if toolbarItem === generalButton {
            panelViewController = generalViewController
        } else if toolbarItem === advancedButton {
            panelViewController = advancedViewController
        } else if toolbarItem === filesharingButton {
            panelViewController = filesharingViewController
        } else if toolbarItem === privacyButton {
            panelViewController = privacyViewController
        } else if toolbarItem === cleanupButton {
            if !requireAdminRights() {
                restoreToolBarItem()
                return
            }
            panelViewController = cleanupViewController
        } else if toolbarItem === storageButton {
            panelViewController = storageViewController
        }
        
        guard let spp = panelViewController as? SettingsPanelProtocol else {
            restoreToolBarItem()
            return
        }
        
        currentSettingsPanelProtocol = spp
        
        guard let panelView = panelViewController?.view else {
            restoreToolBarItem()
            return
        }
        
        window?.title = toolbarItem.label
        currentToolBarItem = toolbarItem
        
        let fittingSize = panelView.fittingSize
        panelView.setFrameSize(fittingSize)
        
        stateView.setFrameSize(NSSize(width: panelView.frame.size.width, height: stateView.frame.height))
        self.toolbar?.selectedItemIdentifier = toolbarItem.itemIdentifier
        
        let containerView = NSView(frame: NSRect(origin: CGPointZero, size: CGSize(width: panelView.frame.size.width, height: panelView.frame.size.height + stateView.frame.size.height)))
        containerView.addSubview(stateView)
        
        panelView.setFrameOrigin(NSPoint(x: 0, y: stateView.frame.size.height))
        
        containerView.addSubview(panelView)
        self.window?.changeView(containerView, animate: animateNextPanel)
        animateNextPanel = true
    }
    
    private func requireAdminRights() -> Bool {
        if Utils.userIsAdminOrRoot {
            return true
        }
        let script = "do shell script \"\" with administrator privileges"
        guard let appleScript = NSAppleScript(source: script) else {
            return false
        }
        var dict: NSDictionary? = nil
        appleScript.executeAndReturnError(&dict)
        if (dict?["NSAppleScriptErrorNumber"]) != nil {
            return false
        }
        return true
    }
    
    // MARK: NSWindowController implementation
    
    // when the settings window gets displayed, populate UI components
    // with actual settings provided by the backend.
    override func showWindow(sender: AnyObject?) {
        super.showWindow(sender)
        // set window's position
        window?.center()
        
        dockerStateViewController?.addGreyBackground()
        
        if let panelRequest = sender as? SettingsPanelRequest {
            
            var toolBarItem: NSToolbarItem?
            
            switch panelRequest.value {
            case .Advanced:
                toolBarItem = advancedButton
            case .Cleanup:
                toolBarItem = cleanupButton
            case .FileSharing:
                toolBarItem = filesharingButton
            case .Privacy:
                toolBarItem = privacyButton
            case .Storage:
                toolBarItem = storageButton
            default:
                toolBarItem = generalButton
            }
            
            if let tbi = toolBarItem {
                animateNextPanel = panelRequest.animated
                toolbarButtonAction(tbi)
            }
        }
        
        // if the toolbar has no selected item, we select the default item ("general" for now)
        // This occurs the first time the user opens the settings window after launching Docker for Mac
        else if toolbar?.selectedItemIdentifier == nil {
            if let generalButton = generalButton {
                animateNextPanel = false
                toolbarButtonAction(generalButton)
            }
        }
    }
    
    
    // MARK: NSWindowDelegate implementation
    
    // This is part of the NSWindowDelegate protocol.
    // when closing the settings window, we check if some advanced settings
    // have not been applied, and prompt the user with a warning, offering
    // different options: don't apply, cancel, or apply.
    func windowShouldClose(sender: AnyObject) -> Bool {
        return applyAndContinue()
    }
    
    private func applyAndContinue() -> Bool {
        if let currentSettingsPanelProtocol = self.currentSettingsPanelProtocol {
            if currentSettingsPanelProtocol.shouldApply() {
                var b = false
                if let message = WizardMessageGeneric(message: "settingsApplySettingsTitle".localize(), details: "settingsApplySettingsMessage".localize(), icon: "MsgWarningIcon") {
                    
                    message.addButton("settingsApplySettingsBtnDont".localize()).blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                        b = true
                        message?.close()
                    }
                    
                    message.addButton("settingsApplySettingsBtnCancel".localize()).blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                        message?.close()
                    }
                    
                    message.addButton("settingsApplySettingsBtnApply".localize()).makeDefault().blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                        b = true
                        currentSettingsPanelProtocol.apply()
                        message?.close()
                    }
                    Wizard.show(message)
                }
                return b
            }
        }
        return true
    }
    
    // This is part of the NSWindowDelegate protocol.
    // This is called before the window is closed.
    func windowWillClose(notification: NSNotification) {
        // Unregister <self> from the NotificationCenter.
        NSNotificationCenter.defaultCenter().removeObserver(self, name: Backend.dockerStateNotificationName, object: nil)
    }
    
    // MARK: utility functions
    
    func showErrorPopupWithTitle(title: String, andMessage message: String) {
        guard let msg = WizardMessageGeneric(message: title, details: message, icon: "MsgWarningIcon") else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
            return
        }
        let closeBtn = msg.addCloseButton("Ok".localize())
        closeBtn.makeDefault()
        Wizard.show(msg)
    }
    
    func showErrorPopupAndCloseWithTitle(title: String, andMessage message: String) {
        showErrorPopupWithTitle(title, andMessage: message)
        self.close()
    }
}