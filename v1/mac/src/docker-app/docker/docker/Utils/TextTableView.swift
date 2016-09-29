//
//  TextTableView.swift
//  docker
//
//  Created by Doby Mock on 7/1/16.
//  Copyright © 2016 docker. All rights reserved.
//

// TextTableView defines a specific usage of NSTableView, where we only
// need to display one column of editable string entries.

import Foundation
import Cocoa

public protocol TextTableViewEventListener {
    func textTableViewGainedFocus(textTableView: TextTableView)
    func textTableViewLostFocus(textTableView: TextTableView)
    func textTableViewContentDidChange(textTableView: TextTableView)
    // can be used to enable/disable '+' & '-' buttons
    func textTableViewCanRemoveRow(textTableView: TextTableView, canRemoveRow: Bool)
    func textTableViewCanAddRow(textTableView: TextTableView, canAddRow: Bool)
}

public class TextTableView: NSTableView, NSTableViewDelegate, NSTableViewDataSource, NSTextFieldDelegate {
 
    var eventListener: TextTableViewEventListener?
    var placeholderString: String = ""
    var editable: Bool = true
    
    public required init?(coder: NSCoder) {
        super.init(coder: coder)
        self.delegate = self
        self.dataSource = self
    }
    
    override init(frame frameRect: NSRect) {
        super.init(frame: frameRect)
        self.delegate = self
        self.dataSource = self
    }
    
    private var _strings = [String]()
    var strings: [String] {
        get {return self._strings}
        set {
            self._strings = newValue
            self.reloadData()
            self.eventListener?.textTableViewContentDidChange(self)
        }
    }
    
    func appendString(s: String, edit: Bool = false) {
        _strings.append(s)
        self.reloadData()
        eventListener?.textTableViewContentDidChange(self)
        if edit {
            editStringAtIndex(_strings.count - 1)
            if s == "" {
                eventListener?.textTableViewCanAddRow(self, canAddRow: false)
            }
        }
    }
    
    func removeStringAtIndex(index: Int) {
        _strings.removeAtIndex(index)
        self.reloadData()
        eventListener?.textTableViewContentDidChange(self)
    }
    
    /// removes selected string if any is selected
    func removeSelectedString() {
        let selectedRow = self.selectedRow
        // test if a row is actually selected
        if selectedRow > -1 {
            // get the text field at the selected row
            if let textfield = self.viewAtColumn(0, row: selectedRow, makeIfNecessary: false) as? NSTextField {
                // test whether the selected text field is being edited
                if let editor: NSText = textfield.currentEditor() {
                    // being edited
                    editor.string = ""
                    // force text field to end editing
                    textfield.window?.makeFirstResponder(nil)
                    // Ending editing on a text field that has an empty editor,
                    // will automatically remove the row from the table view.
                    // See TextTableView.controlTextDidEndEditing(:)
                } else {
                    // text field is not being edited
                    self.removeStringAtIndex(selectedRow)
                }
            }
        }
    }
    
    func editStringAtIndex(index: Int) {
        self.editColumn(0, row: index, withEvent: nil, select: true)
        self.selectRowIndexes(NSIndexSet.init(index: index), byExtendingSelection: false)
    }
    
    // MARK: NSTableViewDelegate implementation
    
    public func tableViewSelectionDidChange(notification: NSNotification) {
        eventListener?.textTableViewCanRemoveRow(self, canRemoveRow: self.selectedRow != -1)
    }
    
    public func tableView(tableView: NSTableView, didRemoveRowView rowView: NSTableRowView, forRow row: Int) {
        eventListener?.textTableViewCanRemoveRow(self, canRemoveRow: false)
        eventListener?.textTableViewCanAddRow(self, canAddRow: true)
    }
    
    public func tableView(tableView: NSTableView, viewForTableColumn tableColumn: NSTableColumn?, row: Int) -> NSView? {
        // check if there is a column
        guard let column = tableColumn else {
            return nil
        }
        let text: NSTextField = NSTextField(frame: NSRect(x: 0, y: 0, width: column.width, height: tableView.rowHeight))
        text.bordered = false
        text.drawsBackground = false
        text.editable = editable
        text.placeholderString = self.placeholderString
        text.delegate = self
        if _strings.count > row {
            text.stringValue = _strings[row]
        } else {
            text.stringValue = "error"
        }
        return text
    }
    
    // MARK: NSTextFieldDelegate Protocol implementation
    
    public override func controlTextDidChange(notification: NSNotification) {
        if let object: AnyObject = notification.object {
            if let textfield: NSTextField = object as? NSTextField {
                for i in 0..<self.numberOfRows {
                    if textfield === self.viewAtColumn(0, row: i, makeIfNecessary: false) { // ===
                        eventListener?.textTableViewCanAddRow(self, canAddRow: textfield.stringValue != "")
                        _strings[i] = textfield.stringValue
                        eventListener?.textTableViewContentDidChange(self)
                        return
                    }
                }
            }
        }
    }
    
    public override func controlTextDidEndEditing(notification: NSNotification) {
        if let object: AnyObject = notification.object {
            if let textfield: NSTextField = object as? NSTextField {
                for i in 0..<self.numberOfRows {
                    if textfield === self.viewAtColumn(0, row: i, makeIfNecessary: false) {
                        if textfield.stringValue == "" {
                            _strings.removeAtIndex(i)
                            self.removeRowsAtIndexes(NSIndexSet(index: i), withAnimation: .EffectNone)
                        }
                        eventListener?.textTableViewCanAddRow(self, canAddRow: true)
                        return
                    }
                }
            }
        }
    }
    
    // MARK: NSTableViewDataSource implementation
    
    public func numberOfRowsInTableView(tableView: NSTableView) -> Int {
        return _strings.count
    }

    // MARK: TextTableViewEventListener events
    
    override public func becomeFirstResponder() -> Bool {
        eventListener?.textTableViewGainedFocus(self)
        return super.becomeFirstResponder()
    }
    
    override public func resignFirstResponder() -> Bool {
        eventListener?.textTableViewLostFocus(self)
        return super.resignFirstResponder()
    }
}


