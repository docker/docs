//
//  DiagnoseSheetViewController.swift
//  Docker
//
//  Created by Gaetan de Villele on 8/15/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Cocoa

class DiagnoseSheetViewController: NSViewController {
    
    var diagnoseAndUploadCallback: (() -> Void)?
    var diagnoseOnlyCallback: (() -> Void)?
    
    // subviews
    @IBOutlet weak var titleTextView: TextViewInView?
    @IBOutlet weak var leftButton: NSButton?
    @IBOutlet weak var rightButton: NSButton?
    @IBOutlet weak var leftTextView: TextViewInView?
    @IBOutlet weak var rightTextView: TextViewInView?
    
    override func viewDidLoad() {
        self.titleTextView?.text = "DiagnoseAndFeedbackSheetTitle".localize()
        self.leftButton?.title = "DiagnoseAndFeedbackSheetLeftButton".localize()
        self.rightButton?.title = "DiagnoseAndFeedbackSheetRightButton".localize()
        self.leftTextView?.text = "DiagnoseAndFeedbackSheetLeftDescription".localize()
        self.rightTextView?.text = "DiagnoseAndFeedbackSheetRightDescription".localize()
    }
        
    @IBAction func diagnoseAndUploadButtonActionReceived(sender: NSButton) {
        self.dismissController(sender)
        if let callback = diagnoseAndUploadCallback {
            callback()
        }
    }
    
    @IBAction func diagnoseOnlyButtonActionReceived(sender: NSButton) {
        self.dismissController(sender)
        if let callback = diagnoseOnlyCallback {
            callback()
        }
    }
}
