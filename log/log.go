package log

import (
	"encoding/json"
	"fmt"
	"github.com/isaac1102/common-log/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Logger       = GetLogger()
	podName, gid string
)

type CustomWriterForFile struct {
	file *os.File
}
type LogData struct {
	Time    time.Time `json:"time"`
	Key     string    `json:"key"`
	Message string    `json:"message"`
}

// InnerLogData는 내부 메시지를 위한 구조체입니다.
type InnerLogData struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Caller  string    `json:"caller"`
	Message string    `json:"message"`
}

func Trace(msg string) {
	Logger.Trace().Msg(msg)
}
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}
func Info(msg string) {
	Logger.Info().Msg(msg)
}
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}
func Error(msg string) {
	Logger.Error().Msg(msg)
}
func Fatal(msg string) {
	Logger.Fatal().Msg(msg)
}

func GetLogger() zerolog.Logger {
	configuration := config.Cfg

	lev := configuration.Env.Level
	level, err := zerolog.ParseLevel(strings.ToLower(lev))
	zerolog.SetGlobalLevel(level)

	if err != nil {
		level = zerolog.InfoLevel
	}

	podName = os.Getenv(configuration.Env.Pod)
	if podName == "" {
		podName = "DEFAULT"
	}

	output := defineConsoleWriter()

	var writer io.Writer
	if isFilePrint(configuration.Env.PrintType) {
		logFile, err := os.OpenFile(
			fmt.Sprintf("%s%s%s_%s_%s_%s_%s_%s.log",
				configuration.Env.FilePath,
				configuration.Env.System,
				configuration.Env.Area,
				configuration.Env.Group,
				podName,
				configuration.Env.LogType,
				level,
				time.Now().Format("20060102-15")),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664)

		if err != nil {
			log.Err(err)
		}

		fileWriter := zerolog.New(CustomWriterForFile{logFile}).
			With().
			Logger()

		writer = zerolog.MultiLevelWriter(output, fileWriter)
	} else {
		writer = output
	}

	return zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger()
}

func defineConsoleWriter() io.Writer {
	gid = config.Cfg.Env.Gid
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("[[[((%s)) %s ]]]", gid, i)
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
		FormatTimestamp: func(i interface{}) string {
			t, err := time.Parse(time.RFC3339, i.(string))
			if err != nil {
				panic(err)
			}
			formatted := t.Format("06-01-02 15:04:05")
			return formatted
		},
	}
}

func isFilePrint(printTypes []string) bool {
	for _, v := range printTypes {
		if v == "f" {
			return true
		}
	}
	return false
}

// log파일을 write
func (cw CustomWriterForFile) Write(p []byte) (n int, err error) {
	var logData LogData
	err = json.Unmarshal(p, &logData)
	if err != nil {
		return 0, err
	}

	// logData.Message를 InnerLogData 구조체로 언마샬링합니다.
	var innerLogData InnerLogData
	err = json.Unmarshal([]byte(logData.Message), &innerLogData)

	if err != nil {
		return 0, err
	}

	formattedMessage := fmt.Sprintf("%s %s %s [[[((%s)) %s ]]]\n",
		innerLogData.Time.Format("06-01-02 15:04:05"),
		fmt.Sprintf(strings.ToUpper(innerLogData.Level)),
		filepath.Base(innerLogData.Caller),
		gid,
		innerLogData.Message,
	)

	return cw.file.WriteString(formattedMessage)
}
