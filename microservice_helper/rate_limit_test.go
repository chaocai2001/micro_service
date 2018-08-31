package microservice_helper

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

//const WAIT_ANYWAY time.Duration = 0

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
				So(time_escaped.Seconds(), ShouldBeGreaterThanOrEqualTo, 0.5*3)
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
				token, err = GetToken(bucket, time.Millisecond*520)
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
				So(cntSuccessfulInvoking, ShouldEqual, 6)
				So(time_escaped.Seconds(), ShouldBeGreaterThanOrEqualTo, 1)
			})
		})

	})
}

func TestTryToGetTokenBucket(t *testing.T) {
	Convey("Given a service call", t, func() {
		bucket := CreateTokenBucket(3, 2, time.Second*1)
		Convey("When controling the coming request with the token bucket\n", func() {
			cntSuccessfulInvoking := 0
			cntNoTokenErr := 0
			//	t1 := time.Now()
			var err error
			var token time.Time
			for i := 0; i < 5; i++ {
				token, err = TryToGetToken(bucket)
				//put the service logic here
				if err == nil {
					fmt.Println(token)
					cntSuccessfulInvoking++
				} else if err == ErrorNoToken {
					cntNoTokenErr++
					fmt.Println(err)
				}
			}
			time.Sleep(time.Millisecond * 1100)
			token, err = TryToGetToken(bucket)
			if err == nil {
				fmt.Println(token)
				cntSuccessfulInvoking++
			} else if err == ErrorNoToken {
				cntNoTokenErr++
				fmt.Println(err)
			}
			Convey("Then the request rate is controled", func() {
				fmt.Println(cntSuccessfulInvoking)
				fmt.Println(cntNoTokenErr)
				So(cntSuccessfulInvoking, ShouldEqual, 4)
				So(cntNoTokenErr, ShouldEqual, 2)
			})
		})

	})
}

func TestForVerfyingTokenBucket(t *testing.T) {
	bucket := CreateTokenBucket(3, 1, time.Second*1)
	var err error
	var token time.Time
	cntToken := 0
	// clean the bucket firstly
	for i := 0; i < 3; i++ {
		_, err = TryToGetToken(bucket)
		if err == nil {
			cntToken++
		}
	}
	if cntToken != 3 {
		t.Errorf("Should get 3 token, but actual number of token is %d", cntToken)
	}

	time.Sleep(time.Second * 6)
	t1 := time.Now()
	for j := 0; j < 6; {
		token, err = TryToGetToken(bucket)
		nowSecond := time.Now().Second()
		if err == nil {
			j++
			fmt.Printf("Token put at %v and got at %v\n", token, nowSecond)
		} else {
			time.Sleep(time.Millisecond * 1)
		}
	}
	time_escaped := time.Since(t1)
	if time_escaped < 3 {
		t.Error("logic is wrong. Token is blocked when the bucket is full")
	}
}
