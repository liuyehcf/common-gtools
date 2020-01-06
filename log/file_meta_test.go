package log

import (
	"fmt"
	"github.com/liuyehcf/common-gtools/utils"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const (
	secondsPerDay  = 86400
	secondsPerHour = 3600
	dayFormat      = "2006-01-02"
	dayHourFormat  = "2006-01-02 15:04:05"
)

func TestDayDiff(t *testing.T) {
	fromDate, err := time.Parse(dayFormat, "2006-01-02")
	utils.AssertNil(err, "test")

	toDate, err := time.Parse(dayFormat, "2020-01-02")
	utils.AssertNil(err, "test")

	metas := make(fileMetaSlice, 0)
	dayStrings := make([]string, 0)
	for t := fromDate.Unix(); t < toDate.Unix(); t += secondsPerDay {
		unix := time.Unix(t, 0)
		dayString := unix.Format(dayFormat)
		dayStrings = append(dayStrings, dayString)
		metas = append(metas, newFileMeta("/test", dayString, "01", "1"))
	}
	utils.AssertTrue(len(metas) > 3650, "test")
	fmt.Printf("day num=%d\n", len(metas))

	// shuffle
	shuffle(metas)

	// sort
	sort.Sort(metas)

	for i := 0; i < len(metas); i += 1 {
		utils.AssertTrue(metas[i].day == dayStrings[i], "test")
	}
}

func TestHourDiff(t *testing.T) {
	fromDate, err := time.Parse(dayFormat, "2019-12-01")
	utils.AssertNil(err, "test")

	toDate, err := time.Parse(dayFormat, "2020-01-02")
	utils.AssertNil(err, "test")

	metas := make(fileMetaSlice, 0)
	dayStrings := make([]string, 0)
	hourStrings := make([]string, 0)
	for t := fromDate.Unix(); t < toDate.Unix(); t += secondsPerHour {
		unix := time.Unix(t, 0)
		dayHourTimeString := unix.Format(dayHourFormat)
		dayHourTime, err := time.Parse(dayHourFormat, dayHourTimeString)
		utils.AssertNil(err, "test")

		dayString := fmt.Sprintf("%d-%02d-%02d", dayHourTime.Year(), dayHourTime.Month(), dayHourTime.Day())
		dayStrings = append(dayStrings, dayString)

		hourString := fmt.Sprintf("%02d", dayHourTime.Hour())
		hourStrings = append(hourStrings, hourString)

		metas = append(metas, newFileMeta("/test", dayString, hourString, "1"))
	}
	utils.AssertTrue(len(metas) > 240, "test")
	fmt.Printf("hour num=%d\n", len(metas))

	// shuffle
	shuffle(metas)

	// sort
	sort.Sort(metas)

	for i := 0; i < len(metas); i += 1 {
		utils.AssertTrue(metas[i].day == dayStrings[i], "test")
		utils.AssertTrue(metas[i].hour == hourStrings[i], "test")
	}
}

func TestIndexDiff(t *testing.T) {
	fromIndex := 1
	toIndex := 10000

	metas := make(fileMetaSlice, 0)
	for t := fromIndex; t <= toIndex; t += 1 {
		metas = append(metas, newFileMeta("/test", "2020-01-02", "01", fmt.Sprintf("%d", t)))
	}
	fmt.Printf("index num=%d\n", len(metas))

	// shuffle
	shuffle(metas)

	// sort
	sort.Sort(metas)

	for i := 0; i < len(metas); i += 1 {
		utils.AssertTrue(metas[i].indexValue == toIndex-i, "test")
	}
}

func shuffle(slice fileMetaSlice) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
