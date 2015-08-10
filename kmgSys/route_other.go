// +build !darwin

package kmgSys

type Route struct {
}

func GetRouteTable() (routeList []Route, err error) {
	panic("get route table only implement darwin")
}
