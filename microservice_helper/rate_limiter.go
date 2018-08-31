package microservice_helper

import (
	"errors"
	"time"
)

// WAIT_ANYWAY is to disable the waiting timeout
const WAIT_ANYWAY time.Duration = 0

// ErrorGettingTokenTimeout occured when waiting token timeout
var ErrorGettingTokenTimeout = errors.New("Failed to get token for timeout.")

// ErrorNoToken occured when invoking TryToGetToken and no token is available
// that moment
var ErrorNoToken = errors.New("Failed to get token")

// TryToGetToken is try to get a token, the function would be return immediately.
// If the token is not ready, an error (ErrorNoToken) will be thrown.
func TryToGetToken(tokenBucket chan time.Time) (time.Time, error) {
	var token time.Time
	select {
	case token = <-tokenBucket:
		return token, nil
	default:
		return token, ErrorNoToken
	}
}

// GetToken is to get a token, if the token cann't be ready in the setting timeout duration,
// the timeout error (ErrorGettingTokenTimeout) will be thrown.
func GetToken(tokenBucket chan time.Time,
	timeout time.Duration) (time.Time, error) {
	var token time.Time

	if timeout != 0 {
		select {
		case token = <-tokenBucket:
			return token, nil
		case <-time.After(timeout):
			return token, ErrorGettingTokenTimeout
		}
	} else {
		token = <-tokenBucket
		return token, nil
	}

}

//Create a token bucket
//sizeOfBucket is the size of the leaky bucket
//numOfTokens, tokenFillingInterval are used to rate limit,
//the (numOfTokens) tokes would be put into bucket in the  period (tokenFillingInterval)
func CreateTokenBucket(sizeOfBucket int, numOfTokens int,
	tokenFillingInterval time.Duration) chan time.Time {
	bucket := make(chan time.Time, sizeOfBucket)
	//fill the bucket firstly
	for j := 0; j < sizeOfBucket; j++ {
		bucket <- time.Now()
	}
	go func() {
		for t := range time.Tick(tokenFillingInterval / time.Duration(numOfTokens)) {
			select {
			case bucket <- t:
			default:
			}
		}
	}()
	return bucket
}

func ExampleRateLimit() {
	//Create a token bucket
	//sizeOfBucket is the size of the leaky bucket
	//numOfTokens, tokenFillingInterval are used to rate limit,
	//the (numOfTokens) tokes would be put into bucket in the  period (tokenFillingInterval)
	//func CreateTokenBucket(sizeOfBucket int, numOfTokens int,tokenFillingInterval time.Duration) chan time.Time
	bucket := CreateTokenBucket(3, 2, time.Second*1)
	_, err := GetToken(bucket, WAIT_ANYWAY) //set the timeout or waiting any anyway

	if err == nil {
		//put the service logic here
	}
}
