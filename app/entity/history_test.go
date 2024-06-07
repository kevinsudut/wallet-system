package entity

import (
	"testing"
)

func TestHistory_NormalizeAmount(t *testing.T) {
	type fields struct {
		Id           string
		UserId       string
		TargetUserId string
		Amount       float64
		Type         int
		Notes        string
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "+amount",
			fields: fields{
				Type:   2,
				Amount: 100,
			},
			want: -100,
		},
		{
			name: "-amount",
			fields: fields{
				Type:   2,
				Amount: -100,
			},
			want: -100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &History{
				Id:           tt.fields.Id,
				UserId:       tt.fields.UserId,
				TargetUserId: tt.fields.TargetUserId,
				Amount:       tt.fields.Amount,
				Type:         tt.fields.Type,
				Notes:        tt.fields.Notes,
			}
			h.NormalizeAmount()
			if got := h.Amount; got != tt.want {
				t.Errorf("History.Amount = %v, want %v", got, tt.want)
			}
		})
	}
}
