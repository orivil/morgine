// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			log.Println(err)
		}
	}()
}
