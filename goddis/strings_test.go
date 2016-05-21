package goddis

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringCommands(t *testing.T) {
	Convey("Given a (Key, Value) pair", t, func() {
		g := NewGoddis()

		Convey("It should set a value", func() {
			So(g.Set("key", "value"), ShouldBeTrue)
			So(g.Set("key2", "value2"), ShouldBeTrue)
			val, k := g.Get("key")
			So(val, ShouldEqual, "value")
			So(k, ShouldBeTrue)
		})

		Convey("It should increment based on integer value", func() {
			g.Set("int", "12")

			val, err := g.IncrBy("int", "3")
			So(val, ShouldEqual, 15)
			So(err, ShouldBeNil)
			val1, err1 := g.IncrBy("int", "5")
			So(val1, ShouldEqual, 20)
			So(err1, ShouldBeNil)
		})

		Convey("It should create int=0 and increment if nonexistent", func() {
			val, _ := g.IncrBy("inte", "5")
			So(val, ShouldEqual, 5)
		})

	})

	Convey("Given a modifier of NX", t, func() {
		g := NewGoddis()

		Convey("It should not overwrite an existing value", func() {
			So(g.Set("key", "value", "NX"), ShouldBeTrue)
			So(g.Set("key", "value", "NX"), ShouldBeFalse)
			So(g.Set("key1", "value", "NX"), ShouldBeTrue)
		})

	})

	Convey("If both modifiers NX and XX are used", t, func() {
		g := NewGoddis()

		Convey("It should always fail", func() {
			So(g.Set("key", "value", "NX", "XX"), ShouldBeFalse)
			So(g.Set("key", "value", "NX", "XX"), ShouldBeFalse)
			So(g.Set("key1", "value", "NX", "XX"), ShouldBeFalse)
		})

	})

	Convey("Given a modifier of XX", t, func() {
		g := NewGoddis()

		Convey("It should not set a value if not existing", func() {
			So(g.Set("key", "value", "XX"), ShouldBeFalse)
			g.Set("key", "value")
			So(g.Set("key", "value", "XX"), ShouldBeTrue)
			So(g.Set("key", "value1", "XX"), ShouldBeTrue)
			So(g.Set("key1", "value1", "XX"), ShouldBeFalse)
		})

	})

	Convey("Given a modifier of EX", t, func() {
		g := NewGoddis()
		Convey("It should require number of seconds", func() {
			So(g.Set("key", "value", "EX"), ShouldBeFalse)
			So(g.Set("key", "value", "EX", "5"), ShouldBeTrue)
		})

	})

	Convey("Given a modifier of PX", t, func() {
		g := NewGoddis()
		Convey("It should require number of milliseconds", func() {
			So(g.Set("key", "value", "PX"), ShouldBeFalse)
			So(g.Set("key", "value", "PX", "5"), ShouldBeTrue)
		})

	})

	Convey("Given a key", t, func() {
		g := NewGoddis()
		So(g.Set("key", "value"), ShouldBeTrue)

		Convey("it should get a value", func() {
			val, k := g.Get("key")
			So(val, ShouldEqual, "value")
			So(k, ShouldBeTrue)
		})

		Convey("it should increment key by one", func() {
			g.Set("someint", "5")
			val, _ := g.Incr("someint")
			So(val, ShouldEqual, 6)
			g.Incr("someint")
			val1, _ := g.Incr("someint")
			So(val1, ShouldEqual, 8)
		})

	})

	Convey("Given a (...key)", t, func() {
		g := NewGoddis()
		g.Set("key1", "value1")
		g.Set("key2", "value2")
		g.Set("key3", "value3")

		Convey("it should get multiple values", func() {
			vals := g.MGet("key1", "key2", "key3")
			So(vals, ShouldResemble, []string{"value1", "value2", "value3"})

			vals1 := g.MGet("key1", "nonexist", "key3", "nonexist")
			So(vals1, ShouldResemble, []string{"value1", "", "value3", ""})
		})

	})
}
