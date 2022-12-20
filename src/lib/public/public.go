package public

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"io"
	"io/fs"
	"strings"
	"text/template"
)

var CachedIntegrity map[string]string = map[string]string{}
var Integrity template.FuncMap = template.FuncMap{}

func GenerateIntegrityMap(publicDir fs.FS) {
	fs.WalkDir(publicDir, ".", func(path string, directory fs.DirEntry, err error) error {
		parts := strings.Split(directory.Name(), ".")
		extension := parts[len(parts)-1]

		if extension == "js" && !directory.IsDir() {
			key := strings.Join(parts[:len(parts)-1], ".") + "PublicIntegrity"

			Integrity[key] = func() string {
				if CachedIntegrity[key] == "" {
					hashAlgo := "sha384"
					hash := sha512.New384()
					buffer := make([]byte, 30*1024)

					file, _ := publicDir.Open(path)
					defer file.Close()

					for {
						n, err := file.Read(buffer)
						if n > 0 {
							hash.Write(buffer[:n])
						}

						if err == io.EOF {
							break
						}
					}

					sum := hash.Sum(nil)
					resultBuffer := bytes.NewBuffer([]byte{})
					base64.NewEncoder(base64.StdEncoding, resultBuffer).Write(sum)
					value := hashAlgo + "-" + resultBuffer.String()

					CachedIntegrity[key] = value

					return value
				} else {
					return CachedIntegrity[key]
				}
			}
		}
		return nil
	})
}
