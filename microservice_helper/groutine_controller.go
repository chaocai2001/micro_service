package microservice_helper

import (
	"errors"
	"time"
)

type token struct{}

var TokenPoolNotAvailableError error = errors.New("Token pool is not available.")
var TokenWaitingTimeoutError error = errors.New("Token waiting is time out.")

// GoroutingController is the concurrency settings of groutine
type GoroutineController struct {
	maxNumOfRoutines int
	tokenPools       chan token
}

// CreateGoroutineController is to create a concurrency controller
func CreateGoroutineController(maxNumOfRoutings int) *GoroutineController {
	controller := GoroutineController{maxNumOfRoutings, make(chan token, maxNumOfRoutings)}
	for i := 0; i < maxNumOfRoutings; i++ {
		controller.tokenPools <- token{}
	}
	return &controller
}

// ApplyToken is to apply a token. Before starting a groutine, you should apply a token
func (controller *GoroutineController) ApplyToken(waitingTimeout time.Duration) error {
	select {
	case _, ok := <-controller.tokenPools:
		if !ok {
			return TokenPoolNotAvailableError
		}
		return nil
	case <-time.After(waitingTimeout): //timeout mechanism is to avoid of being blocked forever/goroutine leaking
		return TokenWaitingTimeoutError
	}
}

// ReleaseToken is to release the token. It should be invoked when a groutine is over
func (controller *GoroutineController) ReleaseToken() {
	controller.tokenPools <- token{}
}

// NumOfGoroutingCanBeCreated is to get the number of groutine can be created
func (controller *GoroutineController) NumOfGoroutingCanBeCreated() int {
	return len(controller.tokenPools)
}

// StartGorouting is a little bit tricky, if you don't really understand why declare innerParam, please don't use it
func (controller *GoroutineController) StartGorouting(
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
	controller := CreateGoroutineController(maxNumOfGrouting) //create a controller
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
