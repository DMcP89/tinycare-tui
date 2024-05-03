package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_HumanizeDuration(t *testing.T) {
	lessThan24hours := time.Duration(16) * time.Hour
	greaterThan24hours := time.Duration(48) * time.Hour

	expectedResult := fmt.Sprintf("[green]%d h(s) ago[white]", int(lessThan24hours.Hours()))
	result := HumanizeDuration(lessThan24hours)
	if expectedResult != result {
		t.Errorf("Error testing HumanizeDuration: Expected: %s, Got: %s", expectedResult, result)
	}

	expectedResult = fmt.Sprintf("[green]%d d(s) ago[white]", int(greaterThan24hours.Hours()/24))
	result = HumanizeDuration(greaterThan24hours)
	if expectedResult != result {
		t.Errorf("Error testing HumanizeDuration: Expected: %s, Got: %s", expectedResult, result)
	}
}
