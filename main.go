package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Host struct {
	HostName  string
	CheckType string
}

// Precompiled regex for config lines: 3-char check type + whitespace + hostname
var reLine = regexp.MustCompile(`^([a-zA-Z0-9]{3})\s+(.+)$`)

func parseHostString(input string) (*Host, error) {
	input = strings.TrimSpace(input)
	matches := reLine.FindStringSubmatch(input)

	if matches == nil {
		return nil, fmt.Errorf("invalid format: must be '3-char-checktype hostname'")
	}

	return &Host{
		CheckType: strings.ToUpper(matches[1]),
		HostName:  matches[2],
	}, nil
}

// Stream directly from config file to hosts to avoid keeping all lines in memory
func hostsFromConfig(path string) ([]Host, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	hosts := make([]Host, 0, 128)
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("starting up")

	cfgPath := flag.String("config", "netcheck.txt", "path to config file")
	flag.Parse()

	hosts, err := hostsFromConfig(*cfgPath)
	if err != nil {
		log.Fatal().Err(err).Str("config", *cfgPath).Msg("failed to load config")
	}
	for _, host := range hosts {
		log.Info().Str("host", host.HostName).Str("checkType", host.CheckType).Msg("checking host")
	}
	log.Info().Int("hostCount", len(hosts)).Str("config", *cfgPath).Msg("config parsed")

	fmt.Print("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
}
