package board

import (
	"testing"
)

// func Test_fill_arrival(t *testing.T) {
// 	type args struct {
// 		entity  *gtfs.FeedEntity
// 		line    string
// 		stop_id string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want Arrival
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for i, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := fill_arrival(tt.args.entity, tt.args.line, tt.args.stop_id, i); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("fill_arrival() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_request(t *testing.T) {
// 	type args struct {
// 		line string
// 		stop string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []Arrival
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := request(tt.args.line, tt.args.stop)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("request() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("request() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGetArrivals(t *testing.T) {
// 	type args struct {
// 		line string
// 		stop string
// 		N    int
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []Arrival
// 	}{
// 		{name: "test1", args: args{line: "6", stop: "635S", N: 1000}, want: []Arrival{}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GetArrivals(tt.args.line, tt.args.stop, tt.args.N); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetArrivals() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_init_linestoURI_linestoColors(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			init_linestoURI_linestoColors()
// 		})
// 	}
// }

func Test_init_stopstoFriendlies(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init_stopstoFriendlies()

		})
	}
}
