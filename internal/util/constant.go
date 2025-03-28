// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"strings"

	"github.com/milvus-io/milvus/internal/proto/commonpb"
)

// Meta Prefix consts
const (
	SegmentMetaPrefix    = "queryCoord-segmentMeta"
	ChangeInfoMetaPrefix = "queryCoord-sealedSegmentChangeInfo"
	HandoffSegmentPrefix = "querycoord-handoff"
	HeaderAuthorize      = "authorization"
	// HeaderSourceID identify requests from Milvus members and client requests
	HeaderSourceID = "sourceId"
	// MemberCredID id for Milvus members (data/index/query node/coord component)
	MemberCredID        = "@@milvus-member@@"
	CredentialSeperator = ":"
	UserRoot            = "root"
	DefaultRootPassword = "Milvus"
	DefaultTenant       = ""
	RoleAdmin           = "admin"
	RolePublic          = "public"

	PrivilegeWord = "Privilege"
)

// StringSet convert array to map for conveniently check if the array contains an element
func StringSet(strings []string) map[string]struct{} {
	stringsMap := make(map[string]struct{})
	for _, str := range strings {
		stringsMap[str] = struct{}{}
	}
	return stringsMap
}

func StringList(stringMap map[string]struct{}) []string {
	strs := make([]string, 0, len(stringMap))
	for k := range stringMap {
		strs = append(strs, k)
	}
	return strs
}

// GetObjectPrivileges get the mapping between object types and privileges. This kind of data is constant and doesn't support to CRUD
func GetObjectPrivileges() map[string][]string {
	return map[string][]string{
		commonpb.ObjectType_Collection.String(): {
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeLoad.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeRelease.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeCompaction.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeInsert.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeDelete.String()),

			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeGetStatistics.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeCreateIndex.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeIndexDetail.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeDropIndex.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeSearch.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeFlush.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeQuery.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeLoadBalance.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeImport.String()),
		},
		commonpb.ObjectType_Global.String(): {
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeAll.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeCreateCollection.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeDropCollection.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeDescribeCollection.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeShowCollections.String()),

			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeCreateOwnership.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeDropOwnership.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeSelectOwnership.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeManageOwnership.String()),
		},
		commonpb.ObjectType_User.String(): {
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeUpdateUser.String()),
			MetaStore2API(commonpb.ObjectPrivilege_PrivilegeSelectUser.String()),
		},
	}
}

// MetaStore2API convert meta-store's privilege name to api's
// example: PrivilegeAll -> All
func MetaStore2API(name string) string {
	return name[strings.Index(name, PrivilegeWord)+len(PrivilegeWord):]
}

func PrivilegeNameForAPI(name string) string {
	_, ok := commonpb.ObjectPrivilege_value[name]
	if !ok {
		return ""
	}
	return MetaStore2API(name)
}

func PrivilegeNameForMetastore(name string) string {
	dbPrivilege := PrivilegeWord + name
	_, ok := commonpb.ObjectPrivilege_value[dbPrivilege]
	if !ok {
		return ""
	}
	return dbPrivilege
}
