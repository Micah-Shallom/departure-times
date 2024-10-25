package nextbus

import (
	"github.com/Micah-Shallom/departure-times/external"
	"github.com/Micah-Shallom/departure-times/utility"
)

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod string
	Logger       *utility.Logger
}

func (r *RequestObj) getNewSendRequestObject(data interface{}, headers map[string]string, urlprefix string) *external.SendRequestObject {
	return external.GetNewSendRequestObject(
		r.Logger,
		r.Name,
		r.Path,
		r.Method,
		urlprefix,
		r.DecodeMethod,
		headers,
		r.SuccessCode,
		data,
	)
}
