package timeservice

import (
	"time"

	"github.com/mikiasgoitom/caching-proxy/internal/contract"
)

type TimeService struct{
	timeservice contract.ITimeService
}

func NewTimeService(timeserv contract.ITimeService) *TimeService {
	return &TimeService{
		timeservice: timeserv,
	}
}

var _ contract.ITimeService = (*TimeService)(nil)

func (ts *TimeService) Now() time.Time {
	return time.Now().UTC()
}

func (ts *TimeService) NowUnix() int64 {
	return time.Now().UTC().Unix()
}