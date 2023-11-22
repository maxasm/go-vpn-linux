package main

import (
	"os/exec"
	"errors"
)

func run_cmd(pipein []byte,command string, args []string) (error,[]byte) {
	cmd := exec.Command(command, args...)
	cmd_stdin_pipe, err__get_cmd_stdin_pipe := cmd.StdinPipe()
	if err__get_cmd_stdin_pipe != nil {
		return err__get_cmd_stdin_pipe, nil
	}

	cmd_stdin_pipe.Write(pipein)
	cmd_stdin_pipe.Close()

	// combined output might be separating `stdout` and `stderr` with a new line 
	co, err__run_cmd := cmd.Output()
	if err__run_cmd != nil {
		// TODO: text has a new line at the end?
		// TODO: What error should we return?
		return errors.New(string(co)), co
	}

	return nil, co	
}
