package orlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"
)

const (
	backupFormat      = "2006-01-02"
	backupFormatRegex = "\\d{4}-\\d{2}-\\d{2}"
)

var _ io.WriteCloser = (*Logger)(nil)

type Logger struct {
	filename     string
	maxBackups   int
	lastedSuffix string
	tokenBucket  *TokenBucket
	file         *os.File
	mu           sync.Mutex
	millCh       chan bool
	startMill    sync.Once
}

func NewLogger(filename string, maxBackups int, tokenBucket *TokenBucket) *Logger {
	logger := &Logger{
		filename:    filename,
		maxBackups:  maxBackups,
		tokenBucket: tokenBucket,
	}
	return logger
}

var (
	currentTime = time.Now
	os_Stat     = os.Stat
	megabyte    = 1024 * 1024
)

func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		if err = l.openExistingOrNew(); err != nil {
			return 0, err
		}
	}
	err = l.rollover()
	if err != nil {
		return 0, err
	}
	if l.tokenBucket != nil && l.tokenBucket.Consume() {
		return 0, nil
	}
	n, err = l.file.Write(p)
	return n, err
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *Logger) rollover() error {
	cur := currentTime()
	curSuffix := cur.Format(backupFormat)
	if curSuffix == l.lastedSuffix {
		return nil
	} else {
		err := l.rotate()
		if err == nil {
			l.lastedSuffix = curSuffix
		}
		return err
	}
}

func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.openNew(); err != nil {
		return err
	}
	l.mill()
	return nil
}

func (l *Logger) openNew() error {
	err := os.MkdirAll(l.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.filename
	mode := os.FileMode(0644)
	info, err := os_Stat(name)
	if err == nil {
		mode = info.Mode()
		newname := backupName(name, l.lastedSuffix)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}
		if err := chown(name, info); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.lastedSuffix = currentTime().Format(backupFormat)
	l.file = f
	return nil
}

func backupName(name, timeSuffix string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	return filepath.Join(dir, fmt.Sprintf("%s.%s", filename, timeSuffix))
}

func (l *Logger) openExistingOrNew() error {
	l.mill()
	filename := l.filename
	_, err := os_Stat(filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return l.openNew()
	}
	l.file = file
	return nil
}

func (l *Logger) regexPattern() string {
	return l.filename + "(?P<Date>" + backupFormatRegex + ")"
}

func (l *Logger) millRunOnce() error {
	files, err := l.oldLogFiles()
	if err != nil {
		return err
	}
	var remove []logInfo
	diff := time.Duration(int64(24*time.Hour) * int64(l.maxBackups))
	baseTime, _ := time.Parse(backupFormat, l.lastedSuffix)
	cutoff := baseTime.Add(-1 * diff)

	for _, f := range files {
		if f.timestamp.Before(cutoff) {
			remove = append(remove, f)
		}
	}

	for _, f := range remove {
		errRemove := os.Remove(filepath.Join(l.dir(), f.Name()))
		if err == nil && errRemove != nil {
			err = errRemove
		}
	}
	return err
}

func (l *Logger) millRun() {
	for _ = range l.millCh {
		_ = l.millRunOnce()
	}
}

func (l *Logger) mill() {
	l.startMill.Do(func() {
		l.millCh = make(chan bool, 1)
		go l.millRun()
	})
	select {
	case l.millCh <- true:
	default:
	}
}

func (l *Logger) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(l.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	logFiles := []logInfo{}
	r := regexp.MustCompile(l.regexPattern())
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fssm := r.FindStringSubmatch(f.Name())
		if fssm != nil {
			t, err := time.Parse(backupFormat, fssm[1])
			if err == nil {
				logFiles = append(logFiles, logInfo{t, f})
				continue
			}

		}
	}
	sort.Sort(byFormatTime(logFiles))
	return logFiles, nil
}

func (l *Logger) dir() string {
	return filepath.Dir(l.filename)
}

type logInfo struct {
	timestamp time.Time
	os.FileInfo
}

type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
