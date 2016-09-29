//
//  Operators.swift
//  docker-installer
//
//  Created by Michel Courtine on 12/5/15.
//  Copyright © 2015 David Scott. All rights reserved.
//

import Foundation

infix operator ~> {}

private let queue = dispatch_queue_create("docker-worker", DISPATCH_QUEUE_SERIAL)

func ~> (
    backgroundClosure: () -> (),
    mainClosure:       () -> ())
{
    dispatch_async(queue) {
        backgroundClosure()
        dispatch_async(dispatch_get_main_queue(), mainClosure)
    }
}
