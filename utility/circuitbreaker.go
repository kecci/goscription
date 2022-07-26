package utility

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/sirupsen/logrus"
)

func CallUsingCircuitBreaker(breakername string, request *http.Request, body []byte) ([]byte, error) {
	hystrix.ConfigureCommand(breakername, hystrix.CommandConfig{
		Timeout:                5000,
		SleepWindow:            5000,
		RequestVolumeThreshold: 10,
	})

	output := make(chan []byte, 1) // declare the channel where the hystrix goroutine will put success responses.

	errors := hystrix.Go(breakername, // pass the name of the circuit breaker as first parameter.

		// 2nd parameter, the inlined func to run inside the breaker.
		func() error {
			// for hystrix, forward the err from the retrier. it's nil if successful.
			return CallWithRetries(request, output)
		},

		// 3rd parameter, the fallback func. in this case, we just do a bit of logging and return the error.
		func(err error) error {
			logrus.Errorf("in fallback function for breaker %v, error: %v", breakername, err.Error())
			circuit, _, _ := hystrix.GetCircuit(breakername)
			logrus.Errorf("circuit state is: %v", circuit.IsOpen())
			return err
		})

	// response and error handling. if the call was successful, the output channel gets the response. otherwise,
	// the errors channel gives us the error.
	select {
	case out := <-output:
		logrus.Debugf("call in breaker %v successful", breakername)
		return out, nil

	case err := <-errors:
		return nil, err
	}
}

func CallWithRetries(req *http.Request, output chan []byte) error {
	// Retries Attempt
	retries := 3

	// create a retrier with constant backoff, retries number of attempts (3) with a 100ms sleep between retries.
	r := retrier.New(retrier.ConstantBackoff(retries, 100*time.Millisecond), nil)

	// this counter is just for getting some logging for showcasing, remove in production code.
	attempt := 0

	// retrier works similar to hystrix, we pass the actual work (doing the http request) in a func.
	err := r.Run(func() error {
		attempt++

		// do http request and handle response. if successful, pass resp.body over output channel,
		// otherwise, do a bit of error logging and return the err.
		var client = &http.Client{}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			responsebody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				defer resp.Body.Close()
				output <- responsebody
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("status was %v", resp.StatusCode)
		}

		logrus.Errorf("retrier failed, attempt %v", attempt)
		return err
	})
	return err
}
