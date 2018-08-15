package microservice_helper

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMaxNumGroutingControl(t *testing.T) {
	Convey("Given the max number of groutings setting", t, func() {
		maxNumOfGrouting := 10
		controller := CreateGoroutineController(maxNumOfGrouting)
		Convey("When the number of the goroutings beyond the limit", func() {
			var err error
			for i := 0; i < maxNumOfGrouting+1; i++ {
				err = controller.ApplyToken(time.Second * 2)
				if err == nil {
					go func(d int) {
						defer controller.ReleaseToken()
						fmt.Println(d)
						time.Sleep(time.Second * 5)
					}(i)
				}
			}
			Convey("then can't create new gorouting", func() {
				So(controller.NumOfGoroutingCanBeCreated(), ShouldEqual, 0)
				So(err, ShouldEqual, TokenWaitingTimeoutError)
			})
		})
	})
}

func TestMaxNumGroutingControlByCall_StartGrouting(t *testing.T) {
	Convey("Given the max number of groutings setting", t, func() {
		maxNumOfGrouting := 10
		controller := CreateGoroutineController(maxNumOfGrouting)
		Convey("When the number of the goroutings beyond the limit", func() {
			var err error
			for i := 0; i < maxNumOfGrouting+1; i++ {
				err = controller.StartGorouting(func(ii interface{}) {
					fmt.Println(ii.(int))
					time.Sleep(time.Second * 5)
				}, time.Second*2, i)
			}
			Convey("then can't create new gorouting", func() {
				So(controller.NumOfGoroutingCanBeCreated(), ShouldEqual, 0)
				So(err, ShouldEqual, TokenWaitingTimeoutError)
			})
		})
	})

}

func TestExampleGroutingNumberControl(t *testing.T) {
	Convey("Given an example", t, func() {
		Convey("When run the example", func() {
			err := ExampleGroutingNumberControl()
			Convey("Then the result is expected", func() {
				So(err == TokenWaitingTimeoutError, ShouldBeTrue)
			})
		})
	})
}
