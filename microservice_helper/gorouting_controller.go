package microservice_helper

import (
	"errors"
	"time"
)

type token struct{}

var TokenPoolNotAvailableError error = errors.New("Token pool is not available.")
var TokenWaitingTimeoutError error = errors.New("Token waiting is time out.")

type GoroutingController struct {
	maxNumOfRoutings int
	tokenPools       chan token
}

func CreateGoroutingController(maxNumOfRoutings int) *GoroutingController {
	controller := GoroutingController{maxNumOfRoutings, make(chan token, maxNumOfRoutings)}
	for i := 0; i < maxNumOfRoutings; i++ {
		controller.tokenPools <- token{}
	}
	return &controller
}

func (controller *GoroutingController) ApplyToken(waitingTimeout time.Duration) error {
	select {
	case _, ok := <-controller.tokenPools:
		if !ok {
			return TokenPoolNotAvailableError
		}
		return nil
	case <-time.After(waitingTimeout): //timeout mechanism is to avoid of being blocked forever/gorouting leaking
		return TokenWaitingTimeoutError
	}
}

func (controller *GoroutingController) ReleaseToken() {
	controller.tokenPools <- token{}
}

func (controller *GoroutingController) NumOfGoroutingCanBeCreated() int {
	return len(controller.tokenPools)
}

//The function is a little bit tricky, if you don't really understand why declare innerParam, please don't use it
func (controller *GoroutingController) StartGorouting(
	f func(param interface{}), waitingTimeout time.Duration, innerParam interface{}) error {
	if err := controller.ApplyToken(waitingTimeout); err == nil {
		go func() {
			defer func() { controller.tokenPools <- token{} }()
			f(innerParam)

		}()
	} else {
		return err
	}

	return nil
}

func ExampleGroutingNumberControl() error {
	maxNumOfGrouting := 10
	controller := CreateGoroutingController(maxNumOfGrouting) //create a controller
	var err error
	for i := 0; i < maxNumOfGrouting+1; i++ {
		//Get a token, before starting a new gorouting
		err = controller.ApplyToken(time.Second * 2)
		//set 'time.Second * 2' as the longest waiting time for getting the token
		//if can't get the toke in 2 seconds, a timeout error will be returned
		if err == nil {
			go func(d int) {
				defer controller.ReleaseToken() //release the token when gorouting is over
				time.Sleep(time.Second * 5)
			}(i)
		}
	}
	return err
	//Output: Token waiting is time out.
}
