package goddis

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestKeysCommands(t *testing.T) {
	g := NewGoddis()
	Convey("Exists", t, func() {

		Convey("Given a key it should return true if found", func() {
			g.Set("key", "value")
			So(g.Exists("key"), ShouldBeTrue)
			So(g.Exists("980usdjsdfj89sdf89sdjf"), ShouldBeFalse)
		})

	})

	Convey("Keys", t, func() {

		Convey("Given a pattern it should find keys matching", nil)

	})

	Convey("Expire", t, func() {

		Convey("Given a key and TTL it should set expiry", func() {
			g.Set("key", "value")
			So(g.Expire("key", "0"), ShouldBeTrue)
			So(g.Exists("key"), ShouldBeFalse)
		})

	})

	Convey("Del", t, func() {

		Convey("Given a range of keys it should delete them", func() {
			g.Set("key1", "value")
			g.Set("key2", "value")
			//g.HSet("key", "field", "value")
			So(g.Del("key1", "key2"), ShouldEqual, 2)
			So(g.Exists("key1"), ShouldBeFalse)
			So(g.Exists("key2"), ShouldBeFalse)
			So(g.Del(""), ShouldBeZeroValue)
		})

	})

	Convey("Rename", t, func() {

		Convey("Given a key and a new key name it should rename", nil)

	})

}
