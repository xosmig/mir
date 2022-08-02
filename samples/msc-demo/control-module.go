package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"

	availabilityevents "github.com/filecoin-project/mir/pkg/availability/events"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
)

type controlModule struct {
	eventsOut chan *events.EventList
}

func newControlModule() modules.ActiveModule {
	return &controlModule{
		eventsOut: make(chan *events.EventList),
	}
}

func (m *controlModule) ImplementsModule() {}

func (m *controlModule) ApplyEvents(ctx context.Context, events *events.EventList) error {
	iter := events.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		switch event.Type.(type) {

		case *eventpb.Event_Init:
			go func() {
				err := m.readConsole()
				if err != nil {
					panic(err)
				}
			}()

		case *eventpb.Event_Availability:
			availabilityEvent := event.Type.(*eventpb.Event_Availability).Availability
			switch availabilityEvent.Type.(type) {

			case *availabilitypb.Event_NewCert:
				newCertEv := availabilityEvent.Type.(*availabilitypb.Event_NewCert).NewCert
				certBytes, err := proto.Marshal(newCertEv.Cert)
				if err != nil {
					return fmt.Errorf("error marshalling certificate: %w", err)
				}

				fmt.Println(base64.StdEncoding.EncodeToString(certBytes))

			case *availabilitypb.Event_ProvideTransactions:
				txsEv := availabilityEvent.Type.(*availabilitypb.Event_ProvideTransactions).ProvideTransactions
				for _, txBytes := range txsEv.Txs {
					tx := new(requestpb.Request)
					err := proto.Unmarshal(txBytes, tx)
					if err != nil {
						return fmt.Errorf("erro unmarshalling transaction: %w", err)
					}

					fmt.Println(string(tx.Data))
				}
			}

		}
	}

	return nil
}

func (m *controlModule) EventsOut() <-chan *events.EventList {
	return m.eventsOut
}

func (m *controlModule) readConsole() error {
	// Read the user input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Type in the command ('createBatch', 'readBatch')")
		scanner.Scan()
		if scanner.Err() != nil {
			return fmt.Errorf("error reading from console: %w", scanner.Err())
		}

		text := scanner.Text()

		switch cmd := strings.TrimSpace(text); cmd {
		case "createBatch":
			err := m.createBatch(scanner)
			return err

		case "readBatch":
			err := m.readBatch(scanner)
			return err

		default:
			fmt.Println("Unknown command: ", cmd)
		}
	}
}

func (m *controlModule) createBatch(scanner *bufio.Scanner) error {
	fmt.Println("Type in 1 transaction per line, then type 'send!' and press Enter")

	for {
		scanner.Scan()
		if scanner.Err() != nil {
			return fmt.Errorf("error reading user data: %w", scanner.Err())
		}

		text := scanner.Text()
		if strings.TrimSpace(text) == "send!" {
			break
		}

		request := &requestpb.Request{Data: []byte(text)}
		m.eventsOut <- events.ListOf(events.NewClientRequests("mempool", []*requestpb.Request{request}))
	}

	m.eventsOut <- events.ListOf(availabilityevents.RequestCert("availability", &availabilitypb.RequestCertOrigin{
		Module: "control",
		Type:   &availabilitypb.RequestCertOrigin_ContextStore{},
	}))

	return nil
}

func (m *controlModule) readBatch(scanner *bufio.Scanner) error {
	fmt.Println("type in the availability certificate and press Enter")

	scanner.Scan()
	if scanner.Err() != nil {
		return fmt.Errorf("error reading batch id: %w", scanner.Err())
	}

	certBase64 := strings.TrimSpace(scanner.Text())
	certBytes, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return fmt.Errorf("error decoding certificate: %w", err)
	}

	cert := new(availabilitypb.Cert)
	err = proto.Unmarshal(certBytes, cert)
	if err != nil {
		return fmt.Errorf("error unmarshalling certificate: %w", err)
	}

	m.eventsOut <- events.ListOf(availabilityevents.RequestTransactions("availability", cert,
		&availabilitypb.RequestTransactionsOrigin{
			Module: "control",
			Type:   &availabilitypb.RequestTransactionsOrigin_ContextStore{},
		}))

	return nil
}
