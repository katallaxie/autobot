package modules

import (
	"fmt"
)

var K8sTag = "k8s"

func K8sReplacer(locale, entry, response, _ string) (string, string) {
	return K8sTag, fmt.Sprintf(response, "tp.infra.cluster.ionos.com-paasis-dev")
}
