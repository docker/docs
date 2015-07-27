// package compact provides function to work with json compact serialization format
package compact

import (
	"strings"
	"github.com/dvsekhvalnov/jose2go/base64url"
)

// Parse splitting & decoding compact serialized json web token, returns slice of byte arrays, each representing part of token
func Parse(token string) [][]byte {
	parts:=strings.Split(token,".")
	
	result:=make([][]byte,len(parts))
	var e error		
	
	for i,part:=range parts	{
		if result[i],e=base64url.Decode(part);e!=nil {
			panic(e)
		}		
	}
		
	return result
}

// Serialize converts given parts into compact serialization format
func Serialize(parts ...[]byte) string {
	result:=make([]string,len(parts))
	
	for i,part:=range parts {
		result[i]=base64url.Encode(part)
	}
	
	return strings.Join(result,".")
}
