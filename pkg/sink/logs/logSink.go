package logs

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog"
)

const (
	logsName = ""
	Suffixes = ".logs"
)

type LogSinker struct {
}

/*
*
获取event,输出当天日期.log
*/

func (l *LogSinker) ExportEvent(buffer []*corev1.Event) {
	for _, event := range buffer {
		bs, err := json.Marshal(event)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bs))
		klog.Infoln(string(bs))
	}
}
