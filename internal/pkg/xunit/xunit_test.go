// =====================================================================================================================
// = LICENSE:       Copyright (c) 2023 Kevin De Coninck
// =
// =                Permission is hereby granted, free of charge, to any person
// =                obtaining a copy of this software and associated documentation
// =                files (the "Software"), to deal in the Software without
// =                restriction, including without limitation the rights to use,
// =                copy, modify, merge, publish, distribute, sublicense, and/or sell
// =                copies of the Software, and to permit persons to whom the
// =                Software is furnished to do so, subject to the following
// =                conditions:
// =
// =                The above copyright notice and this permission notice shall be
// =                included in all copies or substantial portions of the Software.
// =
// =                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// =                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// =                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// =                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// =                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// =                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// =                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// =                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Quality assurance: Verify (and measure the performance) of the public API of the "xunit" package.
package xunit_test

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/kdeconinck/dtvisual/internal/pkg/assert"
	"github.com/kdeconinck/dtvisual/internal/pkg/xunit"
)

// UT: Load an XML file containing a .NET test result.
func TestLoad(t *testing.T) {
	t.Parallel() // Enable "parallel" execution.

	for _, tc := range []struct {
		xmlData string
		want    xunit.TestRun
		wantErr bool
	}{
		{
			xmlData: "{}",
			wantErr: true,
		},
		{
			xmlData: "<assemblies computer=\"WIN11\" user=\"Kevin\" timestamp=\"07/10/2023 20:53:19\" start-rtf=\"2000-12-01\" finish-rtf=\"2001-12-01\" timestamp=\"2001-12-02\">\n" +
				"  <assembly name=\"C:\\Parent\\Sub\\App.dll\" errors=\"1\" failed=\"2\" passed=\"3\" not-run=\"4\" total=\"5\" run-date=\"07/10/2023\" run-time=\"20:53:19\" time-rtf=\"2000-12-01\">\n" +
				"  </assembly>\n" +
				"</assemblies>",
			want: xunit.TestRun{
				Computer:     "WIN11",
				User:         "Kevin",
				StartTimeRTF: "2000-12-01",
				EndTimeRTF:   "2001-12-01",
				Timestamp:    "2001-12-02",
				Assemblies: []xunit.Assembly{
					{
						Name:        "App.dll",
						ErrorCount:  1,
						PassedCount: 3,
						FailedCount: 2,
						NotRunCount: 4,
						TotalCount:  5,
						RunDate:     "07/10/2023",
						RunTime:     "20:53:19",
						Time:        "2000-12-01",
						Tests:       make(map[string][]*xunit.TestGroup),
					},
				},
			},
		},

		{
			xmlData: "<assemblies computer=\"WIN11\" user=\"Kevin\" timestamp=\"07/10/2023 20:53:19\" start-rtf=\"2000-12-01\" finish-rtf=\"2001-12-01\" timestamp=\"2001-12-02\">\n" +
				"  <assembly name=\"~/parent/sub/app.dll\" errors=\"1\" failed=\"2\" passed=\"3\" not-run=\"4\" total=\"5\" run-date=\"07/10/2023\" run-time=\"20:53:19\" time-rtf=\"2000-12-01\">\n" +
				"    <collection>\n" +

				// NOTE: A test which has a display name (it contains spaces, and NO `+` sign).
				"      <test name=\"A test with a display name.\" result=\"Pass\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A NON nested test without a display name (it contains NO spaces, and NO `+` character).
				"      <test name=\"NS1.Class.SubClass.TestClass.TestMethod\" result=\"Fail\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A nested test (it contains the `+` character).
				"      <test name=\"NS1.Class.SubClass.TestClass+Method+Scenario+SubScenario.Result\" result=\"Pass\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A nested test (it contains the `+` character), in already existing group.
				"      <test name=\"NS1.Class.SubClass.TestClass+Method+Scenario2+SubScenario.Result\" result=\"Pass\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A test which has a display name (it contains spaces, and NO `+` sign), belonging to a single trait.
				"      <test name=\"A test with a display name (with a trait).\" result=\"Pass\">\n" +
				"        <traits>\n" +
				"		   <trait name=\"Category\" value=\"Unit\" />\n" +
				"        </traits>\n" +
				"      </test>\n" +

				// NOTE: A test which has a display name (it contains spaces, and NO `+` sign), belonging to multiple traits.
				"      <test name=\"A test with a display name (with multiple traits).\" result=\"Pass\">\n" +
				"        <traits>\n" +
				"		   <trait name=\"Category\" value=\"Unit\" />\n" +
				"		   <trait name=\"Timing\" value=\"Slow\" />\n" +
				"        </traits>\n" +
				"      </test>\n" +
				"    </collection>\n" +
				"  </assembly>\n" +
				"</assemblies>",
			want: xunit.TestRun{
				Computer:     "WIN11",
				User:         "Kevin",
				StartTimeRTF: "2000-12-01",
				EndTimeRTF:   "2001-12-01",
				Timestamp:    "2001-12-02",
				Assemblies: []xunit.Assembly{
					{
						Name:        "app.dll",
						ErrorCount:  1,
						PassedCount: 3,
						FailedCount: 2,
						NotRunCount: 4,
						TotalCount:  5,
						RunDate:     "07/10/2023",
						RunTime:     "20:53:19",
						Time:        "2000-12-01",
						Tests: map[string][]*xunit.TestGroup{
							"": {
								{
									Tests: []xunit.TestCase{
										{
											Name:   "A test with a display name.",
											Result: "Pass",
										},
										{
											Name:   "NS1.Class.SubClass.TestClass.TestMethod",
											Result: "Fail",
										},
									},
									Groups: []*xunit.TestGroup{
										{
											Name:  "TestClass",
											Tests: nil,
											Groups: []*xunit.TestGroup{
												{
													Name:  "Method",
													Tests: nil,
													Groups: []*xunit.TestGroup{
														{
															Name:  "Scenario",
															Tests: nil,
															Groups: []*xunit.TestGroup{
																{
																	Name: "SubScenario",
																	Tests: []xunit.TestCase{
																		{
																			Name:   "NS1.Class.SubClass.TestClass+Method+Scenario+SubScenario.Result",
																			Result: "Pass",
																		},
																	},
																},
															},
														},
														{
															Name:  "Scenario2",
															Tests: nil,
															Groups: []*xunit.TestGroup{
																{
																	Name: "SubScenario",
																	Tests: []xunit.TestCase{
																		{
																			Name:   "NS1.Class.SubClass.TestClass+Method+Scenario2+SubScenario.Result",
																			Result: "Pass",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							"Category - Unit": {
								{
									Tests: []xunit.TestCase{
										{
											Name:   "A test with a display name (with a trait).",
											Result: "Pass",
										},
										{
											Name:   "A test with a display name (with multiple traits).",
											Result: "Pass",
										},
									},
								},
							},
							"Timing - Slow": {
								{
									Tests: []xunit.TestCase{
										{
											Name:   "A test with a display name (with multiple traits).",
											Result: "Pass",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		// HELPER FUNCTIONS.
		fmtValue := func(v xunit.TestRun) string {
			b, _ := json.MarshalIndent(v, "", "  ")
			res := strings.Replace(string(b), "\n", "\n            ", -1)

			return res
		}

		fmtXml := func(v string) string {
			v = strings.Replace(v, "\n", "\n            ", -1)

			return v
		}

		// ARRANGE.
		rdr := strings.NewReader(tc.xmlData)

		// ACT.
		got, err := xunit.Load(rdr)

		// ASSERT.
		if tc.wantErr {
			assert.NotNil(t, err, "", "\n\n"+
				"UT Name:    Parse an invalid XML file containing a .NET test result in xUnit's v2+ XML format.\n"+
				"XML Input:  %s\n"+
				"\033[32mExpected:   Error, NOT <nil>\033[0m\n"+
				"\033[31mActual:     Error, %v\033[0m\n\n", fmtXml(tc.xmlData), err)
		}

		if !tc.wantErr {
			assert.Nil(t, err, "", "\n\n"+
				"UT Name:    Parse an invalid XML file containing a .NET test result in xUnit's v2+ XML format.\n"+
				"XML Input:  %s\n"+
				"\033[32mExpected:   Error, <nil>\033[0m\n"+
				"\033[31mActual:     Error, %v\033[0m\n\n", fmtXml(tc.xmlData), err)
		}

		assert.EqualFn(t, got, tc.want, func(got xunit.TestRun, want xunit.TestRun) bool {
			return reflect.DeepEqual(got, want)
		}, "", "\n\n"+
			"UT Name:    Parse an XML file containing a .NET test result in xUnit's v2+ XML format.\n"+
			"XML Input:  %s\n"+
			"\033[32mExpected:   %s\033[0m\n"+
			"\033[31mActual:     %s\033[0m\n\n",
			fmtXml(tc.xmlData), fmtValue(tc.want), fmtValue(got))
	}
}

// Benchmark: Load an XML file containing a .NET test result.
func BenchmarkLoad_MultipleAssemblies(b *testing.B) {
	xmlData := "<assemblies>\n"

	for i := 0; i < 100; i++ {
		xmlData += "  <assembly />\n"
	}

	xmlData += "</assemblies>"

	benchmarkLoad(xmlData, b)
}

// Benchmark: Load an XML file containing a .NET test result.
func BenchmarkLoad_MultipleTraits(b *testing.B) {
	xmlData := "<assemblies>\n" +
		"  <assembly>\n" +
		"    <collection>\n" +
		"      <test>\n" +
		"        <traits>\n"

	for i := 0; i < 100; i++ {
		xmlData += "		   <trait name=\"Idx\" value=\"" + strconv.Itoa(i) + "\" />\n"
	}

	xmlData += "        </traits>\n" +
		"      </test>\n" +
		"    </collection>\n" +
		"  </assembly>\n" +
		"</assemblies>"

	benchmarkLoad(xmlData, b)
}

// Benchmark: Load an XML file containing a .NET test result.
func benchmarkLoad(xmlData string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		rdr := strings.NewReader(xmlData)

		_, _ = xunit.Load(rdr)
	}
}
