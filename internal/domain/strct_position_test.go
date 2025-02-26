// internal/domain/strct_position_test.go
package domain

import (
	"testing"
	"time"
)

func TestPosition_Validate(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		wantErr  bool
	}{
		{
			name: "valid position",
			position: Position{
				Symbol:   "7203",
				Side:     "long",
				Quantity: 100,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid position - empty Symbol",
			position: Position{
				Side:     "long",
				Quantity: 100,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid position - invalid Side",
			position: Position{
				Symbol:   "7203",
				Side:     "invalid",
				Quantity: 100,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid position - zero Quantity",
			position: Position{
				Symbol:   "7203",
				Side:     "long",
				Quantity: 0,
				Price:    1500,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid position - zero Price",
			position: Position{
				Symbol:   "7203",
				Side:     "long",
				Quantity: 100,
				Price:    0,
				OpenDate: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.position.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
