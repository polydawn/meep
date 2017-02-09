package meep_test

import (
	"encoding/json"
	"fmt"

	"."
)

func ExampleSigh() {
	type Envelope struct {
		Protocol string
		MsgType  string
		Msg      interface{}
	}
	type AppleMsg struct{ Opacity int }
	type PearMsg struct{ Cambre string }

	type ErrBadProtocolHandshake struct {
		meep.TraitCausable
		meep.TraitAutodescribing
	}
	type ErrMalformedMessage struct {
		meep.TraitCausable
		meep.TraitAutodescribing
		ExpectedType string
	}

	envelopeRaw := []byte(`{"Protocol":"7", "MsgType":"apple", "Msg":{"Opacity":"stringy"}}`)

	_, err := func() (*Envelope, error) {
		msgRaw := json.RawMessage{}
		msgEnvelope := &Envelope{Msg: &msgRaw}
		if err := json.Unmarshal(envelopeRaw, msgEnvelope); err != nil {
			return nil, meep.New(
				&ErrBadProtocolHandshake{},
				meep.Cause(err),
			)
		}
		var msg interface{}
		switch msgEnvelope.MsgType {
		case "apple":
			msg = &AppleMsg{}
		case "pear":
			msg = &PearMsg{}
		default:
			return msgEnvelope, meep.New(
				&ErrBadProtocolHandshake{},
				meep.Cause(fmt.Errorf("unknown message type")),
			)
		}
		if err := json.Unmarshal(msgRaw, msg); err != nil {
			return msgEnvelope, meep.New(
				&ErrMalformedMessage{ExpectedType: msgEnvelope.MsgType},
				meep.Cause(err),
			)
		}
		msgEnvelope.Msg = msg
		return msgEnvelope, nil
	}()

	fmt.Printf("%s\n", err)

	// Output:
	//
	// Error[meep_test.ErrMalformedMessage]: ExpectedType="apple";
	// 	Caused by: json: cannot unmarshal string into Go value of type int
	//
}
