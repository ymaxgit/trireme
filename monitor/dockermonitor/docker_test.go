package dockermonitor

import (
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStart(t *testing.T) {
	Convey("Given i create a docker monitor isntance", t, func() {
		//Since we are only testing Start most of these the other pointer can be nil
		//We will return very early in case of failure

		Convey("Given the docker daemon is not running", func() {
			syscall.Mkfifo("/tmp/nonexistent.sock", 0755)
			monitor := NewDockerMonitor("unix", "/tmp/nonexistent.sock", nil, defaultDockerMetadataExtractor, nil, false, nil, false)
			err := monitor.Start()
			So(err, ShouldNotBeNil)

		})
		Convey("Given the docker daemon is  running", func() {
			//This assume docker is running on /var/run/docker.sock
			monitor := NewDockerMonitor("unix", "/var/run/docker.sock", nil, defaultDockerMetadataExtractor, nil, false, nil, false)
			err := monitor.Start()
			So(err, ShouldBeNil)
			err = monitor.Stop()
			So(err, ShouldBeNil)

		})
	})
}
