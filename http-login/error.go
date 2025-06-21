package main

type RequestError struct {
	Body     string
	HTTPCode string
	Err      string
}

func (r RequestError) Error() string {
	return r.Err
}
