# nativemessaging
Native messaging host library for go application  

## Usage

**Go host application**

``` go
package main

import (
	"io"
	"os"

	"github.com/qrtz/nativemessaging"
)

func main() {
	// Messaging host with native byte order
	host := nativemessaging.NativeHost(os.Stdin, os.Stdout)

	for {
		var rsp response
		var msg message
		err := host.Receive(&msg)

		if err != nil {
			if err == io.EOF {
				// exit
				return
			}
			rsp.Text = err.Error()
		} else {
			if msg.Text == "ping" {
				rsp.Text = "pong"
				rsp.Success = true
			} else {
				// Echo the message back to the client
				rsp.Text = msg.Text
			}
		}

		if _, err := host.Send(rsp); err != nil {
			// Log the error and exit
			return
		}
	}
}

type message struct {
	Text string `json:"text"`
}

type response struct {
	Text    string `json:"text"`
	Success bool   `json:"success"`
}
```

**Javascript client**

``` js
    chrome.runtime.sendNativeMessage('com.github.qrtz.nativemessaginghost', {text:'ping'}, (response) => {
        console.log('Native messaging host response ', response);
    })
```

**More info:**  

https://developer.chrome.com/extensions/nativeMessaging  
https://developer.mozilla.org/en-US/Add-ons/WebExtensions/Native_messaging
