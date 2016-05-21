package goddis

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHashCommands(t *testing.T) {
	g := NewGoddis()

	Convey("Given a (Key, Field, Value)", t, func() {

		Convey("It should set the hash value", func() {

			Convey("It should return true if new field", func() {
				So(g.HSet("key", "field", "value"), ShouldBeTrue)
				So(g.HSet("key", "field1", "value1"), ShouldBeTrue)
				So(g.HSet("key1", "field", "value"), ShouldBeTrue)
			})

			Convey("It should return false if the field was updated", func() {
				g.HSet("key", "field", "value")
				g.HSet("key", "field1", "value1")
				g.HSet("key1", "field", "value")
				So(g.HSet("key", "field", "new value"), ShouldBeFalse)
				So(g.HSet("key", "field1", "new value1"), ShouldBeFalse)
				So(g.HSet("key1", "field", "new value"), ShouldBeFalse)
			})

		})

	})

	Convey("Given (Key, Field, Integer)", t, func() {

		Convey("It should increment based on integer value", func() {
			g.HSet("ints", "myint", "15")

			val, _ := g.HIncrBy("ints", "myint", "5")
			So(val, ShouldEqual, 20)

			val1, _ := g.HIncrBy("ints", "myint", "18")
			So(val1, ShouldEqual, 38)

			val2, _ := g.HIncrBy("ints", "myint", "-8")
			So(val2, ShouldEqual, 30)
		})

		Convey("It should create int=0 and increment if nonexistent", func() {
			val, _ := g.HIncrBy("ints", "myint1", "5")
			So(val, ShouldEqual, 5)

			val1, _ := g.HIncrBy("ints", "myint2", "-5")
			So(val1, ShouldEqual, -5)

			val2, _ := g.HIncrBy("ints", "myint3", "546")
			So(val2, ShouldEqual, 546)
		})

		Convey("It should return an error", func() {
			val, err := g.HIncrBy("ints", "int", "value")
			So(val, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)

			val1, err1 := g.HIncrBy("ints", "int", "v1alue")
			So(val1, ShouldBeZeroValue)
			So(err1, ShouldNotBeNil)

			g.HSet("key", "field", "value")
			val2, err2 := g.HIncrBy("key", "field", "1")
			So(val2, ShouldBeZeroValue)
			So(err2, ShouldNotBeNil)
		})
	})

	Convey("Given a (Key, Field)", t, func() {
		g.HSet("key", "field", "value")

		Convey("It should get the hash value", func() {
			value, ok := g.HGet("key", "field")
			value1, ok1 := g.HGet("key", "2")
			So(value, ShouldEqual, "value")
			So(ok, ShouldBeTrue)
			So(value1, ShouldBeEmpty)
			So(ok1, ShouldBeFalse)
		})

	})

	Convey("Given a (Key, ...Field)", t, func() {
		g.HSet("key", "field1", "value1")
		g.HSet("key", "field2", "value2")
		g.HSet("key", "field3", "value3")

		Convey("It should return multiple values", func() {
			vals := g.HMGet("key", "field1", "field2", "field3")
			So(vals, ShouldResemble, []string{"value1", "value2", "value3"})

			vals1 := g.HMGet("key", "field1", "nonexist", "field3", "field4")
			So(vals1, ShouldResemble, []string{"value1", "", "value3", ""})
		})

	})

}
