package commands

import (
	"fmt"
	"io"
	"os"
)

type Redirects struct {
	Args   []string
	Stdout io.Writer
	Stderr io.Writer
	Close  func()
}

func parseRedirects(
	parts []string,
	defaultStdout, defaultStderr io.Writer,
) (*Redirects, error) {
	r := &Redirects{
		Stdout: defaultStdout,
		Stderr: defaultStderr,
	}

	var args []string
	var files []*os.File

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		switch part {
		case ">", "1>":
			if i+1 >= len(parts) {
				return nil, fmt.Errorf("missing redirect target")
			}

			file, err := os.Create(parts[i+1])
			if err != nil {
				return nil, err
			}

			r.Stdout = file
			files = append(files, file)
			i++

		default:
			args = append(args, part)
		}
	}

	r.Args = args
	r.Close = func() {
		for _, f := range files {
			f.Close()
		}
	}

	return r, nil
}
