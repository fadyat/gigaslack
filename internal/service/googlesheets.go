package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"googlesheets-slackbot-golang/cmd/config"
	"strings"
)

var (
	ErrValueNotFound        = errors.New("value not found")
	ErrSearchColumnNotFound = errors.New("search column not found")
	ErrTakeColumnNotFound   = errors.New("take column not found")
	ErrHeadersNotFound      = errors.New("headers not found")
)

type GoogleSpreadsheets struct {
	cfg   *config.Google
	creds *google.Credentials
	svc   *sheets.Service
}

func NewGoogleSpreadsheets(cfg *config.Google, creds *google.Credentials) (*GoogleSpreadsheets, error) {
	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &GoogleSpreadsheets{
		cfg:   cfg,
		creds: creds,
		svc:   sheetsService,
	}, nil
}

func (gs *GoogleSpreadsheets) getTableHeaders(tableData [][]any) []any {
	if len(tableData) <= gs.cfg.Custom.HeaderRowIndex {
		return nil
	}

	return tableData[gs.cfg.Custom.HeaderRowIndex]
}

func getColumnIndex(tableHeaders []any, columnName string) int {
	for i, header := range tableHeaders {
		if header == columnName {
			return i
		}
	}

	return -1
}

func (gs *GoogleSpreadsheets) TakeByValue(value string, caseSensitive bool) (any, error) {
	spreadsheetData, err := gs.svc.Spreadsheets.Values.Get(gs.cfg.SpreadsheetID, gs.cfg.SpreadsheetRange).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Google Sheet data: %v", err)
	}

	tableHeaders := gs.getTableHeaders(spreadsheetData.Values)
	if len(tableHeaders) == 0 {
		return nil, ErrHeadersNotFound
	}

	searchHeaderIndex := getColumnIndex(tableHeaders, gs.cfg.Custom.SearchingValueFrom)
	if searchHeaderIndex == -1 {
		return nil, ErrSearchColumnNotFound
	}

	takeHeaderIndex := getColumnIndex(tableHeaders, gs.cfg.Custom.TakingValueFrom)
	if takeHeaderIndex == -1 {
		return nil, ErrTakeColumnNotFound
	}

	return takeByValue(spreadsheetData.Values, searchHeaderIndex, takeHeaderIndex, value, caseSensitive)
}

func takeByValue(
	values [][]any, searchHeaderIndex, takeHeaderIndex int, value string, caseSensitive bool,
) (any, error) {
	for _, row := range values {
		if !rowHaveBothColumns(row, searchHeaderIndex, takeHeaderIndex) {
			continue
		}

		if areTheSame(fmt.Sprintf("%v", row[searchHeaderIndex]), value, caseSensitive) {
			return row[takeHeaderIndex], nil
		}
	}

	return nil, ErrValueNotFound
}

func rowHaveBothColumns(row []any, searchHeaderIndex, takeHeaderIndex int) bool {
	return len(row) > searchHeaderIndex && len(row) > takeHeaderIndex
}

func areTheSame(have, want string, caseSensitive bool) bool {
	if caseSensitive {
		return have == want
	}

	return strings.EqualFold(have, want)
}
