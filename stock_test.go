package main

import "testing"

func Test_updateStockCode(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test1", args{"삼성전자"}, "005930"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateStockCode(tt.args.name); got != tt.want {
				t.Errorf("updateStockCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
