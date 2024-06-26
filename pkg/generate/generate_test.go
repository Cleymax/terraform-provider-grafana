package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/terraform-provider-grafana/v3/internal/testutils"
	"github.com/grafana/terraform-provider-grafana/v3/pkg/generate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestAccGenerate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}
	testutils.CheckOSSTestsEnabled(t)

	cases := []struct {
		name           string
		config         string
		generateConfig func(cfg *generate.Config)
		check          func(t *testing.T, tempDir string)
	}{
		{
			name:   "dashboard",
			config: testutils.TestAccExample(t, "resources/grafana_dashboard/resource.tf"),
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/dashboard", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
		{
			name:   "dashboard-json",
			config: testutils.TestAccExample(t, "resources/grafana_dashboard/resource.tf"),
			generateConfig: func(cfg *generate.Config) {
				cfg.Format = generate.OutputFormatJSON
			},
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/dashboard-json", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
		{
			name:   "dashboard-filter-strict",
			config: testutils.TestAccExample(t, "resources/grafana_dashboard/resource.tf"),
			generateConfig: func(cfg *generate.Config) {
				cfg.IncludeResources = []string{"grafana_dashboard._1_my-dashboard-uid"}
			},
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/dashboard-filtered", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
		{
			name:   "dashboard-filter-wildcard-on-resource-type",
			config: testutils.TestAccExample(t, "resources/grafana_dashboard/resource.tf"),
			generateConfig: func(cfg *generate.Config) {
				cfg.IncludeResources = []string{"*._1_my-dashboard-uid"}
			},
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/dashboard-filtered", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
		{
			name:   "dashboard-filter-wildcard-on-resource-name",
			config: testutils.TestAccExample(t, "resources/grafana_dashboard/resource.tf"),
			generateConfig: func(cfg *generate.Config) {
				cfg.IncludeResources = []string{"grafana_dashboard.*"}
			},
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/dashboard-filtered", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
		{
			name:   "filter-all",
			config: testutils.TestAccExample(t, "resources/grafana_dashboard/resource.tf"),
			generateConfig: func(cfg *generate.Config) {
				cfg.IncludeResources = []string{"doesnot.exist"}
			},
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/empty", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
		{
			name: "alerting-in-org",
			config: func() string {
				content, err := os.ReadFile("testdata/generate/alerting-in-org.tf")
				require.NoError(t, err)
				return string(content)
			}(),
			check: func(t *testing.T, tempDir string) {
				assertFiles(t, tempDir, "testdata/generate/alerting-in-org", []string{
					".terraform",
					".terraform.lock.hcl",
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV5ProviderFactories: testutils.ProtoV5ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: tc.config,
						Check: func(s *terraform.State) error {
							tempDir := t.TempDir()
							config := generate.Config{
								OutputDir:       tempDir,
								Clobber:         true,
								Format:          generate.OutputFormatHCL,
								ProviderVersion: "v3.0.0",
								Grafana: &generate.GrafanaConfig{
									URL:  "http://localhost:3000",
									Auth: "admin:admin",
								},
							}
							if tc.generateConfig != nil {
								tc.generateConfig(&config)
							}

							require.NoError(t, generate.Generate(context.Background(), &config))
							tc.check(t, tempDir)

							return nil
						},
					},
				},
			})
		})
	}
}

// assertFiles checks that all files in the "expectedFilesDir" directory match the files in the "gotFilesDir" directory.
func assertFiles(t *testing.T, gotFilesDir, expectedFilesDir string, ignoreDirEntries []string) {
	t.Helper()
	assertFilesSubdir(t, gotFilesDir, expectedFilesDir, "", ignoreDirEntries)
}

func assertFilesSubdir(t *testing.T, gotFilesDir, expectedFilesDir, subdir string, ignoreDirEntries []string) {
	t.Helper()

	originalGotFilesDir := gotFilesDir
	originalExpectedFilesDir := expectedFilesDir
	if subdir != "" {
		gotFilesDir = filepath.Join(gotFilesDir, subdir)
		expectedFilesDir = filepath.Join(expectedFilesDir, subdir)
	}

	// Check that all generated files are expected (recursively)
	gotFiles, err := os.ReadDir(gotFilesDir)
	if err != nil {
		t.Logf("folder %s was not generated as expected", subdir)
		t.Fail()
		return
	}
	for _, gotFile := range gotFiles {
		relativeName := filepath.Join(subdir, gotFile.Name())
		if slices.Contains(ignoreDirEntries, relativeName) {
			continue
		}

		if gotFile.IsDir() {
			assertFilesSubdir(t, originalGotFilesDir, originalExpectedFilesDir, filepath.Join(subdir, gotFile.Name()), ignoreDirEntries)
			continue
		}

		if _, err := os.Stat(filepath.Join(expectedFilesDir, gotFile.Name())); err != nil {
			t.Logf("file %s was generated but wasn't expected", relativeName)
			t.Fail()
		}
	}

	// Verify the contents of the generated files (recursively)
	// All files in the expected directory should be present in the generated directory
	expectedFiles, err := os.ReadDir(expectedFilesDir)
	if err != nil {
		t.Logf("folder %s was generated but wasn't expected", subdir)
		t.Fail()
		return
	}
	for _, expectedFile := range expectedFiles {
		if expectedFile.IsDir() {
			assertFilesSubdir(t, originalGotFilesDir, originalExpectedFilesDir, filepath.Join(subdir, expectedFile.Name()), ignoreDirEntries)
			continue
		}
		expectedContent, err := os.ReadFile(filepath.Join(expectedFilesDir, expectedFile.Name()))
		require.NoError(t, err)

		gotContent, err := os.ReadFile(filepath.Join(gotFilesDir, expectedFile.Name()))
		require.NoError(t, err)

		assert.Equal(t, strings.TrimSpace(string(expectedContent)), strings.TrimSpace(string(gotContent)))
	}
}
