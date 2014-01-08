// package test contains utilities to test topo.Server
// implementations. If you are testing your implementation, you will
// want to call CheckAll in your test method. For an example, look at
// the tests in github.com/youtube/vitess/go/vt/zktopo.
package test

import (
	"testing"

	"github.com/youtube/vitess/go/vt/topo"
)

func CheckKeyspace(t *testing.T, ts topo.Server) {
	keyspaces, err := ts.GetKeyspaces()
	if err != nil {
		t.Errorf("GetKeyspaces(empty): %v", err)
	}
	if len(keyspaces) != 0 {
		t.Errorf("len(GetKeyspaces()) != 0: %v", keyspaces)
	}

	if err := ts.CreateKeyspace("test_keyspace", &topo.Keyspace{}); err != nil {
		t.Errorf("CreateKeyspace: %v", err)
	}
	if err := ts.CreateKeyspace("test_keyspace", &topo.Keyspace{}); err != topo.ErrNodeExists {
		t.Errorf("CreateKeyspace(again) is not ErrNodeExists: %v", err)
	}

	keyspaces, err = ts.GetKeyspaces()
	if err != nil {
		t.Errorf("GetKeyspaces: %v", err)
	}
	if len(keyspaces) != 1 || keyspaces[0] != "test_keyspace" {
		t.Errorf("GetKeyspaces: want %v, got %v", []string{"test_keyspace"}, keyspaces)
	}

	if err := ts.CreateKeyspace("test_keyspace2", &topo.Keyspace{ShardingColumnName: "user_id", ShardingColumnType: topo.SCT_UINT64}); err != nil {
		t.Errorf("CreateKeyspace: %v", err)
	}
	keyspaces, err = ts.GetKeyspaces()
	if err != nil {
		t.Errorf("GetKeyspaces: %v", err)
	}
	if len(keyspaces) != 2 || keyspaces[0] != "test_keyspace" || keyspaces[1] != "test_keyspace2" {
		t.Errorf("GetKeyspaces: want %v, got %v", []string{"test_keyspace", "test_keyspace2"}, keyspaces)
	}

	ki, err := ts.GetKeyspace("test_keyspace2")
	if err != nil {
		t.Fatalf("GetKeyspace: %v", err)
	}
	if ki.ShardingColumnName != "user_id" || ki.ShardingColumnType != topo.SCT_UINT64 {
		t.Errorf("GetKeyspace: want user_id/uint64, got %v/%v", ki.ShardingColumnName, ki.ShardingColumnType)
	}

	ki.ShardingColumnName = "other_id"
	ki.ShardingColumnType = topo.SCT_BYTES
	err = ts.UpdateKeyspace(ki)
	if err != nil {
		t.Fatalf("UpdateKeyspace: %v", err)
	}
	ki, err = ts.GetKeyspace("test_keyspace2")
	if err != nil {
		t.Fatalf("GetKeyspace: %v", err)
	}
	if ki.ShardingColumnName != "other_id" || ki.ShardingColumnType != topo.SCT_BYTES {
		t.Errorf("GetKeyspace: want other_id/bytes, got %v/%v", ki.ShardingColumnName, ki.ShardingColumnType)
	}
}
