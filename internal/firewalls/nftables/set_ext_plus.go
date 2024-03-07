// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build linux && plus

package nftables

import (
	"time"
)

func (this *Set) initElements() {
	elements, err := this.conn.Raw().GetSetElements(this.rawSet)
	var unixMilli = time.Now().UnixMilli()
	if err == nil {
		this.expiration = NewExpiration()
		for _, element := range elements {
			if element.Expires == 0 {
				this.expiration.AddUnsafe(element.Key, time.Time{})
			} else {
				this.expiration.AddUnsafe(element.Key, time.UnixMilli(unixMilli+element.Expires.Milliseconds()))
			}
		}
	}
}
