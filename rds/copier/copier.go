package copier

import "github.com/jinzhu/copier"

func UpgradeSchema(in1 interface{}, in2 interface{}) (err error) {
	//for redis hash,zset,set,list, upgrade the data schema
	return copier.Copy(in1, in1)
}
