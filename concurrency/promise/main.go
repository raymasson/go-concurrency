package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	po := new(PurchaseOrder)
	po.Value = 42.27

	SavePO(po, false).Then(func(obj interface{}) error {
		po := obj.(*PurchaseOrder)
		fmt.Printf("Purchase Order saved with ID: %d\n", po.Number)

		return nil
		//return errors.New("First promise failed")
	}, func(err error) {
		fmt.Printf("Failed to save Purchase order: " + err.Error() + "\n")
	}).Then(func(obj interface{}) error {
		fmt.Println("Second promise success")

		return nil
	}, func(err error) {
		fmt.Println("Second promise failed: " + err.Error())
	})

	fmt.Scanln()
}

type Promise struct {
	SuccessChannel chan interface{}
	FailureChannel chan error
}

func (this *Promise) Then(success func(interface{}) error, failure func(error)) *Promise {
	result := new(Promise)

	result.SuccessChannel = make(chan interface{}, 1)
	result.FailureChannel = make(chan error, 1)

	timeout := time.After(1 * time.Second)
	go func() {
		select {
		case obj := <-this.SuccessChannel:
			newErr := success(obj)
			if newErr == nil {
				result.SuccessChannel <- obj
			} else {
				result.FailureChannel <- newErr
			}
		case err := <-this.FailureChannel:
			failure(err)
			result.FailureChannel <- err
		case <-timeout:
			failure(errors.New("Promise timed out"))
		}
	}()

	return result
}

type PurchaseOrder struct {
	Number int
	Value  float64
}

func SavePO(po *PurchaseOrder, shouldFail bool) *Promise {
	result := new(Promise)

	result.SuccessChannel = make(chan interface{}, 1)
	result.FailureChannel = make(chan error, 1)

	go func() {
		time.Sleep(2 * time.Second)
		if shouldFail {
			result.FailureChannel <- errors.New("Failed to save purchase order!!!\n")
		} else {
			po.Number = 1234
			result.SuccessChannel <- po
		}
	}()

	return result
}
