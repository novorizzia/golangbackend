package api

import (
	mockdb "backendmaster/db/mock"
	db "backendmaster/db/sqlc"
	"database/sql"
	"io/ioutil"
	"reflect"

	"backendmaster/utils/password"
	"backendmaster/utils/random"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// konversi x menjadi tipe data CreateUserParams
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := password.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, pwd := randomUser(t)
	userName := user.Username
	fullName := user.FullName
	email := user.Email

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(store *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			name: "BadRequest",
			body: gin.H{
				"username":  "!@##$",
				"password":  pwd,
				"full_name": userName,
				"email":     email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			name: "uniqueConstraint",
			body: gin.H{
				"username":  userName,
				"password":  pwd,
				"full_name": userName,
				"email":     email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					// 23505 is a unique constraint violation error code.
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, responseRecorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"username":  userName,
				"password":  pwd,
				"full_name": userName,
				"email":     email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			name: "Ok",
			body: gin.H{
				"username":  userName,
				"password":  pwd,
				"full_name": fullName,
				"email":     email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       userName,
					HashedPassword: user.HashedPassword,
					FullName:       fullName,
					Email:          email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, pwd)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, responseRecorder.Code)
				requireBodyMatchUser(t, responseRecorder.Body, user)
			},
		},
	}

	for i, _ := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := NewServer(store)

			// build STUBS
			testCase.buildStubs(store)

			// buat recorder
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users")

			// marshal body to JSON
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			// http.NewRequest(method, url,requestbody)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			testCase.checkResponse(t, recorder)

		})
	}
}

type fakeUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}

type fakeUserRequest struct {
	Username string
	Password string
	FullName string
	Email    string
}

func randomUser(t *testing.T) (db.User, string) {
	fullname := random.RandomOwner()
	req := fakeUserRequest{
		Username: fullname,
		Password: random.RandomPassword(),
		FullName: fullname + " " + fullname,
		Email:    fullname + "@gmail.com",
	}

	hashPwd, err := password.HashedPassword(req.Password)
	require.NoError(t, err)

	user := db.User{
		Username:       req.Username,
		HashedPassword: hashPwd,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	return user, req.Password

}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
