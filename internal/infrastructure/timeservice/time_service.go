package timeservice

import (
	"time"

	"github.com/mikiasgoitom/RevProx/internal/contract"
)

type TimeService struct {
}

func NewTimeService() contract.ITimeService {
	return &TimeService{}
}

var _ contract.ITimeService = (*TimeService)(nil)

func (ts *TimeService) Now() time.Time {
	return time.Now().UTC()
}

func (ts *TimeService) NowUnix() int64 {
	return time.Now().UTC().Unix()
}
