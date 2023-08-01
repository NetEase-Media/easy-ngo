package xlog

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/valyala/bytebufferpool"

	"g.hz.netease.com/ngo/ngo/adapter/util"

	"github.com/gin-gonic/gin"
)

// OptFunc is the type to use to options to the option struct during initialization
type OptFunc func(*opt)

// opt is the internal struct that holds the options for logging.
type opt struct {
	Output io.Writer
	Time   time.Time
}

// newOpt returns a new struct to hold options, with the default output to stdout.
func newOpt() *opt {
	o := new(opt)
	o.Output = os.Stdout
	return o
}

// WithOutput sets the io.Writer output for the log file.
func WithOutput(out io.Writer) OptFunc {
	return func(o *opt) {
		o.Output = out
	}
}

// responseWriter is the internal struct that will wrap the http.ResponseWriter
// and hold the status and number of bytes written
type responseWriter struct {
	gin.ResponseWriter

	start time.Time
}

// startTime sets the start time to calculate the elapsed time for the %D directive
func (rw *responseWriter) startTime() {
	rw.start = time.Now()
}

const (
	// ApacheCommonLogFormat is the Apache Common Log directives
	ApacheCommonLogFormat = `%h %l %u %t "%r" %>s %b`

	// ApacheCombinedLogFormat is the Apache Combined Log directives
	ApacheCombinedLogFormat = `%h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i"`
)

var timeFmtMap = map[rune]string{
	'a': "Mon", 'A': "Monday", 'b': "Jan", 'B': "January", 'C': "06",
	'd': "02", 'D': "01/02/06", 'e': "_2", 'F': "2006-01-02",
	'f': "15:04:05.000",
	'h': "Jan", 'H': "15", 'I': "3", 'k': "•15", 'l': "_3",
	'm': "01", 'M': "04", 'n': "\n", 'p': "PM", 'P': "pm",
	'r': "03:04:05 PM", 'R': "15:04", 'S': "05",
	't': "\t", 'T': "15:04:05", 'y': "06", 'Y': "2006",
	'z': "-700", 'Z': "MST", '%': "%%",

	// require calculated time
	'G': "%v", 'g': "%v", 'j': "%v", 's': "%v",
	'u': "%v", 'V': "%v", 'w': "%v",

	// Unsupported directives
	'c': "?", 'E': "?", 'O': "?", 'U': "?",
	'W': "?", 'x': "?", 'X': "?", '+': "?",
}

// convertTimeFormat converts strftime formatting directives to a go time.Time format
func convertTimeFormat(now time.Time, format string) string {
	var isDirective bool
	var buf = new(bytes.Buffer)
	for _, r := range format {
		if !isDirective && r == '%' {
			isDirective = true
			continue
		}
		if !isDirective {
			buf.WriteRune(r)
			continue
		}
		if val, ok := timeFmtMap[r]; ok {
			switch val {
			case "%v":
				switch r {
				case 'G':
					y, _ := now.ISOWeek()
					buf.WriteString(strconv.Itoa(y))
				case 'g':
					y, _ := now.ISOWeek()
					y -= (y / 100) * 100
					buf.WriteString(fmt.Sprintf("%02d", y))
				case 'j':
					buf.WriteString(strconv.Itoa(now.YearDay()))
				case 's':
					buf.WriteString(strconv.FormatInt(now.Unix(), 10))
				case 'u':
					w := now.Weekday()
					if w == 0 {
						w = 7
					}
					buf.WriteString(strconv.Itoa(int(w)))
				case 'V':
					_, w := now.ISOWeek()
					buf.WriteString(strconv.Itoa(w))
				case 'w':
					buf.WriteString(strconv.Itoa(int(now.Weekday())))
				}
			default:
				buf.WriteString(now.Format(val))
			}
			isDirective = false
			continue
		}
		buf.WriteString("(%" + string(r) + " is invalid)")
	}
	return buf.String()
}

// line is the type that will hold all of the runtime formating directives for the log line
type line struct {
	time    time.Time
	request *http.Request
	writer  *responseWriter
	keys    map[string]interface{}

	// directives
	a, A, b, B, h, H, l, m, p, q, r, s, S, t, u, U, v, D, T, I string
}

func (ln *line) withTime(o *opt) *line {
	if !o.Time.IsZero() {
		ln.time = o.Time
		return ln
	}
	ln.time = time.Now()
	return ln
}

func (ln *line) withRequest(r *http.Request) *line {
	ln.request = r
	return ln
}

func (ln *line) withResponse(a *responseWriter) *line {
	ln.writer = a
	return ln
}

func (ln *line) withKeys(k map[string]interface{}) *line {
	ln.keys = k
	return ln
}

// remoteAddr - %a
func (ln *line) remoteAddr() string {
	if len(ln.a) == 0 {
		remoteIp, _, _ := net.SplitHostPort(strings.TrimSpace(ln.request.RemoteAddr))
		ln.a = remoteIp
		if len(ln.a) == 0 {
			ln.a = "-"
		}
	}
	return ln.a
}

// localAddr - %A
func (ln *line) localAddr() string {
	if len(ln.A) == 0 {
		ip, _ := util.GetOutBoundIP()
		ln.A = ip
		if len(ln.A) == 0 {
			ln.A = "-"
		}
	}
	return ln.A
}

// bytesWritten - %b
func (ln *line) bytesWritten() string {
	if len(ln.b) == 0 {
		size := ln.writer.Size()
		if size == 0 {
			ln.b = "-"
		} else {
			ln.b = strconv.Itoa(size)
		}
	}
	return ln.b
}

// bytesWritten - %B
func (ln *line) bytesWritten2() string {
	if len(ln.B) == 0 {
		ln.B = strconv.Itoa(ln.writer.Size())
	}
	return ln.B
}

// remoteHostname - %h
func (ln *line) remoteHostname() string {
	if len(ln.h) == 0 {
		ln.h = ln.request.URL.Host
		if len(ln.h) == 0 {
			ln.h = "-"
		}
	}
	return ln.h
}

// protocol - %H
func (ln *line) protocol() string {
	if len(ln.H) == 0 {
		ln.H = ln.request.Proto
	}
	return ln.H
}

// remoteUsername - %l
func (ln *line) remoteUsername() string {
	if len(ln.l) == 0 {
		ln.l = "-"
	}
	return ln.l
}

// method - %m
func (ln *line) method() string {
	if len(ln.m) == 0 {
		ln.m = ln.request.Method
	}
	return ln.m
}

// port - %m
func (ln *line) port() string {
	if len(ln.p) == 0 {
		ln.p = "-"
	}
	return ln.p
}

// queryString - %q
func (ln *line) queryString() string {
	if len(ln.q) == 0 {
		query := ln.request.URL.RawQuery
		if len(query) > 0 {
			ln.q = "?" + query
		} else {
			ln.q = "-"
		}
	}
	return ln.q
}

// requestLine - %r
func (ln *line) requestLine() string {
	if len(ln.r) == 0 {
		ln.r = strings.ToUpper(ln.request.Method) + " " + ln.request.RequestURI + " " + ln.request.Proto
	}
	return ln.r
}

// status - %s
func (ln *line) status() string {
	if len(ln.s) == 0 {
		ln.s = strconv.Itoa(ln.writer.Status())
	}
	return ln.s
}

// sessionId - %S
func (ln *line) sessionId() string {
	if len(ln.S) == 0 {
		ln.S = "-"
	}
	return ln.S
}

// timeFormatted - %t
func (ln *line) timeFormatted(format string) string {
	if len(ln.t) == 0 {
		ln.t = ln.time.Format(format)
	}
	return ln.t
}

// username - %u
func (ln *line) username() string {
	if len(ln.u) == 0 {
		ln.u = "-"
		if s := strings.SplitN(ln.request.Header.Get("Authorization"), " ", 2); len(s) == 2 {
			if b, err := base64.StdEncoding.DecodeString(s[1]); err == nil {
				if pair := strings.SplitN(string(b), ":", 2); len(pair) == 2 {
					ln.u = pair[0]
				}
			}
		}
	}
	return ln.u
}

// url - %U
func (ln *line) url() string {
	if len(ln.U) == 0 {
		ln.U = ln.request.URL.Path
	}
	return ln.U
}

// serviceName - %v
func (ln *line) serviceName() string {
	if len(ln.v) == 0 {
		ln.v = "-"
	}
	return ln.v
}

// timeElapsedMs - %D
func (ln *line) timeElapsedMs(cost time.Duration) string {
	if len(ln.D) == 0 {
		ln.D = strconv.FormatInt(cost.Milliseconds(), 10)
	}
	return ln.D
}

// timeElapsedSeconds - %T
func (ln *line) timeElapsedSeconds(cost time.Duration) string {
	if len(ln.T) == 0 {
		ln.T = strconv.FormatFloat(cost.Seconds(), 'f', 3, 64)
	}
	return ln.T
}

// threadName - %I
func (ln *line) threadName() string {
	if len(ln.I) == 0 {
		ln.I = "-"
	}
	return ln.I
}

// flatten takes two slices and merges them into one
func flatten(o *opt, a, b []string) func(w *responseWriter, r *http.Request, k map[string]interface{}) string {
	return func(w *responseWriter, r *http.Request, k map[string]interface{}) string {
		ln := new(line)
		ln.withTime(o).withRequest(r).withResponse(w).withKeys(k)
		cost := time.Since(ln.writer.start)

		buf := bytebufferpool.Get()
		for i, s := range a {
			switch s {
			case "":
				buf.WriteString(b[i])
			case "%a":
				buf.WriteString(ln.remoteAddr())
			case "%A":
				buf.WriteString(ln.localAddr())
			case "%b":
				buf.WriteString(ln.bytesWritten())
			case "%B":
				buf.WriteString(ln.bytesWritten2())
			case "%h":
				buf.WriteString(ln.remoteHostname())
			case "%H":
				buf.WriteString(ln.protocol())
			case "%l":
				buf.WriteString(ln.remoteUsername())
			case "%m":
				buf.WriteString(ln.method())
			case "%p":
				buf.WriteString(ln.port())
			case "%q":
				buf.WriteString(ln.queryString())
			case "%r":
				buf.WriteString(ln.requestLine())
			case "%s", "%>s":
				buf.WriteString(ln.status())
			case "%S":
				buf.WriteString(ln.sessionId())
			case "%t":
				buf.WriteString(ln.timeFormatted("[02/01/2006:15:04:05 -0700]"))
			case "%u":
				buf.WriteString(ln.username())
			case "%U":
				buf.WriteString(ln.url())
			case "%v":
				buf.WriteString(ln.serviceName())
			case "%D":
				buf.WriteString(ln.timeElapsedMs(cost))
			case "%T":
				buf.WriteString(ln.timeElapsedSeconds(cost))
			case "%I":
				buf.WriteString(ln.threadName())
			default:
				if len(s) > 4 && s[:2] == "%{" && s[len(s)-2] == '}' {
					label := s[2 : len(s)-2]
					switch s[len(s)-1] {
					case 'i':
						reqHeader := r.Header.Get(label)
						if len(reqHeader) == 0 {
							reqHeader = "-"
						}
						buf.WriteString(reqHeader)
					case 'o':
						rspHeader := w.Header().Get(label)
						if len(rspHeader) == 0 {
							rspHeader = "-"
						}
						buf.WriteString(rspHeader)
					case 'c':
						if cookie, err := r.Cookie(label); err != nil && cookie != nil && len(cookie.Value) > 0 {
							buf.WriteString(cookie.Value)
						} else {
							buf.WriteString("-")
						}
					case 't':
						buf.WriteString(convertTimeFormat(ln.time, label))
					case 'r':
						if v, ok := ln.keys[label]; ok {
							buf.WriteString(v.(string))
						} else {
							buf.WriteString("-")
						}
					}
				}
			}
		}
		s := buf.String()
		bytebufferpool.Get()
		return s
	}
}

// FormatWith accepts a format string using Apache formatting directives with
// option functions and returns a function that can handle standard Gin middleware.
func FormatWith(format string, opts ...OptFunc) func(c *gin.Context) {
	if strings.EqualFold("common", format) {
		format = ApacheCommonLogFormat
	} else if strings.EqualFold("combined", format) {
		format = ApacheCombinedLogFormat
	}

	options := newOpt()
	for _, opt := range opts {
		opt(options)
	}

	var directives, betweens = make([]string, 0, 50), make([]string, 0, 50)
	var cBuf *bytes.Buffer // current buffer
	aBuf, bBuf := new(bytes.Buffer), new(bytes.Buffer)
	cBuf = bBuf

	var isDirective, isEnclosure bool
	for i, r := range format {
		switch r {
		case '%':
			if isDirective {
				cBuf.WriteRune(r)
				continue
			}
			isDirective = true
			if i != 0 {
				directives = append(directives, aBuf.String())
				betweens = append(betweens, bBuf.String())
				aBuf.Reset()
				bBuf.Reset()
			}
			cBuf = aBuf
		case '{':
			isEnclosure = true
		case '}':
			isEnclosure = false
		case '>':
			// nothing - no change in status
		default:
			if isDirective && !isEnclosure && !unicode.IsLetter(r) {
				isDirective = false
				isEnclosure = false
				if i != 0 {
					directives = append(directives, aBuf.String())
					betweens = append(betweens, bBuf.String())
					aBuf.Reset()
					bBuf.Reset()
				}
				cBuf = bBuf
			}
		}
		cBuf.WriteRune(r)
	}

	directives = append(directives, aBuf.String())
	betweens = append(betweens, bBuf.String())
	aBuf.Reset()
	bBuf.Reset()

	logFunc := flatten(options, directives, betweens)

	return func(c *gin.Context) {
		rw := &responseWriter{ResponseWriter: c.Writer}
		rw.startTime()
		c.Next()
		_, err := fmt.Fprintln(options.Output, logFunc(rw, c.Request, c.Keys))
		if err != nil {
			panic(err)
		}
	}
}
