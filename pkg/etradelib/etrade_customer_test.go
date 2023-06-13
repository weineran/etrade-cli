package etradelib

import (
	"errors"
	"github.com/jerryryle/etrade-cli/pkg/etradelib/client"
	"github.com/jerryryle/etrade-cli/pkg/etradelib/responses"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ETradeCustomerTestSuite struct {
	suite.Suite
	clientFake   client.ETradeClientFake
	testCustomer ETradeCustomer
}

func (s *ETradeCustomerTestSuite) SetupTest() {
	// Create an empty client fake. Tests will need to set required members before
	// making client calls
	s.clientFake = client.ETradeClientFake{}

	// Create a test customer using the fake client
	s.testCustomer = CreateETradeCustomer(&s.clientFake, "TestCustomerName")
}

func (s *ETradeCustomerTestSuite) TestETradeCustomer_GetAllAccounts() {
	// Initialize a client response with just enough data to verify the objects are created correctly
	testAccounts := []responses.AccountListAccount{
		{AccountId: "TestAccountId1"},
		{AccountId: "TestAccountId2"},
	}

	// Configure the client fake to return the test data
	s.clientFake.ListAccountsFn = func() (*responses.AccountListResponse, error) {
		return &responses.AccountListResponse{
			Accounts: testAccounts,
		}, nil
	}

	// Get all accounts and ensure all account objects are returned.
	accounts, err := s.testCustomer.GetAllAccounts()
	s.Assert().Nil(err)
	for idx, account := range accounts {
		s.Assert().Equal(testAccounts[idx].AccountId, account.GetAccountInfo().AccountId)
	}

	// Configure the client fake to return an error
	s.clientFake.ListAccountsFn = func() (*responses.AccountListResponse, error) {
		return nil, errors.New("test error")
	}
	// Get all accounts and ensure the client failure error is propagated
	accounts, err = s.testCustomer.GetAllAccounts()
	s.Assert().Nil(accounts)
	s.Assert().EqualError(err, "test error")
}

func (s *ETradeCustomerTestSuite) TestETradeCustomer_GetAccountById() {
	// Initialize a client response with just enough data to verify the objects are created correctly
	testAccounts := []responses.AccountListAccount{
		{AccountId: "TestAccountId1"},
		{AccountId: "TestAccountId2"},
	}

	// Configure the client fake to return the test data
	s.clientFake.ListAccountsFn = func() (*responses.AccountListResponse, error) {
		return &responses.AccountListResponse{
			Accounts: testAccounts,
		}, nil
	}

	// Get an account by ID and ensure the correct account object is returned.
	account, err := s.testCustomer.GetAccountById("TestAccountId2")
	s.Assert().Nil(err)
	s.Assert().Equal("TestAccountId2", account.GetAccountInfo().AccountId)

	// Get an account by non-existent ID and ensure an error is returned.
	account, err = s.testCustomer.GetAccountById("TestAccountNonExistent")
	s.Assert().Nil(account)
	s.Assert().Error(err)

	// Configure the client fake to return an error
	s.clientFake.ListAccountsFn = func() (*responses.AccountListResponse, error) {
		return nil, errors.New("test error")
	}
	// Get an account by ID and ensure the client failure error is propagated
	account, err = s.testCustomer.GetAccountById("TestAccountId2")
	s.Assert().Nil(account)
	s.Assert().EqualError(err, "test error")
}

func (s *ETradeCustomerTestSuite) TestETradeCustomer_GetAllAlerts() {
	// Initialize a client response with just enough data to verify the objects are created correctly
	testAlerts := []responses.AlertsAlert{
		{Id: 1},
		{Id: 2},
	}

	// Configure the client fake to return the test data
	s.clientFake.ListAlertsFn = func() (*responses.AlertsResponse, error) {
		return &responses.AlertsResponse{
			Alerts: testAlerts,
		}, nil
	}

	// Get all alerts and ensure all alert objects are returned.
	alerts, err := s.testCustomer.GetAllAlerts()
	s.Assert().Nil(err)
	for idx, alert := range alerts {
		s.Assert().Equal(testAlerts[idx].Id, alert.GetAlertInfo().Id)
	}

	// Configure the client fake to return an error
	s.clientFake.ListAlertsFn = func() (*responses.AlertsResponse, error) {
		return nil, errors.New("test error")
	}
	// Get all alerts and ensure the client failure error is propagated
	alerts, err = s.testCustomer.GetAllAlerts()
	s.Assert().Nil(alerts)
	s.Assert().EqualError(err, "test error")
}

func (s *ETradeCustomerTestSuite) TestETradeCustomer_GetAlertById() {
	// Initialize a client response with just enough data to verify the objects are created correctly
	testAlerts := []responses.AlertsAlert{
		{Id: 1},
		{Id: 2},
	}

	// Configure the client fake to return the test data
	s.clientFake.ListAlertsFn = func() (*responses.AlertsResponse, error) {
		return &responses.AlertsResponse{
			Alerts: testAlerts,
		}, nil
	}

	// Get an account by ID and ensure the correct account object is returned.
	alert, err := s.testCustomer.GetAlertById(2)
	s.Assert().Nil(err)
	s.Assert().Equal(int64(2), alert.GetAlertInfo().Id)

	// Get an account by non-existent ID and ensure an error is returned.
	alert, err = s.testCustomer.GetAlertById(999)
	s.Assert().Nil(alert)
	s.Assert().Error(err)

	// Configure the client fake to return an error
	s.clientFake.ListAlertsFn = func() (*responses.AlertsResponse, error) {
		return nil, errors.New("test error")
	}
	// Get an account by ID and ensure the client failure error is propagated
	alert, err = s.testCustomer.GetAlertById(2)
	s.Assert().Nil(alert)
	s.Assert().EqualError(err, "test error")
}

func TestETradeCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(ETradeCustomerTestSuite))
}
