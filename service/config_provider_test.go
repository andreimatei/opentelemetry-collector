// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
)

type errConfigUnmarshaler struct {
	err error
}

func (ecu *errConfigUnmarshaler) Unmarshal(*config.Map, component.Factories) (*config.Config, error) {
	return nil, ecu.err
}

func TestConfigProviderUnmarshalError(t *testing.T) {
	factories, errF := componenttest.NopFactories()
	require.NoError(t, errF)

	set := newDefaultConfigProviderSettings([]string{filepath.Join("testdata", "otelcol-nop.yaml")})
	set.Unmarshaler = &errConfigUnmarshaler{err: errors.New("unmarshal_err")}

	cfgW, err := NewConfigProvider(set)
	assert.NoError(t, err)

	_, err = cfgW.Get(context.Background(), factories)
	assert.Error(t, err)

	assert.NoError(t, cfgW.Shutdown(context.Background()))
}

func TestConfigProviderValidationError(t *testing.T) {
	factories, errF := componenttest.NopFactories()
	require.NoError(t, errF)

	set := newDefaultConfigProviderSettings([]string{filepath.Join("testdata", "otelcol-invalid.yaml")})
	set.Unmarshaler = &errConfigUnmarshaler{err: errors.New("unmarshal_err")}

	cfgW, err := NewConfigProvider(set)
	assert.NoError(t, err)

	_, err = cfgW.Get(context.Background(), factories)
	assert.Error(t, err)

	assert.NoError(t, cfgW.Shutdown(context.Background()))
}
