package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`^INSERT INTO (\w+) `)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		log.Fatal("directory not specified")
	}

	// ディレクトリの存在チェック
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	cnf, err := readConfig("testsql-sorter.yml")
	if err != nil {
		log.Fatal(err)
	}

	// 指定ディレクトリのファイルを再帰的に検索
	files := []string{}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(filepath.Base(path), "Test") || filepath.Ext(path) != ".sql" {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		sqls, err := readTestSQL(f)
		if err != nil {
			log.Fatalf("%s %s", f, err)
		}

		err = checkTestSQL(cnf.Tables, sqls)
		if err != nil {
			log.Fatalf("%s %s", f, err)
		}

		sqls, err = sortTestSQL(cnf.Tables, sqls)
		if err != nil {
			log.Fatalf("%s %s", f, err)
		}

		err = writeTestSQL(f, sqls)
		if err != nil {
			log.Fatalf("%s %s", f, err)
		}
	}
}

// SQLファイルのINSERT文を抽出する
func readTestSQL(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sqls := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "INSERT") {
			sqls = append(sqls, line)
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return sqls, nil
}

// Configのテーブル一覧に不足がないかチェックする
func checkTestSQL(tables []string, sqls []string) error {
	m := map[string]bool{}

	for _, s := range sqls {
		matches := reg.FindStringSubmatch(s)
		if matches == nil {
			return errors.New("regexp does not match")
		}
		table := matches[1]
		m[table] = false
	}

	for _, t := range tables {
		if _, ok := m[t]; ok {
			m[t] = true
		}
	}

	missing := []string{}
	for k, v := range m {
		if !v {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		msg := fmt.Sprintf("tables do not exist in config (%s)", strings.Join(missing, ", "))
		return errors.New(msg)
	}

	return nil
}

// Configに準じてINSERT文を並べ替える
func sortTestSQL(tables []string, sqls []string) ([]string, error) {
	sorted := []string{}

	for _, t := range tables {
		prefix := fmt.Sprintf("INSERT INTO %s ", t)

		for _, s := range sqls {
			if strings.HasPrefix(s, prefix) {
				sorted = append(sorted, s)
			}
		}
	}

	if len(sqls) != len(sorted) {
		return nil, errors.New("failed to sort sql")
	}

	return sorted, nil
}

// SQLファイルにINSERT文を上書きする
func writeTestSQL(path string, sqls []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, s := range sqls {
		_, err = file.WriteString(s + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
