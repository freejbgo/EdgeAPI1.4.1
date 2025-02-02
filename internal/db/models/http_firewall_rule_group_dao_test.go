package models

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/dbs"
)

func TestHTTPFirewallRuleGroupDAO_FindRuleGroupIdWithRuleSetId(t *testing.T) {
	dbs.NotifyReady()

	var tx *dbs.Tx
	groupId, err := SharedHTTPFirewallRuleGroupDAO.FindRuleGroupIdWithRuleSetId(tx, 22)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("groupId:", groupId)
}