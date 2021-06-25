package yubikey

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"simple-nginx-otp/utils/rand"
	"strconv"
	"strings"
)

var API_COUNT = 5

func Validate(otp string, id string) bool {
	length := len(otp)
	if length < 32 || length > 48 || !strings.HasPrefix(otp, id) {
		return false
	}
	ch := make(chan bool, API_COUNT)
	for i := 1; i <= API_COUNT; i++ {
		go func(i int) {
			api := ""
			if i > 1 {
				api = strconv.Itoa(i)
			}
			nonce, err := rand.GenerateRandomString(16)
			if err != nil {
				log.Printf("failed to generate random string failed\n%s", err)
				return
			}
			resp, err := http.Get(fmt.Sprintf("https://api%s.yubico.com/wsapi/2.0/verify?id=1&otp=%s&nonce=%s", api, otp, nonce))
			if err != nil {
				log.Printf("failed to contact api %d\n%s", i, err)
				return
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("failed to read api %d body\n%s", i, err)
				return
			}
			err = resp.Body.Close()
			if err != nil {
				log.Printf("failed to close api %d body\n%s", i, err)
				return
			}
			lines := strings.Split(string(body), "\n")
			checks := 0
			for _, line := range lines {
				line = strings.TrimSpace(line)
				delim := strings.Index(line, "=")
				if delim == -1 {
					continue
				}
				key := line[:delim]
				val := line[delim+1:]

				switch key {
				case "otp":
					if val != otp {
						log.Printf("otp mismatch from api %d", i)
						return
					}
					checks++
				case "nonce":
					if val != nonce {
						log.Printf("nonce mismatch from api %d", i)
						return
					}
					checks++
				case "status":
					if val != "OK" {
						ch<-false
						return
					}
					checks++
				}
				if checks == 3 {
					ch<-true
				}
			}
		}(i)
	}
	return <-ch
}
