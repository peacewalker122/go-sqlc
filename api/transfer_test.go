package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	mockdb "sqlc/db/mock"
	db "sqlc/db/sqlc"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	amount := int64(10)

	user1,_ := randomUser(t)
	user2,_ := randomUser(t)
	user3,_ := randomUser(t)

	account1 := randomacc(user1.Username)
	account2 := randomacc(user2.Username)
	account3 := randomacc(user3.Username)

	account1.Currency = "USD"
	account2.Currency = "USD"
	account3.Currency = "EUR"
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mock *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_id":  account1.ID,
				"to_id":    account2.ID,
				"amount":   amount,
				"currency": "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferctxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferCtx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_id":  account1.ID,
				"to_id":    account2.ID,
				"amount":   amount,
				"currency": "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
		name: "ToAccountNotFound",
		body: gin.H{
			"from_id":  account1.ID,
			"to_id":    account2.ID,
			"amount":   amount,
			"currency": "USD",
		},
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
			store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{},sql.ErrNoRows)
			store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(0)
		},
		checkResponse: func(recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
		},
			{
			name: "From_WrongCurrency",
			body: gin.H{
				"from_id":  account3.ID,
				"to_id":    account2.ID,
				"amount":   amount,
				"currency": "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2,nil)
				store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
				name: "To_WrongCurrency",
				body: gin.H{
					"from_id":  account1.ID,
					"to_id":    account3.ID,
					"amount":   amount,
					"currency": "USD",
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1,nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
					store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_id":  account1.ID,
				"to_id":    account2.ID,
				"amount":   amount,
				"currency": "YEN",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountErr",
			body: gin.H{
				"from_id":  account1.ID,
				"to_id":    account2.ID,
				"amount":   amount,
				"currency": "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{},sql.ErrConnDone)
				store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferError",
			body: gin.H{
				"from_id":  account1.ID,
				"to_id":    account2.ID,
				"amount":   amount,
				"currency": "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferCtx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{},sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			recorder := httptest.NewRecorder()
			tc.buildStubs(store)

			server := NewTestServer(t,store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/transfers"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}

}
