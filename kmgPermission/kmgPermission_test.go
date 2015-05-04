package kmgPermission

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestKmgPermission(ot *testing.T) {
	permission := Or{
		Prefix("FastCms.Admin"),
		Prefix("vk.Admin"),
		Prefix("vk.Field"),
		And{
			Prefix("TableField"),
			Not{Prefix("TableField.StudentList._StudentAddressAndSchoolInfo")},
			Not{Prefix("TableField.StudentList._PromoterAndChannelInfo")},
		},
	}
	kmgTest.Equal(permission.IsAllow(map[string]string{"n": "FastCms.Admin"}), true)
	kmgTest.Equal(permission.IsAllow(map[string]string{"n": "FastCms.Admin1"}), true)
	kmgTest.Equal(permission.IsAllow(map[string]string{"n": "FastCms.Admi"}), false)
	kmgTest.Equal(permission.IsAllow(map[string]string{"n": "TableField.StudentList._PromoterAndChannelInfo"}), false)
	kmgTest.Equal(permission.IsAllow(map[string]string{"n": "TableField.StudentList._PromoterAndChannelInf"}), true)
}
