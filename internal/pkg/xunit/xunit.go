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

// Package xunit contains functions for parsing XML files containing .NET test result(s) in xUnit's v2+ XML format.
// More information regarding this format can be found @ https://xunit.net/docs/format-xml-v2.
package xunit

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/kdeconinck/dtvisual/internal/pkg/slices"
)

// A result is the top-level element of the document. It's the result of a `dotnet test` operation in xUnit's v2+ XML
// format.
type result struct {
	XMLName       xml.Name   `xml:"assemblies"`
	Computer      string     `xml:"computer,attr"`
	FinishRTF     string     `xml:"finish-rtf,attr"`
	ID            string     `xml:"id,attr"`
	SchemaVersion string     `xml:"schema-version,attr"`
	StartRTF      string     `xml:"start-rtf,attr"`
	Timestamp     string     `xml:"timestamp,attr"`
	User          string     `xml:"user,attr"`
	Assemblies    []assembly `xml:"assembly"`
}

// An assembly contains information about the run of a single test assembly.
// This includes environmental information.
type assembly struct {
	ConfigFile      string       `xml:"config-file,attr"`
	Environment     string       `xml:"environment,attr"`
	ErrorCount      int          `xml:"errors,attr"`
	FailedCount     int          `xml:"failed,attr"`
	FinishRTF       string       `xml:"finish-rtf,attr"`
	ID              string       `xml:"id,attr"`
	FullName        string       `xml:"name,attr"`
	NotRunCount     int          `xml:"not-run,attr"`
	PassedCount     int          `xml:"passed,attr"`
	RunDate         string       `xml:"run-date,attr"`
	RunTime         string       `xml:"run-time,attr"`
	SkippedCount    int          `xml:"skipped,attr"`
	StartRTF        string       `xml:"start-rtf,attr"`
	TargetFramework string       `xml:"target-framework,attr"`
	TestFramework   string       `xml:"test-framework,attr"`
	Time            float32      `xml:"time,attr"`
	TimeRTF         string       `xml:"time-rtf,attr"`
	Total           int          `xml:"total,attr"`
	Collections     []collection `xml:"collection"`
	ErrorSet        errorSet     `xml:"errors"`
}

// A collection contains information about the run of a single test collection.
type collection struct {
	ID           string `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	FailedCount  int    `xml:"failed,attr"`
	NotRunCount  int    `xml:"not-run,attr"`
	PassedCount  int    `xml:"passed,attr"`
	SkippedCount int    `xml:"skipped,attr"`
	Time         string `xml:"time,attr"`
	TimeRTF      string `xml:"time-rtf,attr"`
	TotalCount   int    `xml:"total,attr"`
	Tests        []test `xml:"test"`
}

// A test contains information about the run of a single test.
type test struct {
	ID         string     `xml:"id,attr"`
	Method     string     `xml:"method,attr"`
	Name       string     `xml:"name,attr"`
	Result     string     `xml:"result,attr"`
	SourceFile string     `xml:"source-file,attr"`
	SourceLine string     `xml:"source-line,attr"`
	Time       float32    `xml:"time,attr"`
	TimeRTF    string     `xml:"time-rtf,attr"`
	Type       string     `xml:"type,attr"`
	Failure    failure    `xml:"failure"`
	Output     string     `xml:"output"`
	Reason     string     `xml:"reason"`
	TraitSet   traitSet   `xml:"traits"`
	WarningSet warningSet `xml:"warnings"`
}

// A failure contains information a test failure.
type failure struct {
	ExceptionType string `xml:"exception-type,attr"`
	Message       string `xml:"message"`
	StackTrace    string `xml:"stack-trace"`
}

// A traitSet contains a collection of trait elements.
type traitSet struct {
	Traits []trait `xml:"trait"`
}

// A trait contains a single trait name/value pair.
type trait struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// A warningSet contains a collection of warning elements.
type warningSet struct {
	Warnings []string `xml:"warning"`
}

// An errorSet contains a collection of error elements.
type errorSet struct {
	Errors []err `xml:"error"`
}

// An err contains information about an environment failure that happened outside the scope of running a single unit
// test (for example, an exception thrown while disposing of a fixture object).
type err struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// TestRun contains the relevant information stored in xUnit's v2+ XML format.
type TestRun struct {
	Computer     string     // The name of the computer that produced xUnit's v2+ XML format.
	User         string     // The name of the user that produced xUnit's v2+ XML format.
	StartTimeRTF string     // The time the first assembly started running.
	EndTimeRTF   string     // The time the last assembly finished running.
	Timestamp    string     // The time the first assembly started running.
	Assemblies   []Assembly // The assemblies that are part of this test run.
}

// Assembly contains information about the run of a single test assembly.
// This includes environmental information.
type Assembly struct {
	Name        string       // The full name of the assembly.
	ErrorCount  int          // The total number of environmental errors experienced in the assembly.
	PassedCount int          // The total number of test cases in the assembly which passed.
	FailedCount int          // The total number of test cases in the assembly which failed.
	NotRunCount int          // The total number of test cases that weren't run.
	TotalCount  int          // The total number of test cases in the assembly.
	RunDate     string       // The date when the test run started.
	RunTime     string       // The time when the test run started.
	Time        string       // The time spent running the tests in the assembly.
	Tests       []*TestGroup // All the tests of the assembly, grouped by trait.
}

// TestGroup is a group of tests.
type TestGroup struct {
	Name   string       // The name of the group.
	Tests  []TestCase   // The tests that belong to this group.
	Groups []*TestGroup // The subgroups of this group.
}

// TestCase contains information about a single test.
type TestCase struct {
	Name   string // The name of the test, in human-readable format.
	Result string // The status of the test.
}

// Load returns a TestRun constructed from the data in rdr.
func Load(rdr io.Reader) (TestRun, error) {
	data, err := unmarshal(rdr)

	if err != nil {
		return TestRun{}, err
	}

	testRun := TestRun{
		Computer:     data.Computer,
		User:         data.User,
		StartTimeRTF: data.StartRTF,
		EndTimeRTF:   data.FinishRTF,
		Timestamp:    data.Timestamp,
		Assemblies:   make([]Assembly, 0, len(data.Assemblies)),
	}

	// Loop over each assembly.
	for _, assembly := range data.Assemblies {
		testRun.Assemblies = append(testRun.Assemblies, Assembly{
			Name:        assembly.name(),
			ErrorCount:  assembly.ErrorCount,
			PassedCount: assembly.PassedCount,
			FailedCount: assembly.FailedCount,
			NotRunCount: assembly.NotRunCount,
			TotalCount:  assembly.Total,
			RunDate:     assembly.RunDate,
			RunTime:     assembly.RunTime,
			Time:        assembly.TimeRTF,
			Tests:       assembly.groupTests(),
		})
	}

	return testRun, nil
}

// Returns a result, constructed from the data in rdr.
func unmarshal(rdr io.Reader) (result, error) {
	var res result

	if bytes, err := io.ReadAll(rdr); err == nil {
		if err := xml.Unmarshal(bytes, &res); err != nil {
			return result{}, err
		}
	}

	return res, nil
}

// Returns the name of the assembly.
func (assembly *assembly) name() string {
	if strings.Contains(assembly.FullName, "/") {
		return assembly.FullName[strings.LastIndex(assembly.FullName, "/")+1:]
	}

	return assembly.FullName[strings.LastIndex(assembly.FullName, "\\")+1:]
}

// Returns a map of tests, grouped per trait of the assembly.
func (assembly *assembly) groupTests() []*TestGroup {
	uniqueTraits := assembly.uniqueTraits()
	resultSet := make([]*TestGroup, 0, len(uniqueTraits))

	if !assembly.hasTests() {
		return resultSet
	}

	for idx, trait := range uniqueTraits {
		cGroup := &TestGroup{Name: trait}
		resultSet = append(resultSet, cGroup)

		for _, tc := range assembly.testsWithTrait(trait) {
			if tc.hasDisplayName() || !tc.isNested() {
				cGroup.Tests = append(cGroup.Tests, TestCase{Name: tc.Name, Result: tc.Result})
			} else {
				for idx, nn := range tc.nestedNames() {
					var sGroup *TestGroup

					for _, group := range cGroup.Groups {
						if group.Name == nn {
							sGroup = group

							break
						}
					}

					if sGroup == nil {
						sGroup = &TestGroup{Name: nn}
						cGroup.Groups = append(cGroup.Groups, sGroup)
					}

					if idx == len(tc.nestedNames())-1 {
						sGroup.Tests = append(sGroup.Tests, tc)
					}

					cGroup = sGroup
				}

				cGroup = resultSet[idx]
			}
		}
	}

	return resultSet
}

// Returns true if the assembly has tests, false otherwise.
func (assembly *assembly) hasTests() bool {
	for _, collection := range assembly.Collections {
		if len(collection.Tests) > 0 {
			return true
		}
	}

	return false
}

// Returns all all the unique trait(s).
func (assembly *assembly) uniqueTraits() []string {
	resultSet := make([]string, 0)
	resultSet = append(resultSet, "")

	for _, collection := range assembly.Collections {
		for _, t := range collection.Tests {
			for _, tTrait := range t.TraitSet.Traits {
				traitName := fmt.Sprintf("%s - %s", tTrait.Name, tTrait.Value)

				if !slices.Contains(resultSet, traitName) {
					resultSet = append(resultSet, traitName)
				}
			}
		}
	}

	return resultSet
}

// Returns all the tests of the assembly that belong to a given trait.
func (assembly *assembly) testsWithTrait(traitName string) []TestCase {
	resultSet := make([]TestCase, 0)

	for _, collection := range assembly.Collections {
		for _, t := range collection.Tests {
			if traitName == "" && len(t.TraitSet.Traits) == 0 {
				resultSet = append(resultSet, TestCase{Name: t.Name, Result: t.Result})
			} else {
				for _, tTrait := range t.TraitSet.Traits {
					if fmt.Sprintf("%s - %s", tTrait.Name, tTrait.Value) == traitName {
						resultSet = append(resultSet, TestCase{Name: t.Name, Result: t.Result})
					}
				}
			}

		}
	}

	return resultSet
}

// Returns true if test has a display name, false otherwise.
// A test has a display name if it contains spaces and NO plus signs.
func (tc *TestCase) hasDisplayName() bool {
	return strings.Contains(tc.Name, " ") && !strings.Contains(tc.Name, "+")
}

// Returns true if the test is nested, false otherwise.
// A test is nested if the name of the test contains one or more plus signs.
func (tc *TestCase) isNested() bool {
	return strings.Contains(tc.Name, "+")
}

// Returns the name of each nested level.
func (tc *TestCase) nestedNames() []string {
	retVal := strings.Split(tc.Name, "+")

	parts := make([]string, 0, len(retVal[1:len(retVal)-1]))
	parts = append(parts, retVal[0][strings.LastIndex(retVal[0], ".")+1:])
	parts = append(parts, retVal[1:len(retVal)-1]...)
	parts = append(parts, strings.Split(retVal[len(retVal)-1], ".")[0])

	return parts
}
