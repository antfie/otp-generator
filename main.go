package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz:-_%#?@*!+"

func main() {
	_, err := os.Stat("otp.html")

	if !os.IsNotExist(err) {
		log.Fatal("otp.html exists.")
	}

	if len(os.Args) < 2 {
		log.Fatal("Need to specify a dimension e.g. 5 for 5x5 OTP.")
	}

	dimension, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Dimension must be an integer e.g. 5 for 5x5 OTP.")
	}

	otp := generateRandomCode(dimension*dimension, alphabet)

	// "X" is used for no style/color. This "XXX" trickery is done to keep style and colour occurrences down
	styles := generateRandomCode(dimension*dimension, "XXXXXXX1XXXXXXX2XXXXXXX")
	colors := generateRandomCode(dimension*dimension, "XXXX1XXXX2XXXX3XXXX4XXX")

	body := generateOTPHTMLTable(dimension, otp, styles, colors)
	output := strings.Replace(loadHTMLTemplate(), " {{todo}}", body, 1)

	err = os.WriteFile("otp.html", []byte(output), 0600)

	if err != nil {
		log.Fatal(err)
	}
}

func generateOTPHTMLTable(dimension int, otp, styles, colors string) string {
	var body = ""
	var otpIndex = 0

	for range dimension {
		body += "<tr>"

		for range dimension {
			var class = ""

			style := string(styles[otpIndex])

			if style != "X" {
				class += " s" + style
			}

			color := string(colors[otpIndex])

			if color != "X" {
				class += " c" + color
			}

			if len(class) > 0 {
				class = fmt.Sprintf(" class=\"%s\"", strings.TrimSpace(class))
			}

			body += fmt.Sprintf("<td%s>%s</td>", class, string(otp[otpIndex]))
			otpIndex++
		}

		body += "</tr>"
	}

	return body
}

func generateRandomCode(length int, alphabet string) string {
	b := make([]byte, length)
	r, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	if r != len(b) {
		log.Fatal(errors.New("could not generate desired length of random data"))
	}

	var out []string

	for i := 0; i < length; i++ {
		k := 0
		for j := 0; j < int(b[i]); j++ {
			k++
			if k > len(alphabet)-1 {
				k = 0
			}
		}

		out = append(out, string(alphabet[k]))
	}

	return strings.Join(out, "")
}

func loadHTMLTemplate() string {
	htmlData, err := os.ReadFile(path.Clean("template.html"))

	if err != nil {
		log.Fatal("Could not load template")
	}

	return string(htmlData)
}
