// +build !codeanalysis

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/TencentBlueKing/iam-go-sdk"
	"github.com/TencentBlueKing/iam-go-sdk/logger"
)

func main() {
	// create a logger
	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	// do set logger
	logger.SetLogger(log)

	req := iam.NewRequest(
		"bk_paas",
		iam.NewSubject("user", "admin"),
		iam.NewAction("access_developer_center"),
		[]iam.ResourceNode{},
	)

	i := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")

	allowed, err := i.IsAllowed(req)
	fmt.Println("isAllowed:", allowed, err)

	// check 3 times but only call iam backend once
	allowed, err = i.IsAllowedWithCache(req, 10*time.Second)
	allowed, err = i.IsAllowedWithCache(req, 10*time.Second)
	i2 := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")
	allowed, err = i2.IsAllowedWithCache(req, 10*time.Second)
	fmt.Println("isAllowedWithCache:", allowed, err)

	multiReq := iam.NewMultiActionRequest(
		"bk_sops",
		iam.NewSubject("user", "admin"),
		[]iam.Action{
			iam.NewAction("task_delete"),
			iam.NewAction("task_edit"),
			iam.NewAction("task_view"),
		},
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "1", map[string]interface{}{"iam_resource_owner": "admin"}),
		},
	)
	i3 := iam.NewIAM("bk_sops", "bk_sops", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")
	result, err := i3.ResourceMultiActionsAllowed(multiReq)
	fmt.Println("ResourceMultiActionsAllowed: ", result, err)

	multiReq.Resources = iam.Resources{}
	resourcesList := []iam.Resources{
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "1", map[string]interface{}{"iam_resource_owner": "admin"}),
		},
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "2", map[string]interface{}{"iam_resource_owner": "admin2"}),
		},
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "3", map[string]interface{}{"iam_resource_owner": "admin3"}),
		},
	}
	results, err := i3.BatchResourceMultiActionsAllowed(multiReq, resourcesList)
	fmt.Println("BatchResourceMultiActionsAllowed: ", results, err)

	actions := []iam.ApplicationAction{}
	application := iam.NewApplication("bk_paas", actions)

	url, err := i.GetApplyURL(application, "", "admin")
	fmt.Println("GetApplyURL:", url, err)

	err = i.IsBasicAuthAllowed("bk_iam", "3b223guhmzwlq7oto417b4g41rqnboip")
	fmt.Println("IsBasicAuthAllowed:", err)

	token, err := i.GetToken()
	fmt.Println("GetToken:", token, err)
}
