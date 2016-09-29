//
//  checksum.swift
//  nettest
//
//  Created by Magnus Skjegstad on 01/01/2016.
//  Copyright © 2016 Magnus Skjegstad. All rights reserved.
//

import Foundation

// Calculate 16 bit IP checksum
func checksum16(inout data : [UInt8]) -> UInt16
{
    if data.count % 2 != 0 { // exit if not even count
        return 0
    }
    var sum : UInt32 = 0 // leave room for overflow
    for var i = 0; i < data.count; i += 2 {
        let b = UInt32((UInt16(data[i]) << 8) + UInt16(data[i+1]))
        sum = (sum + b) & 0xffff // truncate bits at 16
        if sum < b { // check overflow
            sum++ // add carry
        }
    }
    return ~UInt16(sum & 0xffff)
    
    //NOTE: Swift does not automatically truncate bits on type conversion and will crash on overflow.
    //See e.g. https://medium.com/ios-os-x-development/swift-integer-overflow-issue-2970f3896f59
}

