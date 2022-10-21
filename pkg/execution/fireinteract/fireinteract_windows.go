//go:build implant && (!lp || !teamserver)

package fireinteract

import (
	"bytes"
	"io"
	"shlyuz/pkg/component"
)

// Module loaded, executed, and unloaded when execution ends. Can be with a return code, OR when directed by the loader.
// Provides an IO pipe that can be used. Input is typically generated by the operator
// Can read from pipe, and write responses to the same pipe.
//
//	permits memory injection of code which may continue execution beyond the return while providing output to and receiving input from the user. also permits the loader to trigger the module to quit
//	Returns a null-terminated windows named pipe, suitable for use with WriteFile, ReadFile, etc functions. Format of data is controlled by module's author. This mode supports input. Input is generated by the teamserver, which will send data over a unix pipe to the loader's C2 component on the teamserver. Loader's component then conveys the information to the loader which will then convey data to the pipe on the reomte computer for the module's ingestion. Allows a module to tunnel its communications over the loader rather than opening a new socket
//	  Loader is responsible for closing this pipe if it must exit prior to the module. Module must dtect this closure and cease reading from the pipe in that situation. The loader is also responsible for closing the pipe after module execution ends.
//	  Loader's c2 component will not interact in any way with the data sent over this pipe except to apss it between module's c2 component and the module.
//	MUSt create its own long lived thread rather than rely on any threads created by the loader to avoid a possible access violation.
//
// Behavior is designed to allow the loader to free the memory of the module it is initally executed from. The a module's behavior involves process migration or other execution migration methods that permit the inital execution memolry to be deallocated, then the module does not need to communicate status in this way.
// Pipe Comms Order:
// * Implant creates the pipe
// * Implant opens implant side handle to the pipe
// * Implant passes pipe name to the module via load (execute) invocation
// * Module opens module-side handle to the pipe
// * Module returns from execute
// * Module  and implant communicate via pipe
// * Module closes module side handle to pipe
// * Implant detects that module has closed module-side handle to pipe and closes implant side handle in response
// * Once all handles to the pipe are closed, implant frees the pipe
// In addition to this order, both the implant and modyule must be able to handle a pipe being close unexpectedly. In this situation the implant or module must clsoe its own open handle, and react as appropriate
func ExecuteCmd(cmd string, execChannels *component.ComponentExecutionChannel, inPipe io.Reader, outPipe io.Reader) (*bytes.Reader, error) {

}