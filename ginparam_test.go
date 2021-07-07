package ginparam

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRead(t *testing.T) {
	type nestedOut struct {
		Nested string `query:"nested_val"`
	}
	type out struct {
		StrVal         string    `query:"str_val"`
		IntVal         int       `query:"int_val"`
		BoolVal        bool      `query:"bool_val"`
		StrSliceVal    []string  `query:"str_slice_val"`
		IntSliceVal    []int     `query:"int_slice_val"`
		BoolSliceVal   []bool    `query:"bool_slice_val"`
		IgnoredNestVal nestedOut `query:"-"`
		NestVal        nestedOut
	}
	type args struct {
		params string

		// out must be assignable(pointer or interface)
		out interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "empty param", args: args{params: "", out: &out{StrVal: "abc"}}, want: &out{StrVal: "abc"}, wantErr: false},
		{name: "string value", args: args{params: "?str_val=abc", out: &out{}}, want: &out{StrVal: "abc"}, wantErr: false},
		{name: "int value", args: args{params: "?int_val=123", out: &out{}}, want: &out{IntVal: 123}, wantErr: false},
		{name: "int value (N/A)", args: args{params: "?int_val=abc", out: &out{}}, wantErr: true},
		{name: "bool value(true)", args: args{params: "?bool_val=true", out: &out{}}, want: &out{BoolVal: true}, wantErr: false},
		{name: "bool value(false)", args: args{params: "?bool_val=false", out: &out{BoolVal: true}}, want: &out{BoolVal: false}, wantErr: false},
		{name: "bool value(empty)", args: args{params: "", out: &out{BoolVal: true}}, want: &out{BoolVal: true}, wantErr: false},
		{name: "override default string value", args: args{params: "?str_val=ABC", out: &out{StrVal: "abc"}}, want: &out{StrVal: "ABC"}, wantErr: false},
		{name: "struct value", args: args{params: "?nested_val=abc", out: &out{}}, want: &out{NestVal: nestedOut{Nested: "abc"}}, wantErr: false},
		{name: "string slice value", args: args{params: "?str_slice_val=a,b,c", out: &out{}}, want: &out{StrSliceVal: []string{"a", "b", "c"}}, wantErr: false},
		{name: "int slice value", args: args{params: "?int_slice_val=1,2,3", out: &out{}}, want: &out{IntSliceVal: []int{1, 2, 3}}, wantErr: false},
		{name: "int slice value(N/A)", args: args{params: "?int_slice_val=a,b,c", out: &out{}}, wantErr: true},
		{name: "bool slice value", args: args{params: "?bool_slice_val=true,false,ABC", out: &out{}}, want: &out{BoolSliceVal: []bool{true, false, false}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			reqPath := fmt.Sprint("/test", tt.args.params)
			req := httptest.NewRequest("GET", reqPath, nil)
			ctx.Request = req

			// Run test
			// In this case, tt.args.out wraps assignable value.
			// So pass its value directly(instead of its pointer redundantly).
			err := Read(ctx, &tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Test for nominal case
			if err != nil {
				return
			}

			expect := tt.want
			actual := tt.args.out
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("Read() doesn't return expected result.\nExpect = %+v,\nActual = %+v", expect, actual)
			}
		})
	}
}
