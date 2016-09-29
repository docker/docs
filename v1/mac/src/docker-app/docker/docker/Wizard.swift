//
//  Wizard.swift
//  docker
//
//  Created by Adrien Duermael on 1/19/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

// The Wizard is a class exposing static functions that can be used to 
// to display modal windows containing a single WizardMessage or 
// queued WizardMessages.
public class Wizard {
    static func show(messages: WizardMessage...){
        let wwc = WizardWindowController(windowNibName: "WizardWindow")
        
        for message in messages {
            wwc.pushBack(message)
        }
        
        // set the position of the window to
        // "center horizontally and somewhat above center vertically"
        if let window = wwc.window {
            window.delegate = wwc
            window.center()
            NSApp.runModalForWindow(window)
        } else {
            let content = "Wizard.show could not get window from WizardWindowController"
            Logger.log(level: Logger.Level.Error, content: content)
        }
    }
    
    static private var nonModalDisplayed: [String : WizardWindowController] = [String : WizardWindowController]()
    
    /// tries to select a displayed non-modal WizardMessage and brings
    /// it back to front. Returns false if the message can't be found.
    static func selectNonModal(key: String) -> Bool {
        
        if let wwc: WizardWindowController = nonModalDisplayed[key] {
            wwc.showWindow(wwc.window)
            wwc.window?.center()
            return true
        }
        
        return false
    }
    
    static func showNonModal(message: WizardMessage, key: String) {
        
        if let _: WizardWindowController = nonModalDisplayed[key] {
            Logger.log(level: ASLLogger.Level.Error, content: "Wizard.showNonModal: key already used")
            return
        }
        
        let wwc = WizardWindowController(windowNibName: "WizardWindow")
        if let window = wwc.window {
            
            wwc.nonModalKey = key
            nonModalDisplayed[key] = wwc
            
            window.delegate = wwc
            wwc.pushBack(message)
            wwc.showWindow(window)
            window.center()
        } else {
            let content = "Wizard.showNonModal could not get window from WizardWindowController"
            Logger.log(level: Logger.Level.Error, content: content)
        }
    }
}


// WizardWindowController for WizardWindow
public class WizardWindowController: NSWindowController, NSWindowDelegate {
    
    // The WizardMessage being displayed (front WizardMessage in queue)
    private var displayed: WizardMessage?
    private var nonModalKey: String?
    
    private func setDisplayed(message: WizardMessage) {
        if displayed == nil {
            displayed = message
        } else {
            displayed?.view.removeFromSuperview()
            message.next = displayed
            displayed?.previous = message
            displayed = message
        }
        if let displayed = displayed {
            window?.contentView = displayed.view
            window?.styleMask = displayed.styleMask
            if let title = displayed.title {
                window?.title = title
            } else {
                window?.title = ""
            }
        }
    }

    
    // pushBack adds a message at the end of the queue
    // it will be displayed as soon as all messages
    // in front are gone.
    private func pushBack(message: WizardMessage) {
        message.wizardWindowController = self
        if displayed == nil {
            setDisplayed(message)
        } else {
            var msg = displayed
            while msg?.next != nil {
                msg = msg?.next
            }
            msg?.next = message
            message.previous = msg
        }
    }
    
    // pushFront adds a message at the front of the queue
    // it will be displayed right away.
    private func pushFront(message: WizardMessage) {
        message.wizardWindowController = self
        setDisplayed(message)
    }
    
    // insertAfter inserts a message after given message
    // it will be displayed after target is gone.
    private func insertAfter(message: WizardMessage, target: WizardMessage) {
        message.wizardWindowController = self
        var msg = displayed
        while msg != nil && msg != target {
            msg = msg?.next
        }
        if msg == target {
            // TODO: INSERT
        }
    }
    
    // insertBefore inserts a message before given message
    // it will be displayed before the target.
    private func insertBefore(message: WizardMessage, target: WizardMessage) {
        message.wizardWindowController = self
        var msg = displayed
        while msg != nil && msg != target {
            msg = msg?.next
        }
        if msg == target {
            // TODO: INSERT
        }
    }
    
    public func windowShouldClose(sender: AnyObject) -> Bool {
        if let displayed = self.displayed {
            close(displayed)
        }
        return false
    }
    
    // close removes a message from the queue. It makes it
    // disappear if currently displayed.
    private func close(message: WizardMessage) {
        if message == displayed {
            displayed?.view.removeFromSuperview()
            displayed = message.next
        }
        
        if let displayed = displayed {
            // display message
            window?.contentView = displayed.view
            window?.styleMask = displayed.styleMask
            if let title = displayed.title {
                window?.title = title
            } else {
                window?.title = ""
            }
        } else {
            // Wizard windows always are modals
            if let key = nonModalKey {
                Wizard.nonModalDisplayed.removeValueForKey(key)
            }
            
            window?.close()
            NSApp.stopModal()
        }
    }
}

// WizardWindow displays single or queued WizardMessages
class WizardWindow: NSWindow {
}

// WizardMessage defines a message that can be displayed in the WizardWindow
class WizardMessage: NSViewController {
    
    // WizardWindow displaying that message
    var wizardWindowController: WizardWindowController?
    // the previous message, the one displayed before.
    var previous: WizardMessage?
    // the next message, will be displayed as soon as self closes.
    var next: WizardMessage?
    
    var viewFrameCache: NSRect
    
    // used by wizard message to enforce window's style
    // (buttons in the window's top bar)
    private var styleMask: NSWindowStyleMask = NSTitledWindowMask
    private var _closeButton = false
    private var _minimizeButton = false
    
    // decide wether or not a close button should be displayed in the
    // window's top bar. This has to be set before the message gets displayed
    var closeButton: Bool {
        get {
            return _closeButton
        }
        set {
            _closeButton = newValue
            updateStyleMask()
        }
    }
    
    var minimizeButton: Bool {
        get {
            return _minimizeButton
        }
        set {
            _minimizeButton = newValue
            updateStyleMask()
        }
    }
    
    private func updateStyleMask() {
        styleMask = NSTitledWindowMask
        if _closeButton {
            styleMask = [styleMask, NSClosableWindowMask]
        }
        if _minimizeButton {
            styleMask = [styleMask, NSMiniaturizableWindowMask]
        }
    }
    
    @IBOutlet weak var buttonViewWidthConstraint: NSLayoutConstraint?
    @IBOutlet weak var buttonViewHeightConstraint: NSLayoutConstraint?
    
    // the view that contains buttons
    @IBOutlet weak var buttonView: NSView?
    // Message's buttons are stored here, buttonView
    // may not be ready when adding them (only available
    // after view did actually load)
    var buttons = [WizardButton]()
    
    // var buttonViewWidthConstraint: NSLayoutConstraint?
    // var buttonViewHeightConstraint: NSLayoutConstraint?
    
    // margin between buttons
    let marginBetweenButtons: CGFloat = 0.0
    
    override init?(nibName: String?, bundle: NSBundle?) {
        viewFrameCache = NSRect()
        super.init(nibName: nibName, bundle: bundle)
    }
    
    required init?(coder: NSCoder) {
        viewFrameCache = NSRect()
        super.init(coder: coder)
    }
    
    override func viewDidLoad() {
        refreshButtons()
    }
    
    // updateViewConstraints is triggered when layout
    // constraints are being updated. The current frame
    // position and size can be saved here to compare
    // with new values in viewWillLayout() and center window
    // based on previous position.
    // Note(aduermae): updateViewConstraints seems to be called
    // only when view gets visible, not for further layout updates
    override func updateViewConstraints() {
        // Logger.log(level: Logger.Level.Notice, content: "updateViewConstraints \(view.frame)")
        viewFrameCache = view.frame
        super.updateViewConstraints()
    }
    
    // viewWillLayout is triggered just before view layout
    // is updated, right after updateViewConstraints(). The
    // view already has its new position and size. It should
    // be updated to keep it centered giving previous situation.
    override func viewWillLayout() {
//         Logger.log(level: Logger.Level.Notice, content: "viewWillLayout \(view.frame)")
//         Logger.log(level: Logger.Level.Notice, content: "viewFrameCache \(viewFrameCache)")
//         Logger.log(level: Logger.Level.Notice, content: "---")
        
        if let window = wizardWindowController?.window {
            let x = window.frame.origin.x + (viewFrameCache.size.width - view.frame.size.width) * 0.5
            let y = window.frame.origin.y - (viewFrameCache.size.height - view.frame.size.height) * 0.5
            let origin = CGPoint(x: x, y: y)
            let rect = NSRect(origin: origin, size: window.frame.size)
            window.setFrame(rect, display: true)
        } else {
            let content = "WizardMessage could not get window from WizardWindowController"
            Logger.log(level: Logger.Level.Error, content: content)
        }
        
        viewFrameCache = view.frame
        
        super.viewWillLayout()
    }
    
    override func viewWillAppear() {
//        Logger.log(level: Logger.Level.Notice, content: "viewWillAppear \(view.frame)")
        refreshButtons()
        super.viewWillAppear()
    }
    
    // close removes the message from the message queue of its WizardWindow.
    // It will disappear if being displayed.
    func close() {
        // connect next & previous
        previous?.next = next
        next?.previous = previous
        wizardWindowController?.close(self)
    }
    
    // addButton adds a button for the Message placed on the left
    // of last button, all buttons are anchored in bottom right corner
    // Window auto size should make sure all buttons fit, including margins.
    func addButton(text: String) -> WizardButton {
        
        var posX: CGFloat = 0.0
        
        // update posX if there are buttons in buttonView already
        if let previousButton = buttons.last {
            posX = previousButton.frame.origin.x + previousButton.frame.size.width + marginBetweenButtons
        }
        
        let button = WizardButton(frame: NSRect(x:posX, y:0, width: 0, height:32.0))
        button.message = self
        button.bezelStyle = NSBezelStyle.Rounded
        button.setButtonType(NSButtonType.MomentaryPushIn)
        button.title = text
        button.sizeToFit()
        
        let hPadding: CGFloat = 20.0
        
        button.setFrameSize(NSSize(width: button.frame.size.width + hPadding, height: button.frame.size.height))
        
        buttons.append(button)
        
        // display buttons that are not displayed, if view is loaded
        refreshButtons()
        
        return button
    }
    
    // addCloseButton adds a button that simply closes the Message (self)
    func addCloseButton(text: String) -> WizardButton {
        let btn = addButton(text)
        btn.target = self
        btn.action = #selector(WizardMessage.close)
        return btn
    }
    
    // pushAfter inserts a new message in the queue, right after self
    func pushAfter(message: WizardMessage) {
        // TODO: fix implementation currently it pushed the message at
        // the end of the queue, not right after self message
        wizardWindowController?.pushBack(message)
    }
    
    // pushAfter inserts a new message in the queue, right in front of self
    func pushBefore(message: WizardMessage) {
        // TODO: fix implementation currently it pushed the message at
        // the very front of the queue, not right in front of self message
        wizardWindowController?.pushFront(message)
    }
    
    // refreshButtons display buttons that are not displayed
    func refreshButtons() {
        if let bView = buttonView {
            
            var totalWidth: CGFloat = 0.0
            
            for btn in buttons {
                if btn.superview == nil {
                    bView.addSubview(btn)
                }
                
                totalWidth += btn.frame.size.width + marginBetweenButtons
            }
            if totalWidth > 0 {
                totalWidth -= marginBetweenButtons
            }
            
            var buttonViewHeight: CGFloat = 32.0
            if buttons.count == 0 {
                buttonViewHeight = 0.0
            }
            
            if let heightConstraint = buttonViewHeightConstraint {
                heightConstraint.constant = buttonViewHeight
            }
            
            if let widthConstraint = buttonViewWidthConstraint {
                widthConstraint.constant = totalWidth
            }
        }
    }
    
}

// WizardButton extends NSButton to be able to execute block actions
class WizardButton: NSButton {
    
    private var message: WizardMessage? = nil
    
    /// Do not use this directly. Use blockAction.
    private var _blockAction: ((button: WizardButton?, message: WizardMessage?) -> ())? = nil
    
    var blockAction: ((button: WizardButton?, message: WizardMessage?) -> ())? {
        get {
            return self._blockAction
        }
        set {
            if self._blockAction != nil || self.action != nil || self.target != nil {
                Logger.log(level: ASLLogger.Level.Fatal, content: "💥 bad use of WizardButton")
                return
            }
            self.target = self
            self.action = #selector(WizardButton.actionReceived)
            _blockAction = newValue
        }
    }
    
    func makeDefault() -> WizardButton {
        self.highlight(true)
        self.keyEquivalent = "\r"
        return self
    }
    
    func actionReceived(sender: WizardButton) {
        if let action = blockAction {
            action(button:self, message: self.message)
        }
    }
}
