/*
 * Copyright 2020 Go YAML Path Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package test

import (
	"io/ioutil"
	"testing"

	"github.com/glyn/go-yamlpath/pkg/yamlpath"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRegressionSuite(t *testing.T) {
	y, err := ioutil.ReadFile("./testdata/regression_suite.yaml")
	if err != nil {
		t.Error(err)
	}

	var suite regressionSuite

	if err = yaml.Unmarshal(y, &suite); err != nil {
		t.Fatal(err)
	}

	focussed := false
	for _, tc := range suite.Testcases {
		if tc.Focus {
			focussed = true
			break
		}
	}

	tests, passed, consensus := 0, 0, 0
	for _, tc := range suite.Testcases {
		if focussed && !tc.Focus {
			continue
		}
		tests++
		if tc.Consensus.Kind > 0 {
			consensus++
		}
		if pass := t.Run(tc.Name, func(t *testing.T) {
			defer func() {
				p := recover()
				if p != nil {
					// fail on panic regardless of whether there is a consensus
					t.Fatalf("Panicked on test: %s: %v", tc.Name, p)
				}
			}()

			path, err := yamlpath.NewPath(tc.Selector)
			// if there is a consensus, check we agree with it
			if tc.Consensus.Kind > 0 {
				require.NoError(t, err, "NewPath failed with selector: %s, test: %s", tc.Selector, tc.Name)
			}

			results, err := path.Find(&tc.Document)
			// if there is a consensus, check we agree with it
			if tc.Consensus.Kind > 0 {
				require.NoError(t, err, "Find failed with selector: %s, test: %s", tc.Selector, tc.Name)
			}

			clearLineColumn(tc.Consensus.Content)
			clearLineColumn(results)

			// if there is a consensus, check we agree with it
			if tc.Consensus.Kind > 0 {
				require.Equal(t, tc.Consensus.Content, results, "Disagreed with consensus, selector: %s, test: %s", tc.Selector, tc.Name)
			}
		}); pass {
			passed++
		}
	}

	t.Logf("%d passed and %d failed of %d tests of which %d had consensus", passed, tests-passed, tests, consensus)

	if focussed {
		t.Fatalf("testcase(s) still focussed")
	}
}

func clearLineColumn(nodes []*yaml.Node) {
	for _, n := range nodes {
		n.Line = 0
		n.Column = 0
		clearLineColumn(n.Content)
	}
}

type testcase struct {
	Name      string `yaml:"id"`
	Selector  string
	Document  yaml.Node
	Consensus yaml.Node
	Focus     bool // if true, run only tests with focus set to true
}

type regressionSuite struct {
	Testcases []testcase `yaml:"queries"`
}
