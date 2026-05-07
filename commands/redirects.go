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

type redirectMode struct {
	append bool
	stdout bool
	stderr bool
}

var redirectModes = map[string]redirectMode{
	">":  {append: false, stdout: true, stderr: false},
	"1>": {append: false, stdout: true, stderr: false},
	"2>": {append: false, stdout: false, stderr: true},

	">>":  {append: true, stdout: true, stderr: false},
	"1>>": {append: true, stdout: true, stderr: false},
	"2>>": {append: true, stdout: false, stderr: true},
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

		mode, ok := redirectModes[part]
		if !ok {
			args = append(args, part)
			continue
		}

		file, err := openRedirectFile(parts, i, &files, mode.append)
		if err != nil {
			return nil, err
		}

		if mode.stdout {
			r.Stdout = file
		}
		if mode.stderr {
			r.Stderr = file
		}

		i++
	}

	r.Args = args
	r.Close = func() {
		for _, f := range files {
			f.Close()
		}
	}

	return r, nil
}

func openRedirectFile(
	parts []string,
	i int,
	files *[]*os.File,
	appendMode bool,
) (*os.File, error) {
	var file *os.File
	var err error

	if i+1 >= len(parts) {
		return nil, fmt.Errorf("missing redirect target")
	}

	if appendMode {
		file, err = os.OpenFile(
			parts[i+1],
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
	} else {
		file, err = os.Create(parts[i+1])
	}

	if err != nil {
		return nil, err
	}

	*files = append(*files, file)
	return file, nil
}
