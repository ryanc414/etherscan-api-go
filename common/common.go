//go:generate go-enum -f=$GOFILE
package common

import "time"

const (
	AccountsModule  = "account"
	ContractsModule = "contract"
	GasModule       = "gastracker"
	LogsModule      = "logs"
	ProxyModule     = "proxy"
	StatsModule     = "stats"
	TokenModule     = "token"
)

// SortingPreference is an enumeration of sorting preferences.
// ENUM(asc,desc)
type SortingPreference int32

// DateRange contains request parameters for requests that span a set of dates.
type DateRange struct {
	StartDate time.Time `etherscan:"startdate,date"`
	EndDate   time.Time `etherscan:"enddate,date"`
	Sort      SortingPreference
}

// BlockParameter is an enumeration of allowed block parameters.
// ENUM(latest,earliest,pending)
type BlockParameter int32
