package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var mu sync.Mutex

// ConsoleWriterWithLevel is a custom writer that implements WriteLevel for zerolog.
type ConsoleWriterWithLevel struct {
	zerolog.ConsoleWriter
}

func (c ConsoleWriterWithLevel) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	return c.ConsoleWriter.Write(p)
}

// FileWriterWithLevel is a custom writer that implements WriteLevel for lumberjack logger.
type FileWriterWithLevel struct {
	*lumberjack.Logger
}

func (f FileWriterWithLevel) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	// Parse the JSON log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal(p, &logEntry); err != nil {
		return 0, err
	}

	// Extract the main fields
	timestamp, _ := logEntry["time"].(string)
	message, _ := logEntry["message"].(string)
	caller, _ := logEntry["caller"].(string)

	funcName := ""
	if parts := strings.Split(caller, "/"); len(parts) > 0 {
		funcName = parts[len(parts)-1]
	}

	// Collect additional fields
	var additionalFields []string
	for key, value := range logEntry {
		if key != "time" && key != "message" && key != "level" && key != "caller" {
			additionalFields = append(additionalFields, fmt.Sprintf("%s=%v", key, value))
		}
	}

	// Safely handle the timestamp string
	var formattedTimestamp string
	if len(timestamp) >= 22 {
		formattedTimestamp = strings.ReplaceAll(timestamp, "T", " ")[:22]
	} else {
		formattedTimestamp = timestamp
	}

	// Format the log entry
	formattedLog := fmt.Sprintf("%s | %s | %s | %s | %s\n",
		formattedTimestamp,
		level.String(),
		funcName,
		message,
		strings.Join(additionalFields, " "),
	)

	// Write the formatted log to the file
	_, err = f.Logger.Write([]byte(formattedLog))
	return len(p), err
}

// deleteOldLogFiles ensures the log folder retains only the latest specified number of files.
func deleteOldLogFiles(logDir string, maxFiles int) error {
	files, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	// Filter only log files
	logFiles := []os.DirEntry{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			logFiles = append(logFiles, file)
		}
	}

	// Sort files by modification time (oldest first)
	sort.Slice(logFiles, func(i, j int) bool {
		infoI, _ := logFiles[i].Info()
		infoJ, _ := logFiles[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Remove files if they exceed the limit
	if len(logFiles) > maxFiles {
		for _, file := range logFiles[:len(logFiles)-maxFiles] {
			err := os.Remove(filepath.Join(logDir, file.Name()))
			if err != nil {
				log.Err(err).Msgf("Failed to delete old log file: %s", file.Name())
			}
		}
	}

	return nil
}

// InitLogger initializes the logger with dynamic filename based on current time and day.
func InitLogger() {
	if err := Configvar.Load(); err != nil {
		log.Err(err).Msgf("Error loading config: %v", err)
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.999"
	//comp := "User"

	// Configure console writer
	consoleWriter := ConsoleWriterWithLevel{
		ConsoleWriter: zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04:05.999",
			FormatCaller: func(i interface{}) string {
				caller := i.(string)
				parts := strings.Split(caller, "/")
				if len(parts) > 0 {
					return "\033[34m" + parts[len(parts)-1] + "\033[0m"
				}
				return ""
			},
		},
	}

	var writers []io.Writer
	writers = append(writers, consoleWriter)

	// Determine if logs should be saved to file
	if Configvar.App.LogToFile {
		// Determine the log directory
		logDir := "./logs"

		// Create the log directory if it doesn't exist
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Err(err).Msgf("Failed to create log directory: %v", err)
		}

		// Perform log rotation (retain only the latest files)
		if err := deleteOldLogFiles(logDir, Configvar.App.MaxLogFiles); err != nil {
			log.Err(err).Msgf("Failed to delete old log files: %v", err)
		}

		// Dynamic log filename
		currentTime := time.Now().Format("02-01-2006_15-04-05")
		logFilename := fmt.Sprintf("eventy_%s.log", currentTime)
		logFilePath := filepath.Join(logDir, logFilename)

		// Configure file writer
		fileWriter := &FileWriterWithLevel{
			Logger: &lumberjack.Logger{
				Filename:   logFilePath,
				MaxSize:    Configvar.App.MaxFileSize,
				MaxBackups: 3,
				MaxAge:     30,
				Compress:   true,
			},
		}
		writers = append(writers, fileWriter)
	}

	// Use MultiLevelWriter for console and file output
	multiWriter := zerolog.MultiLevelWriter(writers...)

	mu.Lock()
	defer mu.Unlock()
	log.Logger = zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}

	// Set global log level
	switch strings.ToLower(Configvar.App.LogLevel) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
