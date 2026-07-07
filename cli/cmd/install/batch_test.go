// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"testing"

	oasfv1alpha1 "buf.build/gen/go/agntcy/oasf/protocolbuffers/go/agntcy/oasf/types/v1alpha1"
	corev1 "github.com/agntcy/dir/api/core/v1"
	"github.com/stretchr/testify/require"
)

func TestValidateInstallInvocationMutuallyExclusive(t *testing.T) {
	orig := opts

	defer func() { opts = orig }()

	opts.filters.Modules = []string{"integration/mcp"}

	err := validateInstallInvocation(true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mutually exclusive")
}

func TestValidateInstallInvocationSingleOK(t *testing.T) {
	orig := opts

	defer func() { opts = orig }()

	opts.filters = options{}.filters

	require.NoError(t, validateInstallInvocation(true))
}

func TestValidateInstallInvocationBatchOK(t *testing.T) {
	orig := opts

	defer func() { opts = orig }()

	opts.filters.Modules = []string{"integration/mcp"}

	require.NoError(t, validateInstallInvocation(false))
}

func TestRecordLabel(t *testing.T) {
	require.Equal(t, "agent-a:1.0.0", recordLabel(corev1.New(&oasfv1alpha1.Record{
		Name:    "agent-a",
		Version: "1.0.0",
	})))
	require.Equal(t, "agent-b", recordLabel(corev1.New(&oasfv1alpha1.Record{Name: "agent-b"})))
}

func TestSelectRecordsLatestByName(t *testing.T) {
	orig := opts

	defer func() { opts = orig }()

	opts.allVersions = false

	recs := []*corev1.Record{
		corev1.New(&oasfv1alpha1.Record{Name: "a", Version: "1.0.0"}),
		corev1.New(&oasfv1alpha1.Record{Name: "a", Version: "2.0.0"}),
		corev1.New(&oasfv1alpha1.Record{Name: "b", Version: "1.0.0"}),
	}

	selected := selectRecords(recs)
	require.Len(t, selected, 2)
	require.Equal(t, "2.0.0", selected[0].GetVersion())
	require.Equal(t, "b", selected[1].GetName())
}

func TestSelectRecordsAllVersions(t *testing.T) {
	orig := opts

	defer func() { opts = orig }()

	opts.allVersions = true

	recs := []*corev1.Record{
		corev1.New(&oasfv1alpha1.Record{Name: "a", Version: "1.0.0"}),
		corev1.New(&oasfv1alpha1.Record{Name: "a", Version: "2.0.0"}),
	}

	require.Len(t, selectRecords(recs), 2)
}

func TestFormatSkippedSummary(t *testing.T) {
	out := formatSkippedSummary([]skippedRecord{
		{label: "bare", reason: "no installable module"},
	})
	require.Contains(t, out, "Skipped records")
	require.Contains(t, out, "bare")
	require.Contains(t, out, "no installable module")
}

func TestBatchInstallSkipsUnsuitableRecords(t *testing.T) {
	orig := opts

	defer func() { opts = orig }()

	_, err := deriveArtifacts(loadRecord(t, "bare.json"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "no installable")
}
