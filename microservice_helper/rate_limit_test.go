package microservice_helper

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const WAIT_ANYWAY time.Duration = 0

func TestTokenBucket(t *testing.T) {
	Convey("Given a service call", t, func() {
		bucket := CreateTokenBucket(3, 2, time.Second*1)
		Convey("When controling the coming request with the token bucket\n", func() {
			cntSuccessfulInvoking := 0
			t1 := time.Now()
			var err error
			var token time.Time
			for i := 0; i < 6; i++ {
				token, err = GetToken(bucket, WAIT_ANYWAY)
				//put the service logic here
				if err == nil {
					fmt.Println(token)
					cntSuccessfulInvoking++
				}
			}
			Convey("Then the request rate is controled", func() {
				time_escaped := time.Since(t1)
				So(err, ShouldBeNil)
				So(time_escaped.Seconds(), ShouldBeGreaterThanOrEqualTo, 1+1)
				So(cntSuccessfulInvoking, ShouldEqual, 6)
			})
		})

	})
}

func TestTokenBucketWithTimeoutSetting(t *testing.T) {
	Convey("Given a service call", t, func() {
		bucket := CreateTokenBucket(3, 2, time.Second*1)
		Convey("When controling the coming request with the token bucket\n", func() {
			cntSuccessfulInvoking := 0
			t1 := time.Now()
			var err error
			var token time.Time
			for i := 0; i < 6; i++ {
				token, err = GetToken(bucket, time.Millisecond*500)
				//put the service logic here
				if err == nil {
					fmt.Println(token)
					cntSuccessfulInvoking++
				} else {
					fmt.Println(err)
				}
			}
			Convey("Then the request rate is controled", func() {
				time_escaped := time.Since(t1)
				So(err, ShouldBeNil)
				So(cntSuccessfulInvoking, ShouldEqual, 3)
				So(time_escaped.Seconds(), ShouldBeGreaterThanOrEqualTo, 1)
			})
		})

	})
}
