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
