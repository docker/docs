//
//  TextField.swift
//  docker
//
//  Created by Doby Mock on 5/26/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class TextField: NSTextField {
    
    var text: String {
        get {
            return self.attributedStringValue.string
        }
        set(newText) {
            let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
            let attributedString = NSMutableAttributedString(string: newText, attributes:attrs)
            attributedString.format(NSFont.systemFontSize())
            self.attributedStringValue = attributedString
        }
    }
    
    func fitHeight(heightLayoutConstraint: NSLayoutConstraint?) {
        let maxHeight: CGFloat = 10000
        let boundingRect = attributedStringValue.boundingRectWithSize(NSSize(width: frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
        setFrameSize(NSSize(width: frame.size.width, height: boundingRect.size.height))
        if let heightLC = heightLayoutConstraint {
            heightLC.constant = boundingRect.size.height
        }
    }
    
}
