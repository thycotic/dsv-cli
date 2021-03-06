package cmd

import (
	e "errors"
	"testing"
	cst "thy/constants"
	"thy/errors"
	"thy/fake"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHandleConfigUpdateCmd(t *testing.T) {
	testCase := []struct {
		name        string
		args        []string
		data        []byte
		out         []byte
		expectedErr *errors.ApiError
	}{
		{
			"Happy path",
			[]string{"conf"},
			[]byte(`test`),
			[]byte(`test`),
			nil},
		{
			"UTF16",
			[]string{"conf"},
			[]byte{255, 254, 91, 0, 83, 0, 99, 0, 114, 0, 105, 0, 112, 0, 116, 0, 32, 0, 73, 0, 110, 0, 102, 0, 111, 0, 93, 0, 13, 0},
			[]byte(`test`),
			nil,
		},
		{
			"No input data",
			[]string{"@conf"},
			[]byte(""),
			nil,
			nil,
		},
		{
			"invalid UTF16",
			[]string{"update", "@conf"},
			[]byte("0xFEFFa\xc5z"),
			[]byte(`test`),
			nil,
		},
	}

	_, err := GetConfigUpdateCmd()
	assert.Nil(t, err)

	viper.Set(cst.Version, "v1")
	for _, tt := range testCase {

		t.Run(tt.name, func(t *testing.T) {
			writer := &fake.FakeOutClient{}
			var data []byte
			var err *errors.ApiError
			writer.WriteResponseStub = func(bytes []byte, apiError *errors.ApiError) {
				data = bytes
				err = apiError
			}

			writer.FailEStub = func(apiError *errors.ApiError) {
				writer.WriteResponseStub(nil, apiError)
			}

			viper.Set(cst.Data, string(tt.data))

			req := &fake.FakeClient{}
			req.DoRequestStub = func(s string, s2 string, i interface{}) (bytes []byte, apiError *errors.ApiError) {
				return tt.out, tt.expectedErr
			}

			r := Config{req, writer, nil, nil}
			_ = r.handleConfigUpdateCmd(tt.args)
			if tt.expectedErr == nil {
				assert.Equal(t, data, tt.out)
			} else {
				assert.Equal(t, err, tt.expectedErr)
			}

			viper.Set(cst.Data, "")
		})
	}
}

func TestHandleConfigReadCmd(t *testing.T) {
	testCase := []struct {
		name        string
		args        string
		apiResponse []byte
		out         []byte
		expectedErr *errors.ApiError
	}{
		{
			"Happy path",
			"conf",
			[]byte(`test`),
			[]byte(`test`),
			nil},
		{
			"api Error",
			"conf",
			[]byte(`test`),
			[]byte(`test`),
			errors.New(e.New("error")),
		},
	}

	_, err := GetConfigReadCmd()
	assert.Nil(t, err)

	viper.Set(cst.Version, "v1")
	for _, tt := range testCase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writer := &fake.FakeOutClient{}
			var data []byte
			var err *errors.ApiError
			writer.WriteResponseStub = func(bytes []byte, apiError *errors.ApiError) {
				data = bytes
				err = apiError
			}

			req := &fake.FakeClient{}
			req.DoRequestStub = func(s string, s2 string, i interface{}) (bytes []byte, apiError *errors.ApiError) {
				return tt.out, tt.expectedErr
			}

			c := Config{req, writer, nil, nil}
			_ = c.handleConfigReadCmd([]string{tt.args})
			if tt.expectedErr == nil {
				assert.Equal(t, data, tt.out)
			} else {
				assert.Equal(t, err, tt.expectedErr)
			}
		})

	}
}

func TestHandleConfigEditCmd(t *testing.T) {
	testCase := []struct {
		name         string
		arg          []string
		out          []byte
		editResponse []byte
		expectedErr  *errors.ApiError
		apiError     *errors.ApiError
		editError    *errors.ApiError
	}{
		{
			"Happy Path",
			[]string{"path1"},
			[]byte(`test data`),
			[]byte(`test data`),
			nil,
			nil,
			nil,
		},
		{
			"error get permission",
			[]string{"path1"},
			[]byte(`test data`),
			nil,
			errors.New(e.New("error")),
			errors.New(e.New("error")),
			nil,
		},
	}

	_, err := GetConfigEditCmd()
	assert.Nil(t, err)

	for _, tt := range testCase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writer := &fake.FakeOutClient{}
			var data []byte
			var err *errors.ApiError
			writer.WriteResponseStub = func(bytes []byte, apiError *errors.ApiError) {
				data = bytes
				err = apiError
			}

			req := &fake.FakeClient{}
			req.DoRequestStub = func(s string, s2 string, i interface{}) (bytes []byte, apiError *errors.ApiError) {
				return tt.out, tt.apiError
			}

			c := Config{req, writer, nil, nil}
			c.edit = func(bytes2 []byte, d dataFunc, apiError *errors.ApiError, retry bool) (bytes []byte, apiError2 *errors.ApiError) {
				_, _ = d([]byte(`config`))
				return tt.editResponse, tt.editError
			}

			_ = c.handleConfigEditCmd(tt.arg)
			if tt.expectedErr == nil {
				assert.Equal(t, data, tt.out)
			} else {
				assert.Equal(t, err, tt.expectedErr)
			}
		})
	}
}

func TestGetConfigCmd(t *testing.T) {
	_, err := GetConfigCmd()
	assert.Nil(t, err)
}
