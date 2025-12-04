package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"nexus-sds.com/netcheck/pkg/core"
)

// Precompiled regex for config lines: 2-4 char check type + whitespace + hostname
var reLine = regexp.MustCompile(`^([a-zA-Z0-9]{2,4})\s+(.+)$`)

func parseHostString(input string) (*core.Host, error) {
	input = strings.TrimSpace(input)
	matches := reLine.FindStringSubmatch(input)

	if matches == nil {
		return nil, fmt.Errorf("invalid format: must be '2-4 char checktype hostname'")
	}

	return &core.Host{
		CheckType: strings.ToUpper(matches[1]),
		HostName:  matches[2],
	}, nil
}

// Stream directly from config file to hosts to avoid keeping all lines in memory
func hostsFromConfig(path string) ([]core.Host, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	hosts := make([]core.Host, 0, 128)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		h, err := parseHostString(line)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, *h)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	return hosts, nil
}

func main() {
	// Define flags
	cfgPath := flag.String("config", "netcheck.txt", "path to config file")
	cfgPathAlt := flag.String("f", "", "path to config file (alternative to -config)")
	batchMode := flag.Bool("b", false, "batch mode - disable 'press any key' prompt")
	transcriptPath := flag.String("l", "", "path to transcript log file")
	flag.Parse()

	// Use -f flag if provided, otherwise use -config
	configFile := *cfgPath
	if *cfgPathAlt != "" {
		configFile = *cfgPathAlt
	}

	// Setup logging
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	var logWriter io.Writer = consoleWriter
	var transcriptFile *os.File

	// If transcript logging is enabled, write to both console and file
	if *transcriptPath != "" {
		var err error
		transcriptFile, err = os.OpenFile(*transcriptPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal().Err(err).Str("transcript", *transcriptPath).Msg("failed to open transcript file")
		}
		defer transcriptFile.Close()

		// Create multi-writer to output to both console and file
		logWriter = io.MultiWriter(consoleWriter, transcriptFile)
	}

	log.Logger = log.Output(logWriter)
	log.Info().Msg("starting up")

	hosts, err := hostsFromConfig(configFile)
	if err != nil {
		log.Fatal().Err(err).Str("config", configFile).Msg("failed to load config")
	}
	for _, host := range hosts {
		
		checkLabel := "Unknown"
		if label, ok := core.CheckTypeNames[host.CheckType]; ok {
			checkLabel = label
		}

		log.Info().Str("host", host.HostName).Str("checkType", host.CheckType).Str("checkLabel", checkLabel).Msg("checking host")
		checkFunc, ok := core.CheckTypes[host.CheckType]
		if !ok {
			log.Error().Str("host", host.HostName).Str("checkType", host.CheckType).Str("checkLabel", checkLabel).Msg("unknown check type")
			continue
		}

		passed, err := checkFunc(host)
		if err != nil {
			log.Error().Err(err).Str("host", host.HostName).Str("checkType", host.CheckType).Str("checkLabel", checkLabel).Msg("check error")
			continue
		}

		if !passed {
			log.Error().Str("host", host.HostName).Str("checkType", host.CheckType).Str("checkLabel", checkLabel).Msg("host failed check")
		} else {
			log.Info().Str("host", host.HostName).Str("checkType", host.CheckType).Str("checkLabel", checkLabel).Msg("host passed check")
		}
	}
	log.Info().Int("hostCount", len(hosts)).Str("config", configFile).Msg("config parsed")

	// Only prompt if not in batch mode
	if !*batchMode {
		fmt.Print("Press any key to exit...")
		var input string
		fmt.Scanln(&input)
	}
}
