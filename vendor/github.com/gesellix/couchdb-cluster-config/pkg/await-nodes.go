package cluster_config

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func AwaitNodes(addresses []string, timeout <-chan time.Time, check func(address string) (bool, error)) error {
	resc, errc := make(chan string), make(chan error)

	for _, address := range addresses {
		go func(address string) {
			success, err := awaitNode(address, timeout, check)
			if err != nil {
				errc <- err
				return
			}
			resc <- fmt.Sprintf("%v@%s", success, address)
		}(address)
	}

	for i := 0; i < len(addresses); i++ {
		select {
		case res := <-resc:
			fmt.Println(res)
		case err := <-errc:
			//fmt.Println(err)
			return err
		}
	}
	time.Sleep(5 * time.Second)

	return nil
}

func awaitNode(address string, timeout <-chan time.Time, check func(address string) (bool, error)) (bool, error) {
	tick := time.Tick(1000 * time.Millisecond)
	for {
		select {
		case <-timeout:
			//fmt.Println("timeout")
			return false, errors.New(fmt.Sprintf("timeout@%s", address))
		case <-tick:
			fmt.Println(fmt.Sprintf("tick@%s", address))

			ok, err := check(address)
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
		}
	}
}

func Available(address string) (bool, error) {
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() || err.Temporary() {
			return false, nil
		}
		return false, nil
		//return false, err
	}
	defer conn.Close()
	return true, nil
}
