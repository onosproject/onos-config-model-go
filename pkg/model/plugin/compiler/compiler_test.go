// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugincompiler

import (
	"context"
	"github.com/onosproject/onos-config-model/pkg/model"
	plugincache "github.com/onosproject/onos-config-model/pkg/model/plugin/cache"
	pluginmodule "github.com/onosproject/onos-config-model/pkg/model/plugin/module"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCompiler(t *testing.T) {
	t.Skip()

	resolver := pluginmodule.NewResolver(pluginmodule.ResolverConfig{
		Path:   filepath.Join(moduleRoot, "test", "mod"),
		Target: "github.com/onosproject/onos-config@master",
	})

	cache, err := plugincache.NewPluginCache(plugincache.CacheConfig{
		Path: filepath.Join(moduleRoot, "test", "cache"),
	}, resolver)
	assert.NoError(t, err)

	config := CompilerConfig{
		TemplatePath: filepath.Join(moduleRoot, "pkg", "model", "plugin", "compiler", "templates"),
		BuildPath:    filepath.Join(moduleRoot, "test", "build"),
	}

	bytes, err := ioutil.ReadFile(filepath.Join(moduleRoot, "test", "test@2020-11-18.yang"))
	assert.NoError(t, err)

	modelInfo := configmodel.ModelInfo{
		Name:         "test",
		Version:      "1.0.0",
		GetStateMode: configmodel.GetStateExplicitRoPaths,
		Modules: []configmodel.ModuleInfo{
			{
				Name:         "test",
				Organization: "ONF",
				Revision:     "2020-11-18",
				File:         "test.yang",
			},
		},
		Files: []configmodel.FileInfo{
			{
				Path: "test@2020-11-18.yang",
				Data: bytes,
			},
		},
		Plugin: configmodel.PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		},
	}

	entry := cache.Entry("test", "1.0.0")
	err = entry.Lock(context.TODO())
	assert.NoError(t, err)

	compiler := NewPluginCompiler(config, resolver)
	err = compiler.CompilePlugin(modelInfo, entry.Path)
	assert.NoError(t, err)

	plugin, err := entry.Load()
	assert.NoError(t, err)
	assert.NotNil(t, plugin)

	err = entry.Unlock(context.TODO())
	assert.NoError(t, err)
}
