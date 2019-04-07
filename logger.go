package remember

const logPatternRed = "\x1b[31m%s\x1b[39;49m\n"
const logPatternBlue = "\x1b[36m%s\x1b[39;49m\n"

// Logger provides an interface to log extra debug information.
// The glog package can be used or alternatively you can defined your own.
//
// Example:
//
//  import log
//
//  type aLogger struct {}
//
//  func (l aLogger) Log(format string, args ...interface{}) {
//     log.Printf(format, args...)
//  }
type Logger interface {
	// Log follows the same pattern as fmt.Printf( ).
	Log(format string, args ...interface{})
}
