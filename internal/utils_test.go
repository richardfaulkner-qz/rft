package utils

import (
	"os"
	"testing"
)

func TestGetUserInput(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "User types enter", args: args{s: ""}, want: "", wantErr: false},
		{name: "User types a string", args: args{s: "Hello"}, want: "Hello World", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			_, err = w.Write([]byte(tt.want))
			w.Close()

			defer func(v *os.File) { os.Stdin = v }(os.Stdin)

			os.Stdin = r

			got, err := GetUserInput(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetYNUserInput(t *testing.T) {
	tests := []struct {
		name      string
		userInput string
		want      bool
		wantErr   bool
	}{
		{name: "User types enter", userInput: "", want: true, wantErr: false},
		{name: "User types Y", userInput: "Y", want: true, wantErr: false},
		{name: "User types Y", userInput: "Y", want: true, wantErr: false},
		{name: "User types N", userInput: "N", want: false, wantErr: false},
		{name: "User types n", userInput: "n", want: false, wantErr: false},
		{name: "User types a char not in the list", userInput: "k", want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A cool lil way to mock user input to os.Stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			_, err = w.Write([]byte(tt.userInput))
			if err != nil {
				t.Fatal(err)
			}

			w.Close()

			defer func(v *os.File) {
				os.Stdin = v
			}(os.Stdin)

			os.Stdin = r

			got, err := GetYNUserInput()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetYNUserInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetYNUserInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
