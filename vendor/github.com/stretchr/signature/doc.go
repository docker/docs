/*
Signature provides URL signing capabilities for Go.

Encoding

  * Generate the original request URL
  * Add ~key public key value

To generate hash:

  * Create a copy of the request URL
  * Add ~private key parameter
  * Add ~bodyhash value containing an SHA1 hash of the body contents if there is a body - otherwise, skip this step
  * Order parameters alphabetically
  * Prefix it with the HTTP method (in uppercase) followed by an ampersand (i.e. "GET&http://...")
  * Hash it (using SHA-1)
  * Add the hash as ~sign to the END of the original URL

Decoding

  * Strip off the ~sign parameter (and keep it)
  * Lookup the account (using ~key) and get the ~private parameter, and add it to the URL
  * Hash it
  * Compare the generated hash with the ~sign value to decide if it the request is valid or not
*/
package signature
