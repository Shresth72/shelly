package shelly

import (
	"fmt"
	"os"
)

type Redirect struct {
	Type   int
	Target string
	Append bool
}

func applyRedirects(
	redirs []Redirect,
	sh *Shell,
) (*Shell, func(), error) {
	ctx := *sh
	var files []*os.File

	cleanup := func() {
		for _, f := range files {
			_ = f.Close()
		}
	}

	for _, redir := range redirs {
		flags := os.O_CREATE | os.O_WRONLY

		if redir.Append {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		file, err := os.OpenFile(redir.Target, flags, 0644)
		if err != nil {
			return nil, func() {}, err
		}

		files = append(files, file)

		switch redir.Type {
		case 1:
			ctx.Stdout = file
		case 2:
			ctx.Stderr = file
		default:
			cleanup()
			return nil, func() {}, fmt.Errorf("Unsupported redirect type: %d", redir.Type)
		}
	}

	return &ctx, cleanup, nil
}
