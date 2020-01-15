package artisan

import "fmt"

type FuncDescription struct {
	Receiver *FuncTmplReceiver
	Name     string
}

func (f FuncDescription) String() string {
	if f.Receiver == nil {
		return f.Name
	}

	return fmt.Sprintf("(%s %s) %s", f.Receiver.Name, f.Receiver.TypeOf, f.Name)
}
