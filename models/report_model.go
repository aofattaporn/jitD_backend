package models

import "time"

type Report struct {
	UserID     string    `json:"userID"`
	ReportName string    `json:"reportName"`
	ReportType string    `json:"reportType"`
	Date       time.Time `json:"date,omitempty"`
}

type ReportResponse struct {
	ReportID   string    `json:"reportID"`
	UserID     string    `json:"userID"`
	ReportName string    `json:"reportName"`
	ReportType string    `json:"reportType"`
	Date       time.Time `json:"date,omitempty"`
}
