package dispatcher

import (
	"encoding/json"
	"fmt"
)

type callMessage struct {
	Name       string            `json:"name"`
	Args       []json.RawMessage `json:"args"`
	CallbackID string            `json:"callbackID"`
}

func (d *Dispatcher) processCallMessage(message string) error {

	var payload callMessage
	err := json.Unmarshal([]byte(message[1:]), &payload)
	if err != nil {
		return err
	}
	// Lookup method
	registeredMethod := d.bindingsDB.GetMethod(payload.Name)

	// Check we have it
	if registeredMethod == nil {
		return fmt.Errorf("method '%s' not registered", payload.Name)
	}

	args, err := registeredMethod.ParseArgs(payload.Args)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %s", err.Error())
	}

	result, err := registeredMethod.Call(args)
	callbackMessage := &CallbackMessage{
		CallbackID: payload.CallbackID,
	}
	if err != nil {
		callbackMessage.Err = err.Error()
	} else {
		callbackMessage.Result = result
	}
	messageData, err := json.Marshal(callbackMessage)
	d.log.Trace("json call result data: %+v\n", string(messageData))
	if err != nil {
		// what now?
		d.log.Fatal(err.Error())
	}
	d.resultCallback(string(messageData))

	return nil
}

// CallbackMessage defines a message that contains the result of a call
type CallbackMessage struct {
	Result     interface{} `json:"result"`
	Err        string      `json:"error"`
	CallbackID string      `json:"callbackid"`
}
