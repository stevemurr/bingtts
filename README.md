# Installation

`go get -u github.com/stevemurr/bingtts`

# Example

```
package main

import (
    "log"
    "io/ioutil"

    "github.com/stevemurr/bingtts"
)
func main() {
    // See what voices you can use
    for key, value := range bingtts.GetVoices() {
        log.Printf("%s: %s", key, value)
    }
    // Pass in your key from https://www.microsoft.com/cognitive-services/en-us/sign-up
    // Key lasts 10 minutes so if you're doing serious synthesis, write some code to reuse the key
    key := "YOUR KEY HERE"
    token, err := bingtts.IssueToken(key)
    if err != nil {
        log.Println(err)
    }
    // Synthesize
    res, err := bingtts.Synthesize(
        token,
        "oh my god did you hear?  trump is dead!",
        "es-mx",
        "male",
        bingtts.RIFF16Bit16kHzMonoPCM)
    if err != nil {
        log.Println(err)
    }
    // res is []byte
    ioutil.WriteFile("test.wav", res, 0644)
}
```
