//
//  WizardTests.swift
//  docker
//
//  Created by Adrien Duermael on 2/4/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation

extension Wizard {
    // test displays Wizard messages, to test different functionalities
    
    static func test() {
        
        guard let msg1 = WizardMessageGeneric(message: "test", details: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", icon: "AboutIcon") else {return}
        msg1.addCloseButton("close")
        
        guard let msg2 = WizardMessageGeneric(message: "test", details: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", icon: "AboutIcon") else {return}
        msg2.addCloseButton("close")
        msg2.addText("Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.")
        
        guard let msg3 = WizardMessageGeneric(message: "test", details: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", icon: "AboutIcon") else {return}
        msg3.addProgressBar()
        
        show(msg1, msg2, msg3)
    }
}