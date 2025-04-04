package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	// kiem tra mat khau co khop khong
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	//cap nhat mat khau thuc te de so sanh cac truong khac
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password (%v)", e.arg, e.password)
}

func EqCreateUserParamsMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	hashPassword, err := util.HashPassword(password)

	require.NoError(t, err)
	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     user.Email,
				"full_Name": user.FullName,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: hashPassword, // Đảm bảo hashPassword được thiết lập đúng
					FullName:       user.FullName,
					Email:          user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParamsMatcher(arg, password)).
					Times(1).
					Return(user, nil)

			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)

			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)

			//start server
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/users")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(recorder)
		})
	}

}
func randomUser(t *testing.T) (user db.Users, password string) {
	password = util.RandomString(6) // Tạo mật khẩu ngẫu nhiên
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err) // Kiểm tra lỗi trong test

	user = db.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword, //mat khau chua duoc hash
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return user, password // Trả về user và mật khẩu chưa băm
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.Users) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.Users
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	// So sánh các trường mà bạn muốn kiểm tra, không bao gồm mật khẩu
	// Bạn có thể loại bỏ so sánh mật khẩu ở đây
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
	require.Equal(t, user.CreatedAt, gotUser.CreatedAt)
}
