package customerror

type VideoURLNotFound struct {
	Err string
}

func (e VideoURLNotFound) Error() string {
	return e.Err
}

type XSRFTokenNotFound struct {
	Err string
}

func (e XSRFTokenNotFound) Error() string {
	return e.Err
}

type MissingVideoUrl struct {
	Err string
}

func (e MissingVideoUrl) Error() string {
	return e.Err
}
