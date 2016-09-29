//
//  WizardMessageGeneric.swift
//  docker
//
//  Created by Adrian Duermael on 02/28/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Cocoa

// WizardMessageGeneric is a WizardMessage implementation.
// It can optionally display an icon on the top left corner
// And a different components from top to bottom on the
// right side (text, loading bar, input label...etc)
class WizardMessageGeneric: WizardMessage {
    
    @IBOutlet weak var componentsView: NSView?
    @IBOutlet weak var componentsHeightConstraint: NSLayoutConstraint?
    @IBOutlet weak var componentsWidthConstraint: NSLayoutConstraint?

    @IBOutlet weak var iconView: NSImageView?
    // The margin between the icon and all components on the right side
    // if there's no icon, margin constant should be equal to zero.
    @IBOutlet weak var iconRightMarginConstraint: NSLayoutConstraint?
    
    // if message & details are set, a text component will be added
    // at the top, with message in bold and details in regular font.
    private var message: String?
    private var details: String?
    // if icon is set, it will be displayed on the top left side of the message
    private var icon: String?
    
    private var components: [NSView] = []
    
    private let fontSize: CGFloat = 13
    private let maxHeight: CGFloat = 10000
    // TODO: vMargin should be equal to the height of a text line
    private let vMargin: CGFloat = 13
    
    init?(message: String? = nil, details: String? = nil, icon: String? = "MsgGenericIcon") {
        self.message = message
        self.details = details
        self.icon = icon
        super.init(nibName: "WizardMessageGeneric", bundle: NSBundle.mainBundle())
        
        if let title = self.message {
           addText("**\(title)**")
        }
        
        if let details = self.details {
            addText(details)
        }
    }
    
    // addText adds a text component for the Wizard message
    func addText(text: String) {
        let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(fontSize)]
        let attributedString = NSMutableAttributedString(string: text, attributes:attrs)
        attributedString.format(fontSize)
        
        let textView = NSTextView()
        textView.selectable = true
        textView.editable = false
        textView.drawsBackground = false
        textView.textStorage?.setAttributedString(attributedString)
        
        // components.append(textField)
        components.append(textView)
        
    }
    
    // addProgressBar adds a progress bar for the Wizard message
    // by detault it shows indeterminated progression.
    // To display progression, use indeterminate = false parameter
    // and set returned progressIndicator minValue, maxValue & doubleValue
    func addProgressBar(indeterminate: Bool = true) -> NSProgressIndicator {
        let progressIndicator = NSProgressIndicator(frame: NSRect(origin: CGPointZero, size: CGSize(width: 0, height: 20)))
        progressIndicator.style = .BarStyle
        progressIndicator.indeterminate = indeterminate
        if indeterminate {
            progressIndicator.startAnimation(nil)
        }
        components.append(progressIndicator)
        return progressIndicator
    }
    
    // addCustomComponent adds a custom component which is an NSView
    func addCustomComponent(component: NSView) {
        components.append(component)
    }
    
    func addButtonWithinMessage(button: NSButton) {
        // TODO: implementation
    }
    
    func addImageWithinMessage(imageName: String) -> NSImageView? {
        let imageView = NSImageView()
        if let image = NSImage(named: imageName) {
            imageView.image = image
            imageView.setFrameSize(image.size)
            components.append(imageView)
            return imageView
        }
        return nil
    }
    
    private func refreshComponents() {
        guard let componentsView = self.componentsView else { return }
        
        // remove all components
        for subview in componentsView.subviews {
            subview.removeFromSuperview()
        }
        
        var posY: CGFloat = 0
        
        for component in components.reverse() {
            
            if let textField = component as? NSTextField {
                // if textField has an attributed string, it will be autosized
                // otherwize only width is adapted (to fit w/ other components)
                if textField.attributedStringValue.string != "" {
                    let boundingRect = textField.attributedStringValue.boundingRectWithSize(NSSize(width: componentsView.frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
                    textField.setFrameSize(NSSize(width: componentsView.frame.size.width, height: boundingRect.size.height))
                } else {
                    textField.setFrameSize(NSSize(width: componentsView.frame.size.width, height: textField.frame.size.height))
                }
            }
                
            if let textView = component as? NSTextView {
                // if textField has an attributed string, it will be autosized
                // otherwize only width is adapted (to fit w/ other components)
                if let textStorage = textView.textStorage {
                    let boundingRect = textStorage.boundingRectWithSize(NSSize(width: componentsView.frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
                    textView.setFrameSize(NSSize(width: componentsView.frame.size.width, height: boundingRect.size.height))
                } else {
                    textView.setFrameSize(NSSize(width: componentsView.frame.size.width, height: textView.frame.size.height))
                }
            }
            
            else if let progressIndicator = component as? NSProgressIndicator {
                progressIndicator.setFrameSize(NSSize(width: componentsView.frame.size.width, height: progressIndicator.frame.size.height))
            }
            
            else if let imageView = component as? NSImageView {
                imageView.setFrameSize(NSSize(width: componentsView.frame.size.width, height: imageView.frame.size.height))
                // TODO: if frame width is downsized, height should be downsized too
                imageView.imageAlignment = .AlignCenter
                imageView.imageScaling = .ScaleProportionallyDown
            }
            
            component.setFrameOrigin(NSPoint(x: 0.0, y: posY))
            posY += component.frame.size.height + vMargin
            
            componentsView.addSubview(component)
        }
        
        posY -= vMargin
        componentsHeightConstraint?.constant = posY
    }
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        if let icon = self.icon {
            iconView?.image = NSImage(named: icon)
        }
        
        refreshComponents()
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
    }
    
    required init?(coder: NSCoder) {
        super.init(coder: coder)
    }
}