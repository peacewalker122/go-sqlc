package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	mockdb "sqlc/db/mock"
	db "sqlc/db/sqlc"
	"sqlc/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetByID(t *testing.T) {

	account := randomacc()

	//anoynymous function & struct
	testCases := []struct {
		name          string
		ID            int64
		buildstubs    func(mock *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   account.ID,
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(r *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				BodyTest(t, recorder.Body, account)
			},
		},
		{
			name: "Not Found",
			ID:   account.ID,
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(r *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   account.ID,
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(r *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   0,
			buildstubs: func(mock *mockdb.MockStore) {
				//stubs
				mock.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(r *testing.T, recorder *httptest.ResponseRecorder) {
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

			server := Newserver(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%v", tc.ID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestCreateAccount(t *testing.T) {
	account := randomacc()
	//user, _ := randomUser(t)
	//anoynymous function & struct
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			// u need to put inside the "body" with the same argument in createAccountParam. "in this version only owner and currency"
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			// u need to put inside the "body" with the same argument in createAccountParam. "in this version only owner and currency"
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			// u need to put inside the "body" with the same argument in createAccountParam. "in this version only owner and currency"
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "NonUsername",
			// u need to put inside the "body" with the same argument in createAccountParam. "in this version only owner and currency"
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: "23503"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidOwner",
			body: gin.H{
				"owner":    "",
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			// u need to put inside the "body" with the same argument in createAccountParam. "in this version only owner and currency"
			body: gin.H{
				"owner":    account.Owner,
				"currency": "YEN",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1)
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
			tc.buildStubs(store)

			server := Newserver(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomacc() db.Account {
	return db.Account{
		ID:       util.Randomint(1, 1000),
		Owner:    util.Randomowner(),
		Balance:  util.Randommoney(),
		Currency: util.Randomcurrency(),
	}
}

func BodyTest(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var Account db.Account
	err = json.Unmarshal(data, &Account)
	require.NoError(t, err)
	require.Equal(t, account, Account)
}

func Randomuser() (db.User, db.Account) {
	x := db.User{
		Username:       util.Randomowner(),
		HashedPassword: util.Randomstring(6),
		FullName:       util.Randomowner(),
		Email:          util.Randomemail(),
	}
	y := db.Account{
		ID:       util.Randomint(1, 100),
		Owner:    x.Username,
		Balance:  util.Randomint(1, 1000),
		Currency: util.Randomcurrency(),
	}

	return x, y
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
