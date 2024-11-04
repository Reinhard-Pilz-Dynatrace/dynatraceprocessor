/**
 * @license
 * Copyright 2020 Dynatrace LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dynatraceprocessor

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"regexp"
	"strings"
)

var reHostID = regexp.MustCompile(`^HOST-[a-fA-F0-9]+$`)

type CtxKey string

const KeyEntityHost = "dt.entity.host"
const MetaDataKeyDTEntityHost = CtxKey(KeyEntityHost)
const CtxKeyMetaDataPropertiesFilePaths = CtxKey("MetaDataPropertiesFilePaths")
const CtxKeyRuxitHostIDFilePaths = CtxKey("RuxitHostIDFilePaths")

var evaluatedHostID = EvalHostID(context.Background())

// GetHostID attempts to evaluate the HostID based on
// a selected few configuration files on the current host
// If none of these files contains valid content or none of these
// files exists an empty string is getting returned
func GetHostID(ctx context.Context) string {
	if value := ctx.Value(MetaDataKeyDTEntityHost); value != nil {
		if stringValue, ok := value.(string); ok {
			return stringValue
		}
	}
	return evaluatedHostID
}

// EvalHostID attempts to evaluate the HostID based on
// a selected few configuration files on the current host
// If none of these files contains valid content or none of these
// files exists an empty string is getting returned
// If the evaluated doesn't match the format of a valid
// Host ID as expected by Dynatrace, an empty string is
// getting returned
func EvalHostID(ctx context.Context) string {
	hostID := evalHostIDValue(ctx)
	if reHostID.MatchString(hostID) {
		return hostID
	}
	return ""
}

// evalHostIDValue attempts to evaluate the HostID based on
// a selected few configuration files on the current host
// If none of these files contains valid content or none of these
// files exists an empty string is getting returned
func evalHostIDValue(ctx context.Context) string {
	var hostID string
	var err error

	var defaultMetaDataPropertiesFilePaths = []string{
		"dt_metadata_e617c525669e072eebe3d0f08212e8f2.properties",
		"/var/lib/dynatrace/enrichment/dt_metadata.properties",
	}
	metaDataPropertiesFilePaths := defaultMetaDataPropertiesFilePaths
	// productive file paths will be unavailable during unit tests
	// context contains temporary files in that case
	if value := ctx.Value(CtxKeyMetaDataPropertiesFilePaths); value != nil {
		if values, ok := value.([]string); ok {
			metaDataPropertiesFilePaths = values
		}
	}

	for _, metaDataPropertiesFilePath := range metaDataPropertiesFilePaths {
		hostID, err = evalHostIDFromProperties(metaDataPropertiesFilePath)
		if len(hostID) > 0 && err == nil {
			return hostID
		}
	}

	var defaultRuxitHostIDFilePaths = []string{
		"C:\\ProgramData\\dynatrace\\oneagent\\agent\\config\\ruxithost.id",
		"/var/lib/dynatrace/oneagent/agent/config/ruxithost.id",
	}
	ruxitHostIDFilePaths := defaultRuxitHostIDFilePaths
	// productive file paths will be unavailable during unit tests
	// context contains temporary files in that case
	if value := ctx.Value(CtxKeyRuxitHostIDFilePaths); value != nil {
		if values, ok := value.([]string); ok {
			ruxitHostIDFilePaths = values
		}
	}

	for _, ruxitHostIDFilePath := range ruxitHostIDFilePaths {
		hostID, err = evalHostIDFromRuxitHostID(ruxitHostIDFilePath)
		if len(hostID) > 0 && err == nil {
			return hostID
		}
	}
	return ""
}

// evalHostIDFromProperties evaluates the HostID based on a "magic"
// file. That file doesn't exist on the files system.
// OneAgent ensures that the current process is able to read that file.
// If OneAgent isn't running or isn't injected into the running process
// an empty string is getting returned.
// If the contents of the file identified by the parameter `filePath`
// doesn't contain the expected contents an empty string is getting returned
func evalHostIDFromProperties(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	sContent := strings.TrimSpace(string(content))
	if strings.HasSuffix(string(sContent), ".properties") {
		content, err = os.ReadFile(string(sContent))
		if err != nil {
			return "", err
		}
		sContent = strings.TrimSpace(string(content))
	}

	buf := bytes.NewBufferString(sContent)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == string(MetaDataKeyDTEntityHost) {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", nil
}

// evalHostIDFromRuxitHostID evaluates the HostID based on a configuration file
// named `ruxithost.id`.
// The first line of that file is expected to contain a hexadecimal number
// which, when prefixed with `HOST-` is the entity ID of the host monitored
// by the installed Agent
// If the file identified by the parameter `filePath` doesn't exist,
// is unaccessible or doesn't contain any lines
// an empty string is getting returned
func evalHostIDFromRuxitHostID(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		return "HOST-" + strings.TrimSpace(scanner.Text()), nil
	}

	return "", nil
}
