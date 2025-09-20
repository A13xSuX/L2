package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil { //because we have pointer on structer, type not nil
		println("error")
		//fmt.Printf("%T", err)
		return
	}
	println("ok")
}
