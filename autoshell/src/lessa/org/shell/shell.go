package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
    uuid "github.com/nu7hatch/gouuid"
	"time"
	"unicode"
)

type MessageType int

const (
	MT_TEXT		MessageType = iota
	MT_ERROR
	MT_LOG
)

type ShellMessage struct {
	Type	MessageType
	Data	string
}

type ShellMessageListener func (message *ShellMessage)

type dataFragment struct {
	mtype	MessageType
	data	string
	flush	bool
}

type AutoShell struct {
	barrier			*sync.WaitGroup
	fragments		chan *dataFragment
	fragmentsDone	chan bool
	inData			chan string
	listener		ShellMessageListener
	signature		string
	syncSignal		string
	syncStdout		string
	syncStderr		string
}

const (
	commandDone		= "done"
	commandExit		= "exit"
	commandQuit		= "quit"
	prompt			= "\nauto-shell > "
	shellCommand	= "bash"
)

func New(listener ShellMessageListener) *AutoShell {

	var siguuid, _ = uuid.NewV4()
	sig := siguuid.String()
	
	result := &AutoShell {
		barrier:		&sync.WaitGroup{},
		fragments:		make(chan *dataFragment),
		fragmentsDone:	make(chan bool),
		inData:			make(chan string),
		listener:		listener,
		signature:		siguuid.String(),
		syncSignal:		fmt.Sprintf("%s\n", sig),
		syncStdout:		fmt.Sprintf("echo \"%s\"\n", sig),
		syncStderr:		fmt.Sprintf(">&2 echo \"%s\"\n", sig),
	}	

	return result
}

func (as *AutoShell) error(message string) {

	as.fragments<-&dataFragment{MT_ERROR, fmt.Sprintf("[auto-shell] Error: %s\n", message), true}
}

func (as *AutoShell) log(message string) {

	fmt.Printf("[Info] %v %s\n", time.Now().Format(time.RFC3339), message)
}

func (as *AutoShell) message(message string, flush bool) {

	as.fragments<-&dataFragment{MT_TEXT, message, flush}
}

func (as *AutoShell) messagePrompt() {

	as.message(prompt, true)
}

func (as *AutoShell) shouldTerminate(command string) bool {

	return command == commandDone || command == commandExit || command == commandQuit
}

func (as *AutoShell) Run() {
	
	// buffer outputs and sends them out when a flush is requested
	as.barrier.Add(1)
	go as.bufferShellMessage()

	// spawn the shell and wait for its completion
	as.barrier.Add(1)
	go as.spawnAndWaitForShellProcess()

	// push a prompt to the output
	as.messagePrompt()

	var cmd string	
	scanner := bufio.NewScanner(os.Stdin)

	for {

		// break if scan fails
		if !scanner.Scan() {
			as.error(fmt.Sprintf("\nInvalid scan: ", scanner.Err()))
			break
		}

		// scanned text
		cmd = strings.TrimRightFunc(scanner.Text(), unicode.IsSpace)

		// ignore this one
		if cmd == "" {
			continue
		}

		// push the command to the channel
		as.inData<-cmd

		// break if the user is done
		if as.shouldTerminate(cmd) {
			as.log("Exiting the main command loop.")
			break;
		}
	}

	as.log("Closing input channel.")
	close(as.inData)

	as.log("Closing fragments channel.")
	as.fragmentsDone<-true

	close(as.fragments)
	close(as.fragmentsDone)

	as.log("Waiting for goroutines.")
	as.barrier.Wait()

	as.log("Cleanup done.")
}

func (as *AutoShell) spawnAndWaitForShellProcess() {

	// signal when we're done
	defer as.barrier.Done()

	// prepare the shell
	shell := exec.Command(shellCommand)

	// shell input is our writer
	stdin, err := shell.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// shell output is our reader
	stdout, err := shell.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// shell error is our reader
	stderr, err := shell.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// start is non-blocking
	if err := shell.Start(); err != nil {
		log.Fatal(err)
	}

	// run commands by writing on the process' stdin
	as.barrier.Add(1)
	go as.runCommands(stdin)

	// drain stdout to this process' stdin
	as.barrier.Add(1)
	go as.drainOutput(stdout, 64*1024, MT_TEXT)

	// drain stderr to this process' stdin
	as.barrier.Add(1)
	go as.drainOutput(stderr, 8*1024, MT_ERROR)

	// wait for the shell to exit
	if err := shell.Wait(); err != nil {
		log.Fatal(err)
	}

	as.log("Shell process exited.")

	// postcondition: stdin, stdout, and stderr must be closed
}

// exists on a channel close or an exit command from the channel
func (as *AutoShell) runCommands(stdin io.WriteCloser) {

	// signal when we're done
	defer as.barrier.Done()

	// listens to the input channel and pushes commands to the shell
	// the input channel owner must close it after the shell exits
	for {
		select {
		case cmd := <-as.inData:

			// terminate interactive shell
			if as.shouldTerminate(cmd) {

				as.log("Exiting interactive shell.")
				io.WriteString(stdin, "exit\n")

				as.log("Removing case from select in command loop.")
				as.inData = nil

				as.log("Breaking out from command loop.")
				return;
			}

			// extend the received command with a new line terminator
			cmdLine := fmt.Sprintf("%s\n", cmd)

			// push the command to the shell's stdin
			io.WriteString(stdin, cmdLine)

			// write sync to stderr to signal end of output
			io.WriteString(stdin, as.syncStderr)

			// write sync to stdout to signal end of output
			io.WriteString(stdin, as.syncStdout)
		}
	}
}

// exists on a read error; typically, an EOF error happens when the reader closes, 
// which happens when the shell exits, so an EOF error is part of the protocol
func (as *AutoShell) drainOutput(reader io.Reader, sz int, mtype MessageType) {

	// signal when we're done
	defer as.barrier.Done()

	buf := make([]byte, sz)

	for {
		// read output chunks in buffer size units; block if no bytes remain
		n, err := reader.Read(buf)

		// copy output "as-is"
		if n > 0 {
			
			out := string(buf[:n])
			hasSuffix := strings.HasSuffix(out, as.syncSignal)

			if hasSuffix {
				//as.log("hasSuffix")
				out = string(buf[:n-len(as.syncSignal)])
			}

			// write only if there is some payload
			if len(out) > 0 || hasSuffix {
				//as.log("writing")
				as.fragments<-&dataFragment{mtype, out, hasSuffix}
			}

			// udpate the prompt
			if hasSuffix && mtype == MT_TEXT {
				as.messagePrompt()
			}
		}

		// for a non-empty last fragment of output, add a newline if necessary
		if n > 0 && n < sz && buf[n-1] != '\n' {
			as.fragments<-&dataFragment{mtype, "\n", true}
		}

		if err != nil {
			break
		}
	}

	as.log(fmt.Sprintf("Done draining output for %v.", mtype))
}

func (as *AutoShell) bufferShellMessage() {

	// signal when we're done
	defer as.barrier.Done()

	// create local buffers for error and output messages
	errBuffer := &bytes.Buffer{}
	outBuffer := &bytes.Buffer{}

	// variable to use as a selector for the above buffers
	var buffer *bytes.Buffer

	for {
		select {
		case fragment := <-as.fragments:

			// select the buffer to be written			
			if fragment.mtype == MT_TEXT {
				buffer = outBuffer
			} else if fragment.mtype == MT_ERROR {
				buffer = errBuffer
			} 
			
			// write to the buffer
			buffer.WriteString(fragment.data)

			// send the fragment if the flush flag is set
			if fragment.flush {
				message := buffer.String()
				buffer.Reset()
				as.handleShellMessage(&ShellMessage{fragment.mtype, message})
			}
		case <-as.fragmentsDone:
			return
		}
	}

	as.log("Done buffering shell messages.")
}

func (as *AutoShell) handleShellMessage(message *ShellMessage) {
	
	if as.listener != nil {
		as.listener(message)
	} else {
		fmt.Print(message.Data)
	}
}