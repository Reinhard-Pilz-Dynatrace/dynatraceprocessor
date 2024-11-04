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

package dynatraceprocessor_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Reinhard-Pilz-Dynatrace/dynatraceprocessor"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHostIDEvaluation(t *testing.T) {
	tests := []struct {
		name                    string
		metaDataPropertiesFiles []string
		ruxitIDFiles            []string
		expectedHostID          string
	}{
		{
			name:                    "no_files_configured",
			metaDataPropertiesFiles: []string{},
			ruxitIDFiles:            []string{},
			expectedHostID:          "",
		},
		{
			name:                    "valid_ruxit_id_file",
			metaDataPropertiesFiles: []string{},
			ruxitIDFiles: []string{
				`AAF98EFF909EE3F6`,
			},
			expectedHostID: "HOST-AAF98EFF909EE3F6",
		},
		{
			name:                    "valid_multi_line_ruxit_id_file",
			metaDataPropertiesFiles: []string{},
			ruxitIDFiles: []string{
				`AAF98EFF909EE3F6
				asdfsdf`,
			},
			expectedHostID: "HOST-AAF98EFF909EE3F6",
		},
		{
			name:                    "invalid_ruxit_id_file",
			metaDataPropertiesFiles: []string{},
			ruxitIDFiles: []string{
				`ZZF98EFF909EE3F6
				asdfsdf`,
			},
			expectedHostID: "",
		},
		{
			name:                    "invalid_metadata_file",
			metaDataPropertiesFiles: []string{},
			ruxitIDFiles: []string{
				`ZZF98EFF909EE3F6
				asdfsdf`,
			},
			expectedHostID: "",
		},
		{
			name: "valid_metadata_file",
			metaDataPropertiesFiles: []string{
				`dt.entity.host=HOST-AAF98EFF909EE3F6
				asdfsdf`,
			},
			ruxitIDFiles:   []string{},
			expectedHostID: "HOST-AAF98EFF909EE3F6",
		},
		{
			name:                    "non_existent_files",
			metaDataPropertiesFiles: []string{"nil"},
			ruxitIDFiles:            []string{"nil"},
			expectedHostID:          "",
		},
	}

	createConfigFile := func(content string) (*os.File, error) {
		tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.txt", uuid.NewString()))
		if err != nil {
			return nil, err
		}
		if _, err := tempFile.Write([]byte(content)); err != nil {
			tempFile.Close()
			return nil, err
		}
		return tempFile, nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metaDataPropertiesFilePaths := []string{}
			for _, content := range tt.metaDataPropertiesFiles {
				if content != "nil" {
					tempFile, err := createConfigFile(content)
					if err != nil {
						return
					}
					defer os.Remove(tempFile.Name())
					metaDataPropertiesFilePaths = append(metaDataPropertiesFilePaths, tempFile.Name())
				} else {
					metaDataPropertiesFilePaths = append(metaDataPropertiesFilePaths, uuid.NewString())
				}
			}
			ctx := context.WithValue(context.Background(), dynatraceprocessor.CtxKeyMetaDataPropertiesFilePaths, metaDataPropertiesFilePaths)
			ruxitHostIDFilePaths := []string{}
			for _, content := range tt.ruxitIDFiles {
				if content != "nil" {
					tempFile, err := createConfigFile(content)
					if err != nil {
						return
					}
					defer os.Remove(tempFile.Name())
					ruxitHostIDFilePaths = append(ruxitHostIDFilePaths, tempFile.Name())
				} else {
					ruxitHostIDFilePaths = append(ruxitHostIDFilePaths, uuid.NewString())
				}
			}
			ctx = context.WithValue(ctx, dynatraceprocessor.CtxKeyRuxitHostIDFilePaths, ruxitHostIDFilePaths)

			hostID := dynatraceprocessor.EvalHostID(ctx)

			assert.Equal(t, tt.expectedHostID, hostID)
		})
	}
}
