//
//  SettingsFilesharing.swift
//  Docker
//
//  Created by Gaetan de Villele on 6/29/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class SettingsFilesharing: NSViewController, SettingsPanelProtocol, TextTableViewEventListener {
    
    @IBOutlet weak var settingsWindow: Settings?
    @IBOutlet weak var tableView: TextTableView?
    @IBOutlet weak var linkToDocView: TextViewInView?
    
    private var initialDirectories = [String]()
    
    override func viewDidDisappear() {
        super.viewDidDisappear()
        tableView?.deselectAll(self)
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
        
        tableView?.eventListener = self
        tableView?.editable = false
        
        lockUIComponents()
        Backend.getSharedDirectories { (error, directories) in
            if let errorValue: Backend.Error = error {
                // Failed to get the advanced settings value from the backend.
                // We display an error popup and close the settings window.
                let errorMessage = Backend.errorMessage(errorValue)
                self.settingsWindow?.showErrorPopupAndCloseWithTitle("Error".localize(), andMessage: "settingsApiGeneralRetrieveError".localize(errorMessage))
                return
            }
            // no error
            self.initialDirectories = directories
            self.tableView?.strings = directories
            
            self.unlockUIComponents()
        }
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        linkToDocView?.text = "settingsFileSharingDocLink".localize()
    }
    
    // MARK: UI components
    
    @IBOutlet weak var addButton: NSButton?
    @IBAction func addButtonActionReceived(_: NSButton) {
        // use a modal OpenPanel
        let openPanel: NSOpenPanel = NSOpenPanel()
        openPanel.canChooseFiles = false
        openPanel.canChooseDirectories = true
        openPanel.showsHiddenFiles = true
        openPanel.resolvesAliases = false
        openPanel.allowsMultipleSelection = false
        if openPanel.runModal() == NSModalResponseOK {
            if let url: NSURL = openPanel.URLs.first {
                if let path: String = url.path {
                    self.tableView?.appendString(path)
                    return
                }
            }
            self.settingsWindow?.showErrorPopupWithTitle("Error".localize(), andMessage: "settingsApiDirectorySelectionError".localize())
        }
    }
    @IBOutlet weak var removeButton: NSButton?
    @IBAction func removeButtonActionReceived(_: NSButton) {
        self.tableView?.removeSelectedString()
    }
    
    @IBOutlet weak var applyButton: NSButton?
    @IBAction func applyButtonActionReceived(button: NSButton) {
        guard let tableView = self.tableView else { return }
        lockUIComponents()
        
        Backend.setSharedDirectories(tableView.strings) { (error, directories) in
            if let errorValue = error {
                self.settingsWindow?.showErrorPopupWithTitle("Error".localize(), andMessage: Backend.errorMessage(errorValue))
                tableView.strings = self.initialDirectories
            } else {
                tableView.strings = directories
                self.initialDirectories = directories
            }
            self.checkIfFilesharingSettingsChanged()
            self.unlockUIComponents()
        }
        Backend.enforceStarting()
    }
    
    func lockUIComponents() {
        tableView?.enabled = false
        addButton?.enabled = false
        removeButton?.enabled = false
    }
    
    func unlockUIComponents() {
        tableView?.enabled = true
        addButton?.enabled = true
    }
    
    func checkIfFilesharingSettingsChanged() -> Bool {
        guard let tableView = self.tableView else { return false }
        let changed: Bool = initialDirectories != tableView.strings
        applyButton?.enabled = changed
        return changed
    }
    
    // MARK: TextTableViewEventListener implementation
    
    func textTableViewGainedFocus(textTableView: TextTableView) {}
    
    func textTableViewLostFocus(textTableView: TextTableView) {}
    
    func textTableViewContentDidChange(textTableView: TextTableView) {
        checkIfFilesharingSettingsChanged()
    }
    
    func textTableViewCanRemoveRow(textTableView: TextTableView, canRemoveRow: Bool) {
        removeButton?.enabled = canRemoveRow
    }
    
    func textTableViewCanAddRow(textTableView: TextTableView, canAddRow: Bool) {
        addButton?.enabled = canAddRow
    }
    
    // MARK: SettingsPanelProtocol implementation
    
    func shouldApply() -> Bool {
        return checkIfFilesharingSettingsChanged()
    }
    
    func apply() {
        if let applyButton = applyButton {
            applyButtonActionReceived(applyButton)
        }
    }
}
