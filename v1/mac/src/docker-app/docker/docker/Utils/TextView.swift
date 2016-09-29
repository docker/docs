//
//  TextView.swift
//  Docker
//
//  Created by Doby Mock on 5/26/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

// Note(aduermael): Displaying rich text using UI components with dynamic
// height, well integrated with NSLayoutConstraints is a real pain...
// Different implementations below. So far the best options seems to be
// TextView when text areas are added by code, but TextViewInView when 
// using interface builder.

class TextView: NSTextView {
    
    private var heightLayoutConstraint: NSLayoutConstraint? = nil
    
    //link: AnyObject, atIndex charIndex: Int
    var onClickedLink: ((AnyObject, Int) -> Bool) = { (a, b) in true }
    
    override func awakeFromNib() {
        super.awakeFromNib()
        
        for constrain in constraints {
            if let firstItemObject = constrain.firstItem as? NSObject {
                if firstItemObject === self &&
                    constrain.firstAttribute == .Height &&
                    constrain.secondItem == nil &&
                    constrain.secondAttribute == .NotAnAttribute &&
                    constrain.relation == .Equal
                {
                    heightLayoutConstraint = constrain
                }
            }
        }
        
        // create height constraint if non existant
        if heightLayoutConstraint == nil {
            Logger.log(level: Logger.Level.Warning, content: "⚠️ it is highly recommended to give any TextView an height layout constraint! ⚠️")
        }
    }
    
    var text: String {
        get {
            return attributedString().string
        }
        set(newText) {
            let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
            let attributedString = NSMutableAttributedString(string: newText, attributes:attrs)
            attributedString.format(NSFont.systemFontSize())
            textStorage?.setAttributedString(attributedString)
            fitHeight()
        }
    }
    
    private func fitHeight() {
        let maxHeight: CGFloat = 10000000
        if let textStorage = textStorage {
            let boundingRect = textStorage.boundingRectWithSize(NSSize(width: frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
            let size = NSSize(width: frame.size.width, height: boundingRect.size.height)
            setFrameSize(NSSize(width: frame.size.width, height: size.height))
            heightLayoutConstraint?.constant = size.height
        }
    }
    
    override func clickedOnLink(link: AnyObject, atIndex charIndex: Int) {
        if onClickedLink(link, charIndex) {
            super.clickedOnLink(link, atIndex: charIndex)
        }
    }
}



class TextViewInView: NSView {
    
    private var textView: TextView! = nil
    private var heightLayoutConstraint: NSLayoutConstraint? = nil
    
    required init?(coder: NSCoder) {
        super.init(coder: coder)
        
        for constraint in constraints {
            if let firstItemObject = constraint.firstItem as? NSObject {
                if firstItemObject === self &&
                    constraint.firstAttribute == .Height &&
                    constraint.secondItem == nil &&
                    constraint.secondAttribute == .NotAnAttribute &&
                    constraint.relation == .Equal
                {
                    heightLayoutConstraint = constraint
                }
            }
        }
        
        // create height constraint if non existant
        if heightLayoutConstraint == nil {
            Logger.log(level: Logger.Level.Warning, content: "⚠️ it is highly recommended to give any TextView an height layout constraint! ⚠️")
        }
        
        textView = TextView()
        textView.drawsBackground = false
        textView.selectable = true
        textView.editable = false
        textView.textContainer?.lineFragmentPadding = 0
        textView.textContainerInset = NSSize(width: 0, height: 0)
        textView.setFrameSize(frame.size)
        textView.setFrameOrigin(NSPoint(x: 0, y: 0))
        
        self.addSubview(textView)
    }
    
    var text: String {
        get {
            return textView.attributedString().string
        }
        set(newText) {
            let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
            let attributedString = NSMutableAttributedString(string: newText, attributes:attrs)
            attributedString.format(NSFont.systemFontSize())
            textView.textStorage?.setAttributedString(attributedString)
            fitHeight()
            
        }
    }
    
    private func fitHeight() {
        let maxHeight: CGFloat = 10000000
        if let textStorage = textView.textStorage {
            let boundingRect = textStorage.boundingRectWithSize(NSSize(width: frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
            let size = NSSize(width: frame.size.width, height: boundingRect.size.height)
            
            setFrameSize(NSSize(width: frame.size.width, height: size.height))
            heightLayoutConstraint?.constant = size.height
            textView.setFrameSize(size)
            textView.setFrameOrigin(NSPoint(x: 0, y: 0))
        }
    }
    
//    var onClickedLink: ((AnyObject, Int) -> Bool) = { (a, b) in true }
    var onClickedLink: ((AnyObject, Int) -> Bool) {
        get {
            return textView.onClickedLink
        }
        set {
            textView.onClickedLink = newValue
        }
    }

}



// For some reasonm it's impossible to drop NSTextViews using interface builder.
// NSTextViews are always wrapped like this:
// NSScrollView > NSClipView > NSTextView
// FlexibleHeightTextView allows to create NSTextViews with interface builder,
// that are wrapped, but disabling scroll (both vertically and horizontally).
// the width always remains the same, but the height is dynamic to fit content.
class FlexibleHeightTextView: NSScrollView {
    
    private var structureOk = false
    private var textView: NSTextView! = nil
    private var clipView: NSClipView! = nil
    private var heightLayoutConstraint: NSLayoutConstraint? = nil
    
    // check that we actually have the NSScrollView > NSClipView > NSTextView
    // structure
    override func awakeFromNib() {
        super.awakeFromNib()
        
        for subview in subviews {
            if subview is NSScroller {
                // ok, but nothing to do
            } else if let cv = subview as? NSClipView {
                clipView = cv
            } else {
                // NSScrollView should only contain NSScroller and one 
                // NSClipView, this outlet can't be used as FlexibleHeightTextView
                Logger.log(level: ASLLogger.Level.Error, content: "FlexibleHeightTextView awakeFromNib error")
                return
            }
        }
        
        // detect existing height constraint
        for constrain in constraints {
            if let firstItemObject = constrain.firstItem as? NSObject {
                if firstItemObject === self &&
                    constrain.firstAttribute == .Height &&
                    constrain.secondItem == nil &&
                    constrain.secondAttribute == .NotAnAttribute &&
                    constrain.relation == .Equal
                {
                    heightLayoutConstraint = constrain
                }
            }
        }
        
        // create height constraint if non existant
        if heightLayoutConstraint == nil {
           Logger.log(level: Logger.Level.Warning, content: "⚠️ it is highly recommended to give any FlexibleHeightTextView an height layout constraint! ⚠️")
        }

        guard clipView.subviews.count == 1 else { return }
        
        guard let tv = clipView.subviews[0] as? NSTextView else { return }
        
        textView = tv
        
        Swift.print("clipView \(clipView)")
        
        self.hasVerticalScroller = false
        self.hasHorizontalRuler = false
        self.verticalScrollElasticity = .None
        self.horizontalScrollElasticity = .None
        
        textView.drawsBackground = false
        textView.editable = false
        textView.selectable = true
        textView.textContainerInset = NSSize(width: 0, height: 0)
        
        structureOk = true
    }
    
    
    var text: String {
        get {
            guard structureOk else { return "" }
            return textView.attributedString().string
        }
        set(newText) {
            guard structureOk else { return }
            let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
            let attributedString = NSMutableAttributedString(string: newText, attributes:attrs)
            attributedString.format(NSFont.systemFontSize())
            textView.textStorage?.setAttributedString(attributedString)
            fitHeight()
        }
    }
    
    private func fitHeight() {
        guard structureOk else { return }
        let maxHeight: CGFloat = 10000000
        if let textStorage = textView.textStorage {
            let boundingRect = textStorage.boundingRectWithSize(NSSize(width: frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
            let size = NSSize(width: frame.size.width, height: boundingRect.size.height)
            self.setFrameSize(size)
            heightLayoutConstraint?.constant = size.height
        }
    }
}

