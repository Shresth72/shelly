package shelly

import "fmt"

func ParseAST(parts []string) (Node, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("Recieved empty command")
	}

	for i, part := range parts {
		if part != "|" {
			continue
		}

		if i == len(parts)-1 {
			return nil, ErrIncomplete
		}

		if i == 0 {
			return nil, fmt.Errorf("Invalid Pipe")
		}

		left, err := ParseAST(parts[:i])
		if err != nil {
			return nil, err
		}

		right, err := ParseAST(parts[i+1:])
		if err != nil {
			return nil, err
		}

		return &Pipe{
			Left:  left,
			Right: right,
		}, nil
	}

	return parseCommand(parts)
}

func parseCommand(parts []string) (Node, error) {
	cmd := &Command{}

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		switch part {
		case ">", "1>", ">>", "1>>", "2>", "2>>":
			if i+1 >= len(parts) {
				return nil, fmt.Errorf("Missing redirect target")
			}

			target := parts[i+1]
			if isRedirectOperator(target) || target == "|" {
				return nil, fmt.Errorf("Invalid redirect target")
			}

			redirect := Redirect{
				Target: target,
			}

			switch part {
			case ">", "1>":
				redirect.Type = 1
			case ">>", "1>>":
				redirect.Type = 1
				redirect.Append = true
			case "2>":
				redirect.Type = 2
			case "2>>":
				redirect.Type = 2
				redirect.Append = true
			}

			cmd.Redirs = append(cmd.Redirs, redirect)
			i++

		default:
			if cmd.Name == "" {
				cmd.Name = part
			} else {
				cmd.Args = append(cmd.Args, part)
			}
		}
	}

	if cmd.Name == "" {
		return nil, fmt.Errorf("Missing command")
	}

	return cmd, nil
}

func isRedirectOperator(part string) bool {
	switch part {
	case ">", "1>", ">>", "1>>", "2>", "2>>":
		return true
	}
	return false
}
