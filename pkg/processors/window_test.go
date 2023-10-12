package processors_test

import (
	"reflect"
	"testing"

	"github.com/primevprotocol/oracle/pkg/chaintracer"
	"github.com/primevprotocol/oracle/pkg/processors"
)

func TestProcessBlock(t *testing.T) {
	tests := []struct {
		name       string
		block      chaintracer.BlockDetails
		openBids   processors.BidSet
		wantOpen   processors.BidSet
		wantClosed processors.BidSet
		wantErr    bool
	}{
		{
			name: "Normal case",
			block: chaintracer.BlockDetails{
				Transactions: []string{"txn1", "txn2", "txn3"},
			},
			openBids: map[string]struct{}{
				"txn1": {},
				"txn2": {},
				"txn4": {},
			},
			wantOpen: map[string]struct{}{
				"txn4": {},
			},
			wantClosed: map[string]struct{}{
				"txn1": {},
				"txn2": {},
			},
			wantErr: false,
		},
		{
			name: "No closed bids",
			block: chaintracer.BlockDetails{
				Transactions: []string{"txn3", "txn5"},
			},
			openBids: map[string]struct{}{
				"txn1": {},
				"txn2": {},
			},
			wantOpen: map[string]struct{}{
				"txn1": {},
				"txn2": {},
			},
			wantClosed: map[string]struct{}{},
			wantErr:    false,
		},
		{
			name: "No open bids",
			block: chaintracer.BlockDetails{
				Transactions: []string{"txn1", "txn2"},
			},
			openBids: map[string]struct{}{
				"txn1": {},
				"txn2": {},
			},
			wantOpen: map[string]struct{}{},
			wantClosed: map[string]struct{}{
				"txn1": {},
				"txn2": {},
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	w := processors.WindowAlgo{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOpen, gotClosed, err := w.ProcessBlock(tt.block, tt.openBids)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOpen, tt.wantOpen) {
				t.Errorf("ProcessBlock() gotOpen = %v, want %v", gotOpen, tt.wantOpen)
			}
			if !reflect.DeepEqual(gotClosed, tt.wantClosed) {
				t.Errorf("ProcessBlock() gotClosed = %v, want %v", gotClosed, tt.wantClosed)
			}
		})
	}
}
