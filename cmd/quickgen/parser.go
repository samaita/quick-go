package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Table represents a parsed CREATE TABLE statement.
type Table struct {
	Name    string   // original SQL table name (snake_case)
	Columns []Column
}

// GoName returns the PascalCase struct name (e.g. "user_posts" → "UserPost").
func (t Table) GoName() string {
	return toPascal(singular(t.Name))
}

// GoNamePlural returns the PascalCase plural name (e.g. "user_posts" → "UserPosts").
func (t Table) GoNamePlural() string {
	return toPascal(t.Name)
}

// PackageName returns the lowercase package name (e.g. "user_posts" → "userposts").
func (t Table) PackageName() string {
	return strings.ToLower(strings.ReplaceAll(t.Name, "_", ""))
}

// URLPath returns the kebab-case URL segment (e.g. "user_posts" → "user-posts").
func (t Table) URLPath() string {
	return strings.ReplaceAll(t.Name, "_", "-")
}

// PrimaryKey returns the first column marked as primary key, or nil.
func (t Table) PrimaryKey() *Column {
	for i := range t.Columns {
		if t.Columns[i].IsPK {
			return &t.Columns[i]
		}
	}
	return nil
}

// NonPKColumns returns columns that are not the primary key.
func (t Table) NonPKColumns() []Column {
	var cols []Column
	for _, c := range t.Columns {
		if !c.IsPK {
			cols = append(cols, c)
		}
	}
	return cols
}

// Column represents a single column in a CREATE TABLE.
type Column struct {
	Name     string // original SQL column name
	SQLType  string // normalized SQL type
	GoType   string // mapped Go type
	IsPK     bool
	NotNull  bool
	IsUnique bool
}

// GoName returns the PascalCase field name (e.g. "first_name" → "FirstName").
func (c Column) GoName() string {
	return toPascal(c.Name)
}

// JSONName returns the snake_case JSON tag name (same as SQL name).
func (c Column) JSONName() string {
	return c.Name
}

var (
	reCreateTableHeader = regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?[` + "`" + `"]?(\w+)[` + "`" + `"]?\s*\(`)
	reLineComment       = regexp.MustCompile(`--[^\n]*`)
	reBlockComment      = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// ParseDDL parses a DDL SQL file and returns all CREATE TABLE definitions found.
func ParseDDL(path string) ([]Table, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("parser: read %q: %w", path, err)
	}

	src := string(raw)
	src = reBlockComment.ReplaceAllString(src, "")
	src = reLineComment.ReplaceAllString(src, "")

	var tables []Table
	pos := 0
	for pos < len(src) {
		loc := reCreateTableHeader.FindStringSubmatchIndex(src[pos:])
		if loc == nil {
			break
		}
		// loc indices are relative to src[pos:]
		tableName := src[pos+loc[2] : pos+loc[3]]
		bodyStart := pos + loc[1] // character after the opening '('

		// Use paren-depth tracking to find the matching closing ')'
		body, end, ok := extractParenBody(src, bodyStart)
		if !ok {
			return nil, fmt.Errorf("parser: table %q: unmatched parenthesis", tableName)
		}

		cols, err := parseColumns(tableName, body)
		if err != nil {
			return nil, fmt.Errorf("parser: table %q: %w", tableName, err)
		}
		tables = append(tables, Table{Name: tableName, Columns: cols})
		pos = end
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("parser: no CREATE TABLE statements found in %q", path)
	}
	return tables, nil
}

// extractParenBody finds the content between matching parentheses, starting
// at src[start] (i.e. the character immediately after the opening '(').
// Returns the body, the position after the closing ')', and whether it succeeded.
func extractParenBody(src string, start int) (body string, end int, ok bool) {
	depth := 1
	inStr := false
	strChar := byte(0)
	for i := start; i < len(src); i++ {
		ch := src[i]
		if inStr {
			if ch == strChar && (i == 0 || src[i-1] != '\\') {
				inStr = false
			}
			continue
		}
		switch ch {
		case '\'', '"':
			inStr = true
			strChar = ch
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return src[start:i], i + 1, true
			}
		}
	}
	return "", 0, false
}

func parseColumns(tableName, body string) ([]Column, error) {
	lines := splitColumnDefs(body)
	var cols []Column

	// Collect table-level PRIMARY KEY constraint columns
	pkCols := map[string]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		upper := strings.ToUpper(line)
		if strings.HasPrefix(upper, "PRIMARY KEY") {
			// e.g.  PRIMARY KEY (id) or PRIMARY KEY (col1, col2)
			inner := extractParens(line)
			for _, col := range strings.Split(inner, ",") {
				pkCols[strings.TrimSpace(strings.ToLower(col))] = true
			}
		}
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		upper := strings.ToUpper(line)
		// Skip table constraints
		if strings.HasPrefix(upper, "PRIMARY KEY") ||
			strings.HasPrefix(upper, "UNIQUE") ||
			strings.HasPrefix(upper, "FOREIGN KEY") ||
			strings.HasPrefix(upper, "CHECK") ||
			strings.HasPrefix(upper, "INDEX") ||
			strings.HasPrefix(upper, "KEY") {
			continue
		}

		col, err := parseColumn(line)
		if err != nil {
			// Best-effort: skip unparseable lines
			continue
		}
		if pkCols[strings.ToLower(col.Name)] {
			col.IsPK = true
		}
		cols = append(cols, col)
	}

	if len(cols) == 0 {
		return nil, fmt.Errorf("no parseable columns found")
	}
	return cols, nil
}

var reColDef = regexp.MustCompile(`(?i)^[` + "`" + `"]?(\w+)[` + "`" + `"]?\s+(\w+(?:\s*\([^)]*\))?)(.*)$`)

func parseColumn(line string) (Column, error) {
	m := reColDef.FindStringSubmatch(line)
	if m == nil {
		return Column{}, fmt.Errorf("cannot parse column: %q", line)
	}

	name := m[1]
	rawType := strings.TrimSpace(m[2])
	rest := strings.ToUpper(m[3])

	goType := sqlToGoType(rawType)
	col := Column{
		Name:     name,
		SQLType:  normalizeType(rawType),
		GoType:   goType,
		IsPK:     strings.Contains(rest, "PRIMARY KEY"),
		NotNull:  strings.Contains(rest, "NOT NULL") || strings.Contains(rest, "PRIMARY KEY"),
		IsUnique: strings.Contains(rest, "UNIQUE"),
	}
	return col, nil
}

// splitColumnDefs splits the table body on commas, respecting parentheses depth.
func splitColumnDefs(body string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, ch := range body {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, body[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, body[start:])
	return parts
}

func extractParens(s string) string {
	open := strings.Index(s, "(")
	if open == -1 {
		return ""
	}
	body, _, ok := extractParenBody(s, open+1)
	if !ok {
		return ""
	}
	return body
}

// HasTimeColumn returns true if any column maps to time.Time.
func (t Table) HasTimeColumn() bool {
	for _, c := range t.Columns {
		if c.GoType == "time.Time" {
			return true
		}
	}
	return false
}

func normalizeType(rawType string) string {
	return strings.ToUpper(strings.TrimSpace(regexp.MustCompile(`\s*\([^)]*\)`).ReplaceAllString(rawType, "")))
}

func sqlToGoType(rawType string) string {
	t := strings.ToUpper(normalizeType(rawType))
	switch {
	case t == "INTEGER" || t == "INT" || t == "BIGINT" || t == "SMALLINT" || t == "TINYINT" || t == "INT2" || t == "INT8":
		return "int64"
	case t == "REAL" || t == "FLOAT" || t == "DOUBLE" || t == "NUMERIC" || t == "DECIMAL":
		return "float64"
	case t == "BOOLEAN" || t == "BOOL":
		return "bool"
	case t == "DATETIME" || t == "TIMESTAMP" || t == "DATE" || t == "TIME":
		return "time.Time"
	case t == "BLOB":
		return "[]byte"
	default: // TEXT, VARCHAR, CHAR, CLOB, and anything else
		return "string"
	}
}

// toPascal converts snake_case or lowercase to PascalCase.
func toPascal(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return b.String()
}

// singular is a minimal singularization for common English patterns.
func singular(s string) string {
	switch {
	case strings.HasSuffix(s, "ies"):
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "zes"):
		return s[:len(s)-2]
	case strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss"):
		return s[:len(s)-1]
	default:
		return s
	}
}
