package common

type NamedCallback struct {
	name     string
	callback interface{}
}

func NewNamedCallback(name string, callback interface{}) *NamedCallback {
	return &NamedCallback{
		name:     name,
		callback: callback,
	}
}

func (cb *NamedCallback) Name() string {
	return cb.name
}

func (cb *NamedCallback) Callback() interface{} {
	return cb.callback
}

func (cb *NamedCallback) SetCallback(callback interface{}) interface{} {
	oriCallback := cb.callback
	cb.callback = callback
	return oriCallback
}

////

func AddCallback(callbacks []*NamedCallback, name string, callback interface{}, rewrite bool) (
	[]*NamedCallback, interface{}, bool) {

	var oriCallback interface{}
	for _, namedCallback := range callbacks {
		if namedCallback.Name() == name {
			if rewrite {
				oriCallback = namedCallback.SetCallback(callback)
			} else {
				return callbacks, namedCallback.Callback(), false
			}
		}
	}

	if oriCallback == nil {
		callbacks = append(callbacks, NewNamedCallback(name, callback))
	}

	return callbacks, oriCallback, true
}

func DeleteCallback(callbacks []*NamedCallback, name string) ([]*NamedCallback, interface{}) {
	var oriCallback interface{}
	for i, namedCallback := range callbacks {
		if namedCallback.Name() == name {
			oriCallback = namedCallback.Callback()
			callbacks = append(callbacks[:i], callbacks[i+1:]...)
			break
		}
	}

	return callbacks, oriCallback
}