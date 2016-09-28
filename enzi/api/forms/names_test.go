package forms

import (
	"testing"
)

func TestValidateAccountNames(t *testing.T) {
	validNames := []string{
		"alice",                 // Simple alphabetical name.
		"b0b",                   // With numbers.
		"jane.doe",              // With period.
		"john-doe",              // With hyphen.
		"carmen_san_diego",      // With underscore.
		"_charlie_xcx_",         // Starts with underscore.
		"$creative_name_here",   // With dollar sign.
		"micro$oft",             // more dollar sign.
		"user.name@example.com", // Email address.
		"1337-h@ckerz",          // Numbers, hyphen, at-sign.
		"김일성",                   // Kim Il-Sung
		"김정일",                   // Kim Jong-Il
		"김정은",                   // Kim Jong-Un
		"张伟",                    // Zhang Wei
		"张丽",                    // Zhang Li
		"李桂英",                   // Li Gui Ying
		"Александр",             // Alexander
		"Матвей",                // Matvei
		"Юрий",                  // Yury
	}

	for _, validName := range validNames {
		if err := ValidateAccountName(&validName, "name"); err != nil {
			t.Fail()
			t.Errorf("account name %q should be valid: %s", validName, err)
		}
	}

	invalidNames := []string{
		"", // Too short.
		"abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789", // Too long.
		"slim/shady",    // Contains a '/'.
		"back\\slash",   // Contains a '\'.
		"left[bracket",  // Contains a '['.
		"right]bracket", // Contains a ']'.
		"colin:powell",  // Contains a ':'.
		"semi;colon",    // Contains a ';'.
		"pip|pip",       // Contains a '|'.
		"equal=sign",    // Contains a '='.
		"c,s,v",         // Contains a ','.
		"1+1",           // Contains a '+'.
		"a*search",      // Contains a '*'.
		"wut?wut",       // Contains a '?'.
		"left<arrow",    // Contains a '<'.
		"right>arrow",   // Contains a '>'.
	}

	for _, invalidName := range invalidNames {
		if err := ValidateAccountName(&invalidName, "name"); err == nil {
			t.Fail()
			t.Errorf("account name %q should be invalid", invalidName)
		}
	}
}
