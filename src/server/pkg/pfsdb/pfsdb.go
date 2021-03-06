// Package pfsdb contains the database schema that PFS uses.
package pfsdb

import (
	"fmt"
	"path"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/pachyderm/pachyderm/src/client/pfs"
	"github.com/pachyderm/pachyderm/src/client/pkg/uuid"
	col "github.com/pachyderm/pachyderm/src/server/pkg/collection"
)

const (
	reposPrefix          = "/repos"
	repoRefCountsPrefix  = "/repoRefCounts"
	putFileRecordsPrefix = "/putFileRecords"
	commitsPrefix        = "/commits"
	branchesPrefix       = "/branches"
	openCommitsPrefix    = "/openCommits"
)

var (
	// ProvenanceIndex is a secondary index on provenance
	ProvenanceIndex = col.Index{"Provenance", true}
)

// Repos returns a collection of repos
func Repos(etcdClient *etcd.Client, etcdPrefix string) col.Collection {
	return col.NewCollection(
		etcdClient,
		path.Join(etcdPrefix, reposPrefix),
		[]col.Index{ProvenanceIndex},
		&pfs.RepoInfo{},
		nil,
	)
}

// RepoRefCounts returns a collection of repo ref counts
func RepoRefCounts(etcdClient *etcd.Client, etcdPrefix string) col.Collection {
	return col.NewCollection(
		etcdClient,
		path.Join(etcdPrefix, repoRefCountsPrefix),
		nil,
		nil,
		nil,
	)
}

// PutFileRecords returns a collection of putFileRecords
func PutFileRecords(etcdClient *etcd.Client, etcdPrefix string) col.Collection {
	return col.NewCollection(
		etcdClient,
		path.Join(etcdPrefix, putFileRecordsPrefix),
		nil,
		&pfs.PutFileRecords{},
		nil,
	)
}

// Commits returns a collection of commits
func Commits(etcdClient *etcd.Client, etcdPrefix string, repo string) col.Collection {
	return col.NewCollection(
		etcdClient,
		path.Join(etcdPrefix, commitsPrefix, repo),
		[]col.Index{ProvenanceIndex},
		&pfs.CommitInfo{},
		nil,
	)
}

// Branches returns a collection of branches
func Branches(etcdClient *etcd.Client, etcdPrefix string, repo string) col.Collection {
	return col.NewCollection(
		etcdClient,
		path.Join(etcdPrefix, branchesPrefix, repo),
		nil,
		&pfs.BranchInfo{},
		func(key string) error {
			if uuid.IsUUIDWithoutDashes(key) {
				return fmt.Errorf("branch name cannot be a UUID V4")
			}
			return nil
		},
	)
}

// OpenCommits returns a collection of open commits
func OpenCommits(etcdClient *etcd.Client, etcdPrefix string) col.Collection {
	return col.NewCollection(
		etcdClient,
		path.Join(etcdPrefix, openCommitsPrefix),
		nil,
		&pfs.Commit{},
		nil,
	)
}
