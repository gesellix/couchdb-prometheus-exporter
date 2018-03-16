package cluster_setup

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func AwaitNodes(urls []string) error {
	resc, errc := make(chan bool), make(chan error)

	for _, u := range urls {
		go func(u string) {
			success, err := awaitNode(u)
			if err != nil {
				errc <- err
				return
			}
			resc <- success
		}(u)
	}

	for i := 0; i < len(urls); i++ {
		select {
		case res := <-resc:
			fmt.Println(res)
		case err := <-errc:
			//fmt.Println(err)
			return err
		}
	}
	return nil
}

func awaitNode(url string) (bool, error) {
	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			fmt.Println("timeout")
			return false, errors.New(fmt.Sprintf("timed out @%s", url))
		case <-tick:
			fmt.Println(fmt.Sprintf("tick@%s", url))

			ok, err := fetch(url)
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
		}
	}
}

func fetch(url string) (bool, error) {
	conn, err := net.DialTimeout("tcp", url, time.Second)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return false, nil
		}
		return false, err
	}
	defer conn.Close()
	return true, nil

	//res, err := http.Get(url)
	//if err != nil {
	//	return false, err
	//}
	//_, err = ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//if err != nil {
	//	return false, err
	//}
	//return true, nil
}
