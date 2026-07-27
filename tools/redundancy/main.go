// Command redundancy reports structural zombies and unreachable code in the SDK.
//
//	go run ./tools/redundancy
//	make redundancy
//
// Product goal: keep the graph lean so new markets can be added via a standard
// domain shell (watchlist → ingestor → store → facade.Emit) without fighting
// dead plugins, empty modules, or dual APIs.
//
// Requires Go 1.26+ matching go.mod. Installs deadcode on demand.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	// Product surface reachability root (fx graph strategies link).
	_ "github.com/wisp-trading/sdk/wisp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "redundancy: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	if err := ensureDeadcode(); err != nil {
		return err
	}

	fmt.Println("Wisp SDK redundancy report")
	fmt.Println("Goal: AAA market framework — one pattern, minimal time to a live bot.")
	fmt.Println()

	fmt.Println("=== 1. Structural zombies (high signal for DX) ===")
	structural := structuralFindings(root)
	for _, f := range structural {
		fmt.Println("•", f)
	}
	if len(structural) == 0 {
		fmt.Println("(none)")
	}
	fmt.Println()

	fmt.Println("=== 2. Unexported dead code (from product surface root) ===")
	fmt.Println("Root: tools/redundancy main blank-imports github.com/wisp-trading/sdk/wisp")
	fmt.Println("Only unexported symbols (true internal dead). Public API may look unused")
	fmt.Println("because strategy authors call it outside this module — that is expected.")
	fmt.Println()

	cmd := exec.Command("deadcode", "./tools/redundancy")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.CombinedOutput()
	text := string(out)
	if err != nil && len(bytes.TrimSpace(out)) == 0 {
		return fmt.Errorf("deadcode: %w\n%s", err, text)
	}

	unexported := filterUnexportedDead(text)
	if len(unexported) == 0 {
		fmt.Println("(no unexported unreachable funcs)")
	} else {
		for _, l := range unexported {
			fmt.Println(l)
		}
	}
	fmt.Println()

	fmt.Println("=== 3. Export-looking dead (review: maybe public API or truly orphan) ===")
	exported := filterExportedDead(text)
	// Cap noise — list packages with counts
	byPkg := map[string]int{}
	for _, l := range exported {
		pkg := packageOfDeadLine(l)
		byPkg[pkg]++
	}
	var pkgs []string
	for p := range byPkg {
		pkgs = append(pkgs, p)
	}
	sort.Strings(pkgs)
	if len(pkgs) == 0 {
		fmt.Println("(none)")
	} else {
		for _, p := range pkgs {
			fmt.Printf("• %s (%d symbols) — sample: %s\n", p, byPkg[p], firstInPkg(exported, p))
		}
	}
	fmt.Println()

	fmt.Printf("Summary: %d structural · %d unexported dead · %d exported-looking\n",
		len(structural), len(unexported), len(exported))
	fmt.Println()
	fmt.Println("Next (manual): market shell DRY, instruments routing, hard connector errors.")
	fmt.Println("Re-run: go run ./tools/redundancy   or   make redundancy")
	return nil
}

func moduleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", wd)
		}
		dir = parent
	}
}

func ensureDeadcode() error {
	if _, err := exec.LookPath("deadcode"); err == nil {
		return nil
	}
	fmt.Fprintln(os.Stderr, "installing golang.org/x/tools/cmd/deadcode@latest …")
	cmd := exec.Command("go", "install", "golang.org/x/tools/cmd/deadcode@latest")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var deadLineRE = regexp.MustCompile(`^(.+\.go:\d+:\d+): unreachable func: (.+)$`)

func filterUnexportedDead(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if skipDeadLine(line) {
			continue
		}
		m := deadLineRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if isUnexportedFunc(m[2]) {
			out = append(out, line)
		}
	}
	sort.Strings(out)
	return out
}

func filterExportedDead(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if skipDeadLine(line) {
			continue
		}
		m := deadLineRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if !isUnexportedFunc(m[2]) {
			out = append(out, line)
		}
	}
	sort.Strings(out)
	return out
}

func skipDeadLine(line string) bool {
	if line == "" {
		return true
	}
	if strings.Contains(line, "tools/redundancy") {
		return true
	}
	if strings.Contains(line, "/mocks/") {
		return true
	}
	if strings.Contains(line, "requires newer Go") || strings.Contains(line, "no main packages") {
		return true
	}
	return false
}

// isUnexportedFunc reports whether the deadcode func name is unexported.
// Handles "pkg.Type.Method", "Type.method", "funcName".
func isUnexportedFunc(name string) bool {
	// Take last path segment of method (after last '.')
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsLower(r)
}

func packageOfDeadLine(line string) string {
	// path/file.go:line:col: ...
	idx := strings.Index(line, ".go:")
	if idx < 0 {
		return "?"
	}
	path := line[:idx]
	// drop filename
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return path
	}
	return dir
}

func firstInPkg(lines []string, pkg string) string {
	for _, l := range lines {
		if packageOfDeadLine(l) == pkg {
			if m := deadLineRE.FindStringSubmatch(l); m != nil {
				return m[2]
			}
			return l
		}
	}
	return ""
}

func structuralFindings(root string) []string {
	var hits []string
	checkFile := func(rel, why string) {
		if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
			hits = append(hits, fmt.Sprintf("%s — %s", rel, why))
		}
	}
	checkFile("pkg/signal/module.go", "empty fx Module still wired (domain packages own signal builders)")
	checkFile("pkg/testing/module.go", "re-exports wisp.Module only; not a product surface")

	if b, err := os.ReadFile(filepath.Join(root, "pkg/types/execution/hook.go")); err == nil {
		if bytes.Contains(b, []byte("HookPlugin")) {
			hits = append(hits, "pkg/types/execution.HookPlugin — .so plugin interface after plugin removal")
		}
	}
	if b, err := os.ReadFile(filepath.Join(root, "pkg/modules.go")); err == nil {
		if bytes.Contains(b, []byte("signal.Module")) {
			hits = append(hits, "pkg/modules.go wires signal.Module (empty)")
		}
		if bytes.Contains(b, []byte("package packages")) {
			hits = append(hits, `pkg/modules.go package name "packages" (non-idiomatic; prefer sdk/pkg Module re-export)`)
		}
	}
	if b, err := os.ReadFile(filepath.Join(root, "pkg/types/strategy/strategy.go")); err == nil {
		if bytes.Contains(b, []byte("CashCarry")) || bytes.Contains(b, []byte("VolumeMaximizer")) {
			hits = append(hits, "pkg/types/strategy — hardcoded product strategy names (CashCarry, …)")
		}
	}
	if b, err := os.ReadFile(filepath.Join(root, "pkg/types/lifecycle/lifecycle.go")); err == nil {
		if bytes.Contains(b, []byte("WispExecutor")) {
			hits = append(hits, "pkg/types/lifecycle — fossil comment WispExecutor")
		}
	}
	// Nested facades slow “add a market” onboarding
	for _, m := range []string{"spot", "perp", "options"} {
		if st, err := os.Stat(filepath.Join(root, "pkg/markets", m, m)); err == nil && st.IsDir() {
			hits = append(hits, fmt.Sprintf("pkg/markets/%s/%s — nested same-name facade (tax when cloning market shell)", m, m))
		}
	}
	// prediction uses predict/ not prediction/prediction — asymmetry
	if st, err := os.Stat(filepath.Join(root, "pkg/markets/prediction/predict")); err == nil && st.IsDir() {
		hits = append(hits, "pkg/markets/prediction/predict — facade name differs from spot/perp/options pattern")
	}
	// base activity stores that deadcode often flags entirely
	checkFile("pkg/markets/base/store/activity/position/store.go", "legacy position store package (often fully unreachable)")
	checkFile("pkg/markets/base/store/activity/trade/store.go", "legacy trade store package (often fully unreachable)")

	sort.Strings(hits)
	return hits
}
