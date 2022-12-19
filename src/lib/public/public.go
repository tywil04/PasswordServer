package public

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"io"
	"io/fs"
	"strings"
)

var Integrity map[string]string = map[string]string{}

func CalculatePublicJSIntegrity(publicDir fs.FS) {
	fs.WalkDir(publicDir, ".", func(path string, directory fs.DirEntry, err error) error {
		parts := strings.Split(directory.Name(), ".")
		extension := parts[len(parts)-1]

		if extension == "js" && !directory.IsDir() {
			buffer := make([]byte, 30*1024)
			sha384 := sha512.New384()

			file, _ := publicDir.Open(path)
			defer file.Close()

			for {
				n, err := file.Read(buffer)
				if n > 0 {
					sha384.Write(buffer[:n])
				}

				if err == io.EOF {
					break
				}
			}

			sum := sha384.Sum(nil)
			resultBuffer := bytes.NewBuffer([]byte{})
			base64.NewEncoder(base64.StdEncoding, resultBuffer).Write(sum)

			key := strings.Join(parts[:len(parts)-1], ".")
			value := "sha384-" + resultBuffer.String()

			Integrity[key] = value
		}

		return nil
	})
}
