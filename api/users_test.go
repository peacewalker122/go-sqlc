package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "sqlc/db/mock"
	db "sqlc/db/sqlc"
	"sqlc/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type UseranyMatcher struct{
	user db.CreateUserParams
	password string
}

func (u UseranyMatcher) Matches(x interface{}) bool {
	arg,ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPass(u.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	u.user.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(u.user,arg)
}

func (u UseranyMatcher) String() string {
	return fmt.Sprintf("matches arg %v amd password %v",u.user, u.password)
}

func UserEQ(arg db.CreateUserParams, pass string) gomock.Matcher {
	 return UseranyMatcher{arg,pass} 
	}

func TestCreateUser(t *testing.T) {
	account, pass := randomUser(t)
	//anoynymous function & struct
	testCases := []struct {
		name          string
		Body          gin.H
		buildstubs    func(mock *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			Body: gin.H{
				"username":  account.Username,
				"password":  pass,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       account.Username,
					FullName:       account.FullName,
					Email:          account.Email,
				}
				//stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), UserEQ(arg,pass)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				Bodycheck(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			Body: gin.H{
				"username":  account.Username,
				"password":  pass,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			Body: gin.H{
				"username":  account.Username,
				"password":  pass,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			Body: gin.H{
				"username":  "test123*",
				"password":  pass,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			Body: gin.H{
				"username":  account.Username,
				"password":  pass,
				"full_name": account.FullName,
				"email":     "not_email",
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordShort",
			Body: gin.H{
				"username":  account.Username,
				"password":  "pass",
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstubs(store)

			server := NewTestServer(t,store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.Body)
			require.NoError(t, err)

			url := "/user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.Randomstring(6)
	hash,err := util.HashPassword(password)
	require.NoError(t,err)
	user = db.User{
		Username:          util.Randomowner(),
		HashedPassword:    hash,
		FullName:          util.Randomowner(),
		Email:             util.Randomemail(),
	}
	return
}

func Bodycheck(t *testing.T, body *bytes.Buffer, account db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, account.Username, gotUser.Username)
	require.Equal(t, account.FullName, gotUser.FullName)
	require.Equal(t, account.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
