package kmgSys

import "github.com/kardianos/osext"

func GetCurrentExecutePath() (string, error) {
	return osext.Executable()
}
