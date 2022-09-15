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
	"sqlc/dummy"
	"testing"

	"github.com/golang/mock/gomock"
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
		t.Run(tc.name,func(t *testing.T) {
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
		tc.checkResponse(t,recorder)
		})
		
	}
}

func randomacc() db.Account {
	return db.Account{
		ID:       dummy.Randomint(1, 1000),
		Owner:    dummy.Randomowner(),
		Balance:  dummy.Randommoney(),
		Currency: dummy.Randomcurrency(),
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
