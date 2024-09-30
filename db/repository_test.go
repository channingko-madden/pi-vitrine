package db_test

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"testing"
	"time"
)

func TestReportingPeriodWhereFilter(t *testing.T) {
	r := db.PostgresDeviceRepository{}

	start := time.Time{}
	end := time.Time{}

	filter := r.ReportingPeriodWhereFilter(start, end)

	if len(filter) != 0 {
		t.Fatal("The returned filter should be empty if start & end are zero values")
	}

	start = time.Date(2024, 4, 20, 16, 20, 0, 0, time.UTC)

	filter = r.ReportingPeriodWhereFilter(start, end)

	expected := "created_at >= 2024-04-20 16:20:00"

	if filter != expected {
		t.Fatalf("The returned start only filter %s does not match the expected filter", filter)
	}

	start = time.Time{}
	end = time.Date(2024, 4, 20, 16, 20, 0, 0, time.UTC)

	filter = r.ReportingPeriodWhereFilter(start, end)

	expected = "created_at <= 2024-04-20 16:20:00"

	if filter != expected {
		t.Fatalf("The returned end only filter %s does not match the expected filter", filter)
	}

	start = time.Date(2023, 4, 20, 16, 20, 0, 0, time.UTC)
	end = time.Date(2024, 4, 20, 16, 20, 0, 0, time.UTC)

	filter = r.ReportingPeriodWhereFilter(start, end)

	expected = "created_at >= 2023-04-20 16:20:00 and created_at <= 2024-04-20 16:20:00"

	if filter != expected {
		t.Fatalf("The returned start + end filter %s does not match the expected filter", filter)
	}
}
