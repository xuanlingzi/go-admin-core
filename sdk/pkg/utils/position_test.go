package utils

import (
	"github.com/spf13/cast"
	"testing"
)

func TestGetDistance(t *testing.T) {
	type args struct {
		tm1  string
		lat1 float64
		lon1 float64
		tm2  string
		lat2 float64
		lon2 float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test1",
			args: args{
				tm1:  "2024-06-05 12:51:32",
				lat1: 34.922394,
				lon1: 118.743530,
				tm2:  "2024-06-05 12:49:32",
				lat2: 34.907707,
				lon2: 118.749016,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := GetDistance(tt.args.lat1, tt.args.lon1, tt.args.lat2, tt.args.lon2)
			duration := cast.ToTime(tt.args.tm1).Sub(cast.ToTime(tt.args.tm2))
			speed := (distance / 1000) / (duration.Seconds() / 3600)
			t.Errorf("Distance: %v, Duration: %v, Speed: %v", distance, duration.Seconds(), speed)
		})
	}
}
