package inmemory

import (
	"reflect"
	"testing"
)

func TestBoard_Get(t *testing.T) {
	type args struct {
		width  uint
		height uint
	}
	type testCase[T any] struct {
		name    string
		b       Board[T]
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name:    "valid coordinates",
			b:       NewBoard[int](3, 3),
			args:    args{width: 1, height: 1},
			want:    0,
			wantErr: false,
		},
		{
			name:    "out of bounds",
			b:       NewBoard[int](3, 3),
			args:    args{width: 4, height: 1},
			want:    0,
			wantErr: true,
		},
		{
			name:    "out of bounds corner case",
			b:       NewBoard[int](3, 3),
			args:    args{width: 3, height: 2},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Get(tt.args.width, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoard_Height(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Board[T]
		want uint
	}
	tests := []testCase[int]{
		{
			name: "height 3",
			b:    NewBoard[int](3, 3),
			want: 3,
		},
		{
			name: "height 5",
			b:    NewBoard[int](4, 5),
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Height(); got != tt.want {
				t.Errorf("Height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoard_Set(t *testing.T) {
	type args[T any] struct {
		width  uint
		height uint
		t      T
	}
	type testCase[T any] struct {
		name    string
		b       Board[T]
		args    args[T]
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name:    "valid set",
			b:       NewBoard[int](3, 3),
			args:    args[int]{width: 1, height: 1, t: 42},
			wantErr: false,
		},
		{
			name:    "out of bounds set",
			b:       NewBoard[int](3, 3),
			args:    args[int]{width: 4, height: 1, t: 42},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Set(tt.args.width, tt.args.height, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoard_Width(t *testing.T) {
	type testCase[T any] struct {
		name string
		b    Board[T]
		want uint
	}
	tests := []testCase[int]{
		{
			name: "width 3",
			b:    NewBoard[int](3, 3),
			want: 3,
		},
		{
			name: "width 4",
			b:    NewBoard[int](4, 5),
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Width(); got != tt.want {
				t.Errorf("Width() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_coordsAreInBounds(t *testing.T) {
	type args[T any] struct {
		width  uint
		height uint
		grid   *[][]T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "in bounds",
			args: args[int]{width: 1, height: 1, grid: &[][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}},
			want: true,
		},
		{
			name: "out of bounds width",
			args: args[int]{width: 3, height: 1, grid: &[][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}},
			want: false,
		},
		{
			name: "out of bounds height",
			args: args[int]{width: 1, height: 3, grid: &[][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coordsAreInBounds(tt.args.width, tt.args.height, tt.args.grid); got != tt.want {
				t.Errorf("coordsAreInBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}
