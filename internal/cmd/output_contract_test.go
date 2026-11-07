package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/aarondpn/redmine-cli/v2/internal/cmdutil"
)

type outputContract string

const (
	outputContractStructured     outputContract = "structured"
	outputContractRawPassthrough outputContract = "raw_passthrough"
	outputContractInteractive    outputContract = "interactive_only"
)

type commandContract struct {
	Mode   outputContract
	Reason string
}

// commandOutputContracts is an explicit registry for every runnable,
// user-facing command. Adding a new command requires classifying how it
// handles output, which turns "forgotten command" regressions into a test
// failure instead of a review surprise.
var commandOutputContracts = map[string]commandContract{
	"rdflow api":                      {Mode: outputContractStructured},
	"rdflow attachments download":     {Mode: outputContractStructured},
	"rdflow attachments get":          {Mode: outputContractStructured},
	"rdflow attachments update":       {Mode: outputContractStructured},
	"rdflow attachments delete":       {Mode: outputContractStructured},
	"rdflow attachments thumbnail":    {Mode: outputContractStructured},
	"rdflow auth list":                {Mode: outputContractStructured},
	"rdflow auth login":               {Mode: outputContractInteractive, Reason: "runs an interactive Huh form to collect credentials"},
	"rdflow auth logout":              {Mode: outputContractInteractive, Reason: "may prompt for interactive confirmation before deleting a profile"},
	"rdflow auth status":              {Mode: outputContractStructured},
	"rdflow auth switch":              {Mode: outputContractInteractive, Reason: "may prompt with an interactive profile selector"},
	"rdflow categories list":          {Mode: outputContractStructured},
	"rdflow categories get":           {Mode: outputContractStructured},
	"rdflow categories create":        {Mode: outputContractStructured},
	"rdflow categories update":        {Mode: outputContractStructured},
	"rdflow categories delete":        {Mode: outputContractStructured},
	"rdflow completion":               {Mode: outputContractRawPassthrough, Reason: "streams generated shell completion scripts"},
	"rdflow config":                   {Mode: outputContractStructured},
	"rdflow custom-fields get":        {Mode: outputContractStructured},
	"rdflow custom-fields list":       {Mode: outputContractStructured},
	"rdflow files list":               {Mode: outputContractStructured},
	"rdflow files upload":             {Mode: outputContractStructured},
	"rdflow groups add-user":          {Mode: outputContractStructured},
	"rdflow groups create":            {Mode: outputContractStructured},
	"rdflow groups delete":            {Mode: outputContractStructured},
	"rdflow groups get":               {Mode: outputContractStructured},
	"rdflow groups list":              {Mode: outputContractStructured},
	"rdflow groups remove-user":       {Mode: outputContractStructured},
	"rdflow groups update":            {Mode: outputContractStructured},
	"rdflow install-skill":            {Mode: outputContractStructured},
	"rdflow issues assign":            {Mode: outputContractStructured},
	"rdflow issues browse":            {Mode: outputContractInteractive, Reason: "launches the TUI browser"},
	"rdflow issues close":             {Mode: outputContractStructured},
	"rdflow issues comment":           {Mode: outputContractStructured},
	"rdflow issues comment edit":      {Mode: outputContractStructured},
	"rdflow issues create":            {Mode: outputContractStructured},
	"rdflow issues delete":            {Mode: outputContractStructured},
	"rdflow issues get":               {Mode: outputContractStructured},
	"rdflow issues list":              {Mode: outputContractStructured},
	"rdflow issues open":              {Mode: outputContractStructured},
	"rdflow issues relations add":     {Mode: outputContractStructured},
	"rdflow issues relations get":     {Mode: outputContractStructured},
	"rdflow issues relations list":    {Mode: outputContractStructured},
	"rdflow issues relations remove":  {Mode: outputContractStructured},
	"rdflow issues reopen":            {Mode: outputContractStructured},
	"rdflow issues update":            {Mode: outputContractStructured},
	"rdflow issues watchers add":      {Mode: outputContractStructured},
	"rdflow issues watchers list":     {Mode: outputContractStructured},
	"rdflow issues watchers remove":   {Mode: outputContractStructured},
	"rdflow mcp serve":                {Mode: outputContractRawPassthrough, Reason: "runs the MCP JSON-RPC protocol over stdio"},
	"rdflow mcp tools":                {Mode: outputContractRawPassthrough, Reason: "prints a human-formatted catalog of MCP tools and groups"},
	"rdflow memberships create":       {Mode: outputContractStructured},
	"rdflow memberships delete":       {Mode: outputContractStructured},
	"rdflow memberships get":          {Mode: outputContractStructured},
	"rdflow memberships list":         {Mode: outputContractStructured},
	"rdflow memberships update":       {Mode: outputContractStructured},
	"rdflow my-account get":           {Mode: outputContractStructured},
	"rdflow news create":              {Mode: outputContractStructured},
	"rdflow news delete":              {Mode: outputContractStructured},
	"rdflow news get":                 {Mode: outputContractStructured},
	"rdflow news list":                {Mode: outputContractStructured},
	"rdflow news update":              {Mode: outputContractStructured},
	"rdflow my-account update":        {Mode: outputContractStructured},
	"rdflow projects archive":         {Mode: outputContractStructured},
	"rdflow projects create":          {Mode: outputContractStructured},
	"rdflow projects delete":          {Mode: outputContractStructured},
	"rdflow projects get":             {Mode: outputContractStructured},
	"rdflow projects list":            {Mode: outputContractStructured},
	"rdflow projects members":         {Mode: outputContractStructured},
	"rdflow projects unarchive":       {Mode: outputContractStructured},
	"rdflow projects update":          {Mode: outputContractStructured},
	"rdflow queries get":              {Mode: outputContractStructured},
	"rdflow queries list":             {Mode: outputContractStructured},
	"rdflow roles get":                {Mode: outputContractStructured},
	"rdflow roles list":               {Mode: outputContractStructured},
	"rdflow search":                   {Mode: outputContractStructured},
	"rdflow search browse":            {Mode: outputContractInteractive, Reason: "launches the TUI search browser"},
	"rdflow search issues":            {Mode: outputContractStructured},
	"rdflow search messages":          {Mode: outputContractStructured},
	"rdflow search news":              {Mode: outputContractStructured},
	"rdflow search projects":          {Mode: outputContractStructured},
	"rdflow search wiki":              {Mode: outputContractStructured},
	"rdflow statuses list":            {Mode: outputContractStructured},
	"rdflow priorities list":          {Mode: outputContractStructured},
	"rdflow activities list":          {Mode: outputContractStructured},
	"rdflow document-categories list": {Mode: outputContractStructured},
	"rdflow time delete":              {Mode: outputContractStructured},
	"rdflow time get":                 {Mode: outputContractStructured},
	"rdflow time list":                {Mode: outputContractStructured},
	"rdflow time log":                 {Mode: outputContractStructured},
	"rdflow time summary":             {Mode: outputContractStructured},
	"rdflow time update":              {Mode: outputContractStructured},
	"rdflow trackers get":             {Mode: outputContractStructured},
	"rdflow trackers list":            {Mode: outputContractStructured},
	"rdflow update":                   {Mode: outputContractStructured},
	"rdflow users create":             {Mode: outputContractStructured},
	"rdflow users delete":             {Mode: outputContractStructured},
	"rdflow users get":                {Mode: outputContractStructured},
	"rdflow users list":               {Mode: outputContractStructured},
	"rdflow users me":                 {Mode: outputContractStructured},
	"rdflow users update":             {Mode: outputContractStructured},
	"rdflow versions get":             {Mode: outputContractStructured},
	"rdflow versions list":            {Mode: outputContractStructured},
	"rdflow versions create":          {Mode: outputContractStructured},
	"rdflow versions update":          {Mode: outputContractStructured},
	"rdflow versions delete":          {Mode: outputContractStructured},
	"rdflow wiki create":              {Mode: outputContractStructured},
	"rdflow wiki delete":              {Mode: outputContractStructured},
	"rdflow wiki get":                 {Mode: outputContractStructured},
	"rdflow wiki list":                {Mode: outputContractStructured},
	"rdflow wiki update":              {Mode: outputContractStructured},
}

func TestRunnableCommandsHaveDeclaredOutputContracts(t *testing.T) {
	root := NewRootCmd("test")

	var actual []string
	walk(root, func(c *cobra.Command) {
		if c.Runnable() && !c.Hidden {
			actual = append(actual, c.CommandPath())
		}
	})
	slices.Sort(actual)

	var expected []string
	for commandPath, contract := range commandOutputContracts {
		expected = append(expected, commandPath)
		if contract.Mode == outputContractRawPassthrough || contract.Mode == outputContractInteractive {
			if strings.TrimSpace(contract.Reason) == "" {
				t.Fatalf("%s must include a reason for %s classification", commandPath, contract.Mode)
			}
		}
	}
	slices.Sort(expected)

	if !slices.Equal(actual, expected) {
		t.Fatalf("runnable command output contract registry drifted\nactual:   %v\nexpected: %v", actual, expected)
	}
}

func TestDirectStdoutWritesAreAllowlisted(t *testing.T) {
	allowed := map[string]string{
		"completion/completion.go": "shell completion generator streams directly to stdout",
		"issue/get.go":             "human journal rendering still prints directly after detail output",
		"news/get.go":              "section headings between the detail block and the comment/attachment tables",
		"update/check.go":          "background update notice writes directly to an arbitrary writer",
	}
	patterns := []string{
		"fmt.Print(",
		"fmt.Printf(",
		"fmt.Println(",
		"fmt.Fprintf(f.IOStreams.Out",
		"fmt.Fprint(f.IOStreams.Out",
		"fmt.Fprintln(f.IOStreams.Out",
		"os.Stdout",
	}

	rootDir := "."
	var offenders []string

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		src := string(data)
		for _, pattern := range patterns {
			if strings.Contains(src, pattern) {
				if _, ok := allowed[relPath]; !ok {
					offenders = append(offenders, relPath+" contains "+pattern)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk command sources: %v", err)
	}

	if len(offenders) > 0 {
		slices.Sort(offenders)
		t.Fatalf("direct stdout/stderr writes found outside the allowlist:\n%s", strings.Join(offenders, "\n"))
	}
}

// TestRawPassthroughCommandsRejectJQ ensures every command whose stdout is not
// JSON is marked so --jq fails loudly instead of being silently ignored.
func TestRawPassthroughCommandsRejectJQ(t *testing.T) {
	root := NewRootCmd("test")
	walk(root, func(c *cobra.Command) {
		contract, ok := commandOutputContracts[c.CommandPath()]
		if !ok || contract.Mode != outputContractRawPassthrough {
			return
		}
		if !cmdutil.HasNoJSONOutput(c) {
			t.Errorf("%s is %s but not marked with cmdutil.MarkNoJSONOutput", c.CommandPath(), contract.Mode)
		}
	})
}
