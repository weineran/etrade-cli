package client

import (
	"github.com/jerryryle/etrade-cli/pkg/etradelib/client/constants"
	"github.com/jerryryle/etrade-cli/pkg/etradelib/etradelibtest"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func createTestClient(t *testing.T, responseData string, expectMethod string, expectUrl string) ETradeClient {
	httpClient := NewHttpClientFake(
		func(req *http.Request) *http.Response {
			assert.Equal(t, expectMethod, req.Method)
			assert.Equal(t, expectUrl, req.URL.String())
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseData)),
			}
		},
	)
	return CreateETradeClient(GetEndpointUrls(true), httpClient, etradelibtest.CreateNullLogger())
}

func TestETradeClient(t *testing.T) {
	type testFn func(client ETradeClient) ([]byte, error)

	// testResponseData is bogus JSON that's only used to ensure the client returns the exact response from the server
	const testResponseData = `{"testResponse": true}`

	tests := []struct {
		name           string
		testFn         testFn
		expectMethod   string
		expectUrl      string
		expectResponse []byte
		expectErr      bool
	}{
		{
			name: "List Accounts",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListAccounts()
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/list",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Account Balances",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetAccountBalances("1234", true)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/balance?instType=BROKERAGE&realTimeNAV=true",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Account Balances Fails Without Account ID Key",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetAccountBalances("", true)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "List Transactions With All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListTransactions(
					"1234",
					etradelibtest.CreateTime(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					etradelibtest.CreateTime(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
					constants.SortOrderAsc, "5", 6,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/transactions?count=6&endDate=01022023&marker=5&sortOrder=ASC&startDate=01012023",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Transactions Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListTransactions("1234", nil, nil, constants.SortOrderNil, "", -1)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/transactions",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Transactions Fails Without Account ID Key",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListTransactions("", nil, nil, constants.SortOrderNil, "", -1)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "List Transaction Details",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListTransactionDetails("1234", "5678")
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/transactions/5678",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Transaction Details Fails Without Account ID Key",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListTransactionDetails("", "5678")
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "List Transaction Details Fails Without Transaction ID",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListTransactionDetails("1234", "")
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "View Portfolio With All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ViewPortfolio(
					"1234", 5, constants.PortfolioSortBySymbol, constants.SortOrderAsc, "6",
					constants.MarketSessionRegular, true, true, constants.PortfolioViewComplete,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/portfolio?count=5&lotsRequired=true&marketSession=REGULAR&pageNumber=6&sortBy=SYMBOL&sortOrder=ASC&totalsRequired=true&view=COMPLETE",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "View Portfolio Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ViewPortfolio(
					"1234", -1, constants.PortfolioSortByNil, constants.SortOrderNil, "",
					constants.MarketSessionNil, true, true, constants.PortfolioViewNil,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/portfolio?lotsRequired=true&totalsRequired=true",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "View Portfolio Fails Without Account ID Key",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ViewPortfolio(
					"", -1, constants.PortfolioSortByNil, constants.SortOrderNil, "",
					constants.MarketSessionNil, true, true, constants.PortfolioViewNil,
				)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "View Portfolio Fails If Count Is Too Big",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ViewPortfolio(
					"1234", constants.PortfolioMaxCount+1, constants.PortfolioSortByNil, constants.SortOrderNil, "",
					constants.MarketSessionNil, true, true, constants.PortfolioViewNil,
				)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "List Alerts With All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListAlerts(
					1, constants.AlertCategoryAccount, constants.AlertStatusUnread, constants.SortOrderAsc, "FOO",
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/user/alerts?category=ACCOUNT&count=1&direction=ASC&search=FOO&status=UNREAD",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Alerts Fails With Count Too Big",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListAlerts(
					301, constants.AlertCategoryAccount, constants.AlertStatusUnread, constants.SortOrderAsc, "FOO",
				)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "List Alerts Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListAlerts(
					-1, constants.AlertCategoryNil, constants.AlertStatusNil, constants.SortOrderNil, "",
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/user/alerts",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Alert Details",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListAlertDetails("1234", true)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/user/alerts/1234?htmlTags=true",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Delete Alerts",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.DeleteAlerts([]string{"1234", "5678"})
			},
			expectMethod:   "DELETE",
			expectUrl:      "https://api.etrade.com/v1/user/alerts/1234,5678",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Quotes With All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetQuotes([]string{"GOOG"}, constants.QuoteDetailAll, true, false)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/quote/GOOG?detailFlag=ALL&requireEarningsDate=true&skipMiniOptionsCheck=false",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Quotes Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetQuotes([]string{"GOOG"}, constants.QuoteDetailNil, true, false)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/quote/GOOG?requireEarningsDate=true&skipMiniOptionsCheck=false",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Quotes Fails Without Symbols",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetQuotes([]string{}, constants.QuoteDetailNil, true, false)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "Get Quotes Overrides When More Than 25 Symbols",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetQuotes(
					[]string{
						"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17",
						"18", "19", "20", "21", "22", "23", "24", "25", "26",
					}, constants.QuoteDetailAll, true, false,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/quote/1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26?detailFlag=ALL&overrideSymbolCount=true&requireEarningsDate=true&skipMiniOptionsCheck=false",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Quotes Fails With More Than 50 Symbols",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetQuotes(
					[]string{
						"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17",
						"18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33",
						"34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
						"50", "51",
					}, constants.QuoteDetailAll, true, false,
				)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "Lookup Product",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.LookupProduct("A")
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/lookup/A",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Lookup Product Fails With Empty Search String",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.LookupProduct("")
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "Get Option Chains With All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetOptionChains(
					"GOOG",
					1, 2, 3,
					4, 5,
					true, true,
					constants.OptionCategoryAll, constants.OptionChainTypeCall, constants.OptionPriceTypeAll,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/optionchains?chainType=CALL&expiryDay=3&expiryMonth=2&expiryYear=1&includeWeekly=true&noOfStrikes=5&optionCategory=ALL&priceType=ALL&skipAdjusted=true&strikePriceNear=4&symbol=GOOG",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Option Chains Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetOptionChains(
					"GOOG",
					-1, -1, -1,
					-1, -1,
					true, true,
					constants.OptionCategoryNil, constants.OptionChainTypeNil, constants.OptionPriceTypeNil,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/optionchains?includeWeekly=true&skipAdjusted=true&symbol=GOOG",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Option Chains Fails Without Symbol",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetOptionChains(
					"",
					-1, -1, -1,
					-1, -1,
					true, true,
					constants.OptionCategoryNil, constants.OptionChainTypeNil, constants.OptionPriceTypeNil,
				)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "Get Option Expire Date With All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetOptionExpireDates("GOOG", constants.OptionExpiryTypeAll)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/optionexpiredate?expiryType=ALL&symbol=GOOG",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Option Expire Date Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetOptionExpireDates("GOOG", constants.OptionExpiryTypeNil)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/market/optionexpiredate?symbol=GOOG",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "Get Option Expire Date Fails Without Symbol",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.GetOptionExpireDates("", constants.OptionExpiryTypeNil)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
		{
			name: "List Orders With All Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListOrders(
					"1234", "TestMarker", 5, constants.OrderStatusOpen,
					etradelibtest.CreateTime(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					etradelibtest.CreateTime(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
					[]string{"A", "B"},
					constants.OrderSecurityTypeEquity, constants.OrderTransactionTypeBuy,
					constants.MarketSessionRegular,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/orders?count=5&fromDate=01012023&marker=TestMarker&marketSession=REGULAR&securityType=EQ&status=OPEN&symbol=A%2CB&toDate=01022023&transactionType=BUY",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Orders Can Omit All Optional Arguments",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListOrders(
					"1234", "", -1, constants.OrderStatusNil, nil, nil, nil, constants.OrderSecurityTypeNil,
					constants.OrderTransactionTypeNil, constants.MarketSessionNil,
				)
			},
			expectMethod:   "GET",
			expectUrl:      "https://api.etrade.com/v1/accounts/1234/orders",
			expectResponse: []byte(testResponseData),
			expectErr:      false,
		},
		{
			name: "List Orders Fails Without Account ID Key",
			testFn: func(client ETradeClient) ([]byte, error) {
				return client.ListOrders(
					"", "", -1, constants.OrderStatusNil, nil, nil, nil, constants.OrderSecurityTypeNil,
					constants.OrderTransactionTypeNil, constants.MarketSessionNil,
				)
			},
			expectMethod:   "",
			expectUrl:      "",
			expectResponse: nil,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				client := createTestClient(t, testResponseData, tt.expectMethod, tt.expectUrl)
				// Call the Method Under Test
				response, err := tt.testFn(client)
				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.Nil(t, err)
				}
				assert.Equal(t, tt.expectResponse, response)
			},
		)
	}
}
