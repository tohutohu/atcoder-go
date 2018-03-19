package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets9e4e8890dcdbda29f610fc6d38be21150dec3d3e = "package main\n\nimport (\n\t\"io/ioutil\"\n\t\"os\"\n\t\"testing\"\n)\n\nfunc Test_solve(t *testing.T) {\n\tsampleInput := []string{\"sampleInput-placeholder\"}\n\tsampleOutput := []string{\"sampleOutput-placeholder\"}\n\tfor i, input := range sampleInput {\n\t\tstdout, _ := stubIO(input, solve)\n\t\tif got, want := stdout, sampleOutput[i]; got != want {\n\t\t\tt.Fatalf(\"wrong answer: got %s, want %s\", got, want)\n\t\t}\n\t}\n\n}\n\nfunc stubIO(inbuf string, fn func()) (string, string) {\n\tinr, inw, _ := os.Pipe()\n\toutr, outw, _ := os.Pipe()\n\terrr, errw, _ := os.Pipe()\n\n\torgStdin := os.Stdin\n\torgStdout := os.Stdout\n\torgStderr := os.Stderr\n\n\tinw.Write([]byte(inbuf))\n\tinw.Close()\n\tos.Stdin = inr\n\tos.Stdout = outw\n\tos.Stderr = errw\n\tfn()\n\tos.Stdin = orgStdin\n\tos.Stdout = orgStdout\n\tos.Stderr = orgStderr\n\toutw.Close()\n\toutbuf, _ := ioutil.ReadAll(outr)\n\terrw.Close()\n\terrbuf, _ := ioutil.ReadAll(errr)\n\n\treturn string(outbuf), string(errbuf)\n}\n"
var _Assetsfb9734e22b388461e045c5d6206dd76539df1895 = "package main\n\nimport \"fmt\"\n\nfunc solve() {\n\tfmt.Scan()\n}\n\nfunc main() {\n\tsolve()\n}\n\nfunc min(a, b int) int {\n\tif a < b {\n\t\treturn a\n\t}\n\treturn b\n}\n\nfunc max(a, b int) int {\n\tif a > b {\n\t\treturn a\n\t}\n\treturn b\n}\n\nfunc abs(a int) int {\n\tif a < 0 {\n\t\treturn -a\n\t}\n\treturn a\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"main_test.go", "main.go"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520768711, 1520768711666417825),
		Data:     nil,
	}, "/main_test.go": &assets.File{
		Path:     "/main_test.go",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520530588, 1520530588936210392),
		Data:     []byte(_Assets9e4e8890dcdbda29f610fc6d38be21150dec3d3e),
	}, "/main.go": &assets.File{
		Path:     "/main.go",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520768711, 1520768711690418020),
		Data:     []byte(_Assetsfb9734e22b388461e045c5d6206dd76539df1895),
	}}, "")
