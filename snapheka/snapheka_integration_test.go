//
// +build integration

package snapheka

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHekaPublish(t *testing.T) {
	config := make(map[string]ctypes.ConfigValue)

	Convey("snap Plugin integration testing with Heka", t, func() {
		var buf bytes.Buffer
		buf.Reset()
		enc := gob.NewEncoder(&buf)

		config["host"] = ctypes.ConfigValueStr{Value: os.Getenv("SNAP_HEKA_HOST")}
		config["port"] = ctypes.ConfigValueInt{Value: 3242}

		op := NewHekaPublisher()
		cp, _ := op.GetConfigPolicy()
		cfg, _ := cp.Get([]string{""}).Process(config)

		Convey("Publish float metrics to Heka", func() {
			metrics := []plugin.PluginMetricType{
				*plugin.NewPluginMetricType([]string{"/intel/psutil/load/load15"}, time.Now(), "mac1", nil, nil, 23.1),
				*plugin.NewPluginMetricType([]string{"/intel/psutil/vm/available"}, time.Now().Add(2*time.Second), "mac2", nil, nil, 23.2),
				*plugin.NewPluginMetricType([]string{"/intel/psutil/load/load1"}, time.Now().Add(3*time.Second), "linux3", nil, nil, 23.3),
			}
			enc.Encode(metrics)

			err := op.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
			So(err, ShouldBeNil)
		})

		Convey("Publish int metrics to Heka", func() {
			metrics := []plugin.PluginMetricType{
				*plugin.NewPluginMetricType([]string{"/intel/psutil/vm/free"}, time.Now().Add(5*time.Second), "linux7", nil, nil, 88),
			}
			enc.Encode(metrics)

			err := op.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
			So(err, ShouldBeNil)
		})
	})
}
