package models

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/iwind/TeaGo/bootstrap"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestCreateHTTPAccessLogs(t *testing.T) {
	err := NewDBNodeInitializer().loop()
	if err != nil {
		t.Fatal(err)
	}

	accessLog := &pb.HTTPAccessLog{
		ServerId:  1,
		NodeId:    4,
		Status:    200,
		Timestamp: time.Now().Unix(),
	}
	dao := randomAccessLogDAO()
	t.Log("dao:", dao)
	err = SharedHTTPAccessLogDAO.CreateHTTPAccessLogsWithDAO(dao, []*pb.HTTPAccessLog{accessLog})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestHTTPAccessLogDAO_ListAccessLogs(t *testing.T) {
	err := NewDBNodeInitializer().loop()
	if err != nil {
		t.Fatal(err)
	}

	accessLogs, requestId, hasMore, err := SharedHTTPAccessLogDAO.ListAccessLogs("", 10, timeutil.Format("Ymd"), 0, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("requestId:", requestId, "hasMore:", hasMore)
	if len(accessLogs) == 0 {
		t.Log("no access logs yet")
		return
	}
	for _, accessLog := range accessLogs {
		t.Log(accessLog.Id, accessLog.CreatedAt, timeutil.FormatTime("H:i:s", int64(accessLog.CreatedAt)))
	}
}

func TestHTTPAccessLogDAO_ListAccessLogs_Page(t *testing.T) {
	err := NewDBNodeInitializer().loop()
	if err != nil {
		t.Fatal(err)
	}

	lastRequestId := ""

	times := 0 // 防止循环次数太多
	for {
		before := time.Now()
		accessLogs, requestId, hasMore, err := SharedHTTPAccessLogDAO.ListAccessLogs(lastRequestId, 2, timeutil.Format("Ymd"), 0, false)
		cost := time.Since(before).Seconds()
		if err != nil {
			t.Fatal(err)
		}
		lastRequestId = requestId
		if len(accessLogs) == 0 {
			break
		}
		t.Log("===")
		t.Log("requestId:", requestId[:10]+"...", "hasMore:", hasMore, "cost:", cost*1000, "ms")
		for _, accessLog := range accessLogs {
			t.Log(accessLog.Id, accessLog.CreatedAt, timeutil.FormatTime("H:i:s", int64(accessLog.CreatedAt)))
		}

		times++
		if times > 10 {
			break
		}
	}
}

func TestHTTPAccessLogDAO_ListAccessLogs_Reverse(t *testing.T) {
	err := NewDBNodeInitializer().loop()
	if err != nil {
		t.Fatal(err)
	}

	before := time.Now()
	accessLogs, requestId, hasMore, err := SharedHTTPAccessLogDAO.ListAccessLogs("16023261176446590001000000000000003500000004", 2, timeutil.Format("Ymd"), 0, true)
	cost := time.Since(before).Seconds()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("===")
	t.Log("requestId:", requestId[:19]+"...", "hasMore:", hasMore, "cost:", cost*1000, "ms")
	if len(accessLogs) > 0 {
		t.Log("accessLog:", accessLogs[0].RequestId[:19]+"...", len(accessLogs[0].RequestId))
	}
}

func TestHTTPAccessLogDAO_ListAccessLogs_Page_NotExists(t *testing.T) {
	err := NewDBNodeInitializer().loop()
	if err != nil {
		t.Fatal(err)
	}

	lastRequestId := ""

	times := 0 // 防止循环次数太多
	for {
		before := time.Now()
		accessLogs, requestId, hasMore, err := SharedHTTPAccessLogDAO.ListAccessLogs(lastRequestId, 2, timeutil.Format("Ymd", time.Now().AddDate(0, 0, 1)), 0, false)
		cost := time.Since(before).Seconds()
		if err != nil {
			t.Fatal(err)
		}
		lastRequestId = requestId
		if len(accessLogs) == 0 {
			break
		}
		t.Log("===")
		t.Log("requestId:", requestId[:10]+"...", "hasMore:", hasMore, "cost:", cost*1000, "ms")
		for _, accessLog := range accessLogs {
			t.Log(accessLog.Id, accessLog.CreatedAt, timeutil.FormatTime("H:i:s", int64(accessLog.CreatedAt)))
		}

		times++
		if times > 10 {
			break
		}
	}
}
