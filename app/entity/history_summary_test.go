package entity

import (
	"testing"
)

func TestHistorySummary_GetId(t *testing.T) {
	type fields struct {
		Id           string
		UserId       string
		TargetUserId string
		Amount       float64
		Type         int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				UserId:       "id",
				TargetUserId: "target",
				Type:         1,
			},
			want: "id:target:1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hs := HistorySummary{
				Id:           tt.fields.Id,
				UserId:       tt.fields.UserId,
				TargetUserId: tt.fields.TargetUserId,
				Amount:       tt.fields.Amount,
				Type:         tt.fields.Type,
			}
			if got := hs.GetId(); got != tt.want {
				t.Errorf("HistorySummary.GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
