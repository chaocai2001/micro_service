/*
	This package is to provide the basic mechanisms,
	which are the fundations for building the high reliable microservice application

	Author: Chao Cai

*/
package microservice_helper

import (
	"errors"
	"time"
)

var ErrorGettingTokenTimeout = errors.New("Failed to get token for timeout.")

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
		for t := range time.Tick(tokenFillingInterval) {
			for i := 0; i < numOfTokens; i++ {
				bucket <- t
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
	token, err = GetToken(bucket, WAIT_ANYWAY) //set the timeout or waiting any anyway

	if err == nil {
		//put the service logic here
	}
}
