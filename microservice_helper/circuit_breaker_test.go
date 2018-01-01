package microservice_helper

/*

 */
import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	. "github.com/smartystreets/goconvey/convey"
)

const RESULT = "Done"

func TestRunSuccessfully(t *testing.T) {

	Convey("Given a service call", t, func() {
		serviceInvoke := func() (interface{}, error) {
			//time.Sleep(time.Second * 1)
			return RESULT, nil
		}
		Convey("When invoking obeys the settings", func() {
			hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
				Timeout: 2000,
			})

			ret, err := CallDependentService("my_command", serviceInvoke, nil)
			Convey("Then a result should be returned", func() {
				So(err, ShouldBeNil)
				So(ret, ShouldEqual, RESULT)
			})
		})

	})
}

func TestCutOffWhenTimeout(t *testing.T) {

	Convey("Given a long run service call", t, func() {
		serviceInvoke := func() (interface{}, error) {
			time.Sleep(time.Second * 3)
			return RESULT, nil
		}
		Convey("When timeout happens", func() {
			hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
				Timeout: 1000,
			})

			_, err := CallDependentService("my_command", serviceInvoke, nil)
			Convey("Then a timeout error occurs", func() {
				So(IsTimeoutError(err), ShouldBeTrue)
			})
		})

	})
}

func TestErrorWillBeThrownOut(t *testing.T) {

	Convey("Given a service call", t, func() {
		serviceInvoke := func() (interface{}, error) {
			return nil, errors.New("Error occured")
		}

		Convey("When error happens and no fallback set ", func() {
			hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
				Timeout: 1000,
			})

			_, err := CallDependentService("my_command", serviceInvoke, nil)
			Convey("Then the error would be thrown out.", func() {
				fmt.Println(err)
				So(err, ShouldNotBeNil)

			})
		})

	})
}

func TestErrorOccurFallbackMethodWillBeInvoked(t *testing.T) {
	fallBackRet := "Fallback"
	Convey("Given a service call and fallback method", t, func() {
		serviceInvoke := func() (interface{}, error) {
			return nil, errors.New("Error occured")
		}
		fallBack := func(err error) (interface{}, error) {
			return fallBackRet, nil
		}
		Convey("When error happens, ", func() {
			hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
				Timeout: 1000,
			})

			ret, err := CallDependentService("my_command", serviceInvoke, fallBack)
			Convey("Then fallback method executed.", func() {

				So(err, ShouldBeNil)
				So(ret, ShouldEqual, fallBackRet)
			})
		})

	})
}

func TestAutoRetryWhenErrorOccur(t *testing.T) {
	Convey("Given a retryable function", t, func() {
		executeCnt := 0
		retrySettings := RetrySettings{
			retryTimes:             3,
			retryInterval:          time.Second * 1,
			retryIntervalIncrement: time.Millisecond * 500,
		}

		Convey("When retryable error occurred", func() {
			t1 := time.Now()
			AutoRetry(func() (interface{}, error) {
				executeCnt++
				return nil, errors.New("Retryable")
			}, retrySettings, []string{"Retryable"})
			escaped_time := time.Since(t1)
			Convey("Then logic has been retried", func() {
				So(executeCnt, ShouldEqual, 4)
				So(escaped_time.Seconds(), ShouldBeGreaterThanOrEqualTo,
					(1 + (1 + 0.5) + (1 + 0.5 + 0.5)))
			})
		})
	})
}

func TestAutoRetrySucceedAfterRetry(t *testing.T) {
	Convey("Given a retryable function", t, func() {
		retrySettings := RetrySettings{
			retryTimes:             3,
			retryInterval:          time.Second * 1,
			retryIntervalIncrement: time.Millisecond * 500,
		}
		executeCnt := 0
		Convey("When retryable error occurred", func() {

			ret, err := AutoRetry(func() (interface{}, error) {
				executeCnt++
				if executeCnt > 1 {
					return RESULT, nil
				}
				return nil, errors.New("Retryable")
			}, retrySettings, []string{"Retryable"})

			Convey("Then after retrying, the logic would be executed successfully",
				func() {
					So(executeCnt, ShouldEqual, 2)
					So(err, ShouldBeNil)
					So(ret, ShouldEqual, RESULT)
				})
		})
	})
}
