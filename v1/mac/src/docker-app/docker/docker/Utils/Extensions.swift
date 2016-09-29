//
//  Utils.swift
//  docker
//
//  Created by Adrian Duermael on 12/9/15.
//  Copyright © 2015 docker. All rights reserved.
//

import Foundation
import Cocoa

// Extensions for String
extension String {
    
    // Localize extends String and allows for localization.
    // Localization strings should be placed in Localizable.strings
    func localize(args: CVarArgType...) -> String {
        return withVaList(args) { NSString(format: NSLocalizedString(self, comment: self), arguments: $0) as String }
    }
    
    // returns the first position of a char in the string
    func indexOf(target: String) -> Int? {
        let range = (self as NSString).rangeOfString(target)
        guard range.toRange() != nil else {
            return nil
        }
        return range.location
    }
    
    // returns the last position of a char in the string
    func lastIndexOf(target: String) -> Int? {
        let range = (self as NSString).rangeOfString(target, options: NSStringCompareOptions.BackwardsSearch)
        guard range.toRange() != nil else {
            return nil
        }
        return range.location
    }
    
    // return a substring
    func substringWithRange(range: Range<Int>) -> String {
        let start = self.startIndex.advancedBy(range.startIndex)
        let end = self.startIndex.advancedBy(range.endIndex)
        return self.substringWithRange(start..<end)
    }
    
    // returns a string sha1 digest
    func sha1() -> String? {
        guard let data = self.dataUsingEncoding(NSUTF8StringEncoding) else {
            return nil
        }
        var digest = [UInt8](count:Int(CC_SHA1_DIGEST_LENGTH), repeatedValue: 0)
        CC_SHA1(data.bytes, CC_LONG(data.length), &digest)
        let hexBytes = digest.map { String(format: "%02hhx", $0) }
        return hexBytes.joinWithSeparator("")
    }
    
    // returns a binary sha1 digest
    func sha1_binary() -> NSData? {
        guard let data = self.dataUsingEncoding(NSUTF8StringEncoding) else {
            return nil
        }
        var digest = [UInt8](count:Int(CC_SHA1_DIGEST_LENGTH), repeatedValue: 0)
        CC_SHA1(data.bytes, CC_LONG(data.length), &digest)
        return NSData(bytes: digest, length: digest.count)
    }
}

extension NSMutableAttributedString {
    
    func format(fontSize: CGFloat = 10.0) {
        formatColor(fontSize)
        formatBold(fontSize)
        formatCommands(fontSize)
        formatLinks()
    }
    
    private func getColorFromString(webColorString: String) -> NSColor?
    {
        var result: NSColor? = nil
        var colorCode: UInt32 = 0
        var redByte, greenByte, blueByte: UInt8
        
        let index1 = webColorString.endIndex.advancedBy(-6)
        let substring1 = webColorString.substringFromIndex(index1)
        
        let scanner = NSScanner(string: substring1)
        let success = scanner.scanHexInt(&colorCode)
        
        if success == true {
            redByte = UInt8.init(truncatingBitPattern: (colorCode >> 16))
            greenByte = UInt8.init(truncatingBitPattern: (colorCode >> 8))
            blueByte = UInt8.init(truncatingBitPattern: colorCode) // masks off high bits
            
            result = NSColor(calibratedRed: CGFloat(redByte) / 0xff, green: CGFloat(greenByte) / 0xff, blue: CGFloat(blueByte) / 0xff, alpha: 1.0)
        }
        return result
    }

    
    func formatColor(fontSize: CGFloat = 10.0) {
        let regexp: NSRegularExpression
        
        do {
            regexp = try NSRegularExpression(pattern: "#([^#]{6})#([^{##}]*)##", options: NSRegularExpressionOptions(rawValue: 0))
        } catch {
            return
        }
        
        let all = NSRange(location: 0, length: self.length)
        
        var matches = Array<NSRange>()
        
        regexp.enumerateMatchesInString(self.string, options:  NSMatchingOptions(rawValue: 0), range: all) { (result: NSTextCheckingResult?, _, _) in
            if let range = result?.range {
                matches.append(range)
            }
        }
        
        for match in matches.reverse() {
            let subAttrString = self.attributedSubstringFromRange(match)
            let all = NSRange(location: 0, length: subAttrString.length)
        
            let colorCode = regexp.stringByReplacingMatchesInString(subAttrString.string, options: NSMatchingOptions(rawValue: 0), range: all, withTemplate: "$1")
            
            let replacement = regexp.stringByReplacingMatchesInString(subAttrString.string, options: NSMatchingOptions(rawValue: 0), range: all, withTemplate: "$2")
        
            self.replaceCharactersInRange(match, withString: replacement)
            
            if let color = getColorFromString(colorCode) {
                self.addAttribute(NSForegroundColorAttributeName, value: color, range: NSRange(location: match.location, length: replacement.characters.count))
            }
        }
    }

    
    func formatLinks() {
        let regexp: NSRegularExpression
        
        do {
            regexp = try NSRegularExpression(pattern: "\\[([^]]*)\\]\\(([^)]*)\\)", options: NSRegularExpressionOptions(rawValue: 0))
        } catch {
            return
        }
        
        let all = NSRange(location: 0, length: self.length)
        
        var matches = Array<NSRange>()
        
        regexp.enumerateMatchesInString(self.string, options:  NSMatchingOptions(rawValue: 0), range: all) { (result: NSTextCheckingResult?, _, _) in
            if let range = result?.range {
                matches.append(range)
            }
        }
        
        for match in matches.reverse() {
            let subAttrString = self.attributedSubstringFromRange(match)
            let all = NSRange(location: 0, length: subAttrString.length)
            
            let replacement = regexp.stringByReplacingMatchesInString(subAttrString.string, options: NSMatchingOptions(rawValue: 0), range: all, withTemplate: "$1")
            
            let urlStr = regexp.stringByReplacingMatchesInString(subAttrString.string, options: NSMatchingOptions(rawValue: 0), range: all, withTemplate: "$2")
            let url =  NSURL(string: urlStr)
            
            self.replaceCharactersInRange(match, withString: replacement)
            
            if let url = url {
                self.addAttribute(NSLinkAttributeName, value: url, range: NSRange(location: match.location, length: replacement.characters.count))
            }
        }
    }
    
    func formatBold(fontSize: CGFloat = 10.0) {
        let regexp: NSRegularExpression
        
        do {
            regexp = try NSRegularExpression(pattern: "\\*\\*([^{\\*\\*}]*)\\*\\*", options: NSRegularExpressionOptions(rawValue: 0))
        } catch {
            return
        }
        
        let all = NSRange(location: 0, length: self.length)
        
        var matches = Array<NSRange>()
        
        regexp.enumerateMatchesInString(self.string, options:  NSMatchingOptions(rawValue: 0), range: all) { (result: NSTextCheckingResult?, _, _) in
            if let range = result?.range {
                matches.append(range)
            }
        }
        
        for match in matches.reverse() {
            let subAttrString = self.attributedSubstringFromRange(match)
            let all = NSRange(location: 0, length: subAttrString.length)
            
            let replacement = regexp.stringByReplacingMatchesInString(subAttrString.string, options: NSMatchingOptions(rawValue: 0), range: all, withTemplate: "$1")
            
            self.replaceCharactersInRange(match, withString: replacement)
            
            self.addAttribute(NSFontAttributeName, value: NSFont.boldSystemFontOfSize(fontSize), range: NSRange(location: match.location, length: replacement.characters.count))
        }
    }
    
    func formatCommands(fontSize: CGFloat = 10.0) {
        let regexp: NSRegularExpression
        
        do {
            regexp = try NSRegularExpression(pattern: "`([^`]*)`", options: NSRegularExpressionOptions(rawValue: 0))
        } catch {
            return
        }
        
        let all = NSRange(location: 0, length: self.length)
        
        var matches = Array<NSRange>()
        
        regexp.enumerateMatchesInString(self.string, options:  NSMatchingOptions(rawValue: 0), range: all) { (result: NSTextCheckingResult?, _, _) in
            if let range = result?.range {
                matches.append(range)
            }
        }
        
        for match in matches.reverse() {
            let subAttrString = self.attributedSubstringFromRange(match)
            let all = NSRange(location: 0, length: subAttrString.length)
            
            let replacement = regexp.stringByReplacingMatchesInString(subAttrString.string, options: NSMatchingOptions(rawValue: 0), range: all, withTemplate: "$1")
            
            self.replaceCharactersInRange(match, withString: replacement)
            
            if let font = NSFont(name: "Menlo Regular", size: fontSize + 1.0) {
                self.addAttribute(NSFontAttributeName, value: font, range: NSRange(location: match.location, length: replacement.characters.count))
            }
            
            self.addAttribute(NSBackgroundColorAttributeName, value: NSColor(white: 0.0, alpha: 0.05), range: NSRange(location: match.location, length: replacement.characters.count))
            
            self.addAttribute(NSForegroundColorAttributeName, value: NSColor.darkGrayColor(), range: NSRange(location: match.location, length: replacement.characters.count))
        }
    }
}


extension NSWindow {
    
    func changeView(let view: NSView?, animate: Bool) {
        // if the view argument is nil, we do nothing
        guard let newView: NSView = view else {
            return
        }
        // if the window doesn't have a contentView, we do nothing
        guard let currentView = self.contentView else {
            return
        }
        // we store currentView's size as it will be altered in the process
        // so we can restore it to the correct size at the end of the function
        let currentViewOriginalSize = currentView.frame.size
        
        // store the new view's size because setting the content view will
        // resize the view
        let newViewSize: NSSize = newView.frame.size
        
        let toolbarHeight = self.frame.size.height - currentView.frame.size.height
        let newHeight = newViewSize.height + toolbarHeight - titleBarHeight()
        let newWidth  = newViewSize.width
        
        var aFrame = NSWindow.contentRectForFrameRect(self.frame, styleMask: self.styleMask)
        aFrame.origin.y += aFrame.size.height
        aFrame.origin.y -= newHeight
        aFrame.size.height = newHeight
        aFrame.size.width  = newWidth
        aFrame = NSWindow.frameRectForContentRect(aFrame, styleMask: self.styleMask)
        
        // populate the window with the new view
        self.contentView = newView
        
        // resize the windows, eventually using an animation
        self.setFrame(aFrame, display: true, animate: animate)

        // fix the current view's (previous, now) size
        currentView.frame.size = currentViewOriginalSize
    }
    
    private func titleBarHeight() -> CGFloat
    {
        let frame = NSRect(x:0, y:0, width: 100, height: 100)
        let contentRect = NSWindow.contentRectForFrameRect(frame, styleMask: NSTitledWindowMask)
        return frame.size.height - contentRect.size.height
    }
    
}

extension CFArray: SequenceType {
    public func generate() -> AnyGenerator<AnyObject> {
        var index = -1
        let maxIndex = CFArrayGetCount(self)
        return AnyGenerator {
            index += 1
            guard index < maxIndex else {
                return nil
            }
            let unmanagedObject: UnsafePointer<Void> = CFArrayGetValueAtIndex(self, index)
            let rec = unsafeBitCast(unmanagedObject, AnyObject.self)
            return rec
        }
    }
}
