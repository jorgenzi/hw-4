package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Ошибки
var (
	ErrMutuallyExclusive = errors.New("параметры -c, -d, -u взаимозаменяемы")
	ErrNegativeNumber    = errors.New("числовые параметры не могут быть отрицательными")
)

// Options представляет параметры для уникализации
type Options struct {
	Count      bool // -c
	Duplicate  bool // -d
	Unique     bool // -u
	IgnoreCase bool // -i
	NumFields  int  // -f
	NumChars   int  // -s
}

// Validate проверяет корректность опций
func (o *Options) Validate() error {
	// Проверяем взаимоисключающие параметры
	count := 0
	if o.Count {
		count++
	}
	if o.Duplicate {
		count++
	}
	if o.Unique {
		count++
	}

	if count > 1 {
		return ErrMutuallyExclusive
	}

	// Проверяем числовые параметры
	if o.NumFields < 0 {
		return ErrNegativeNumber
	}
	if o.NumChars < 0 {
		return ErrNegativeNumber
	}

	return nil
}

// LineInfo хранит информацию о строке
type LineInfo struct {
	Original  string
	Processed string
	Count     int
	Index     int
}

// Uniq выполняет уникализацию строк согласно опциям
func Uniq(lines []string, opts Options) ([]string, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// Обрабатываем строки
	processed := processLines(lines, opts)

	// Группируем строки
	groups := groupLines(processed)

	// Формируем результат согласно опциям
	return formatResult(groups, opts), nil
}

// processLines обрабатывает строки согласно опциям -f, -s, -i
func processLines(lines []string, opts Options) []LineInfo {
	result := make([]LineInfo, len(lines))

	for i, line := range lines {
		processed := line

		// Применяем -f (игнорируем первые num_fields полей)
		if opts.NumFields > 0 {
			processed = skipFields(processed, opts.NumFields)
		}

		// Применяем -s (игнорируем первые num_chars символов)
		if opts.NumChars > 0 {
			if len(processed) > opts.NumChars {
				processed = processed[opts.NumChars:]
			} else {
				processed = ""
			}
		}

		// Применяем -i (игнорируем регистр)
		if opts.IgnoreCase {
			processed = strings.ToLower(processed)
		}

		result[i] = LineInfo{
			Original:  line,
			Processed: processed,
			Index:     i,
		}
	}

	return result
}

// skipFields пропускает первые n полей в строке
func skipFields(line string, n int) string {
	fields := strings.Fields(line)
	if n >= len(fields) {
		return ""
	}
	return strings.Join(fields[n:], " ")
}

// groupLines группирует одинаковые строки
func groupLines(lines []LineInfo) map[string][]LineInfo {
	groups := make(map[string][]LineInfo)

	for _, line := range lines {
		key := line.Processed
		groups[key] = append(groups[key], line)
	}

	return groups
}

// formatResult форматирует результат согласно опциям

func formatResult(groups map[string][]LineInfo, opts Options) []string {
    // Если нет групп, возвращаем пустой слайс (не nil)
    if len(groups) == 0 {
        return []string{}
    }

    var result []string

    // Собираем все группы для сортировки по первоначальному порядку
    sortedGroups := getSortedGroups(groups)

    for _, group := range sortedGroups {
        result = appendGroupResult(result, group, opts)
    }

    return result
}

// getSortedGroups возвращает группы, отсортированные по индексу первой строки
func getSortedGroups(groups map[string][]LineInfo) [][]LineInfo {
    var sortedGroups [][]LineInfo
    for _, group := range groups {
        sortedGroups = append(sortedGroups, group)
    }

    // Сортируем группы по индексу первой строки
    sort.Slice(sortedGroups, func(i, j int) bool {
        return sortedGroups[i][0].Index < sortedGroups[j][0].Index
    })

    return sortedGroups
}

// appendGroupResult добавляет результат для одной группы согласно опциям
func appendGroupResult(result []string, group []LineInfo, opts Options) []string {
    count := len(group)
    originalLine := group[0].Original

    switch {
    case opts.Count:
        // -c: выводим количество и строку
        result = append(result, formatCountLine(count, originalLine))

    case opts.Duplicate:
        // -d: выводим только дубликаты
        if count > 1 {
            result = append(result, originalLine)
        }

    case opts.Unique:
        // -u: выводим только уникальные строки
        if count == 1 {
            result = append(result, originalLine)
        }

    default:
        // Без параметров: выводим только первую строку из группы
        result = append(result, originalLine)
    }

    return result
}

// formatCountLine форматирует строку с количеством
func formatCountLine(count int, line string) string {
	return formatCount(count) + " " + line
}

// formatCount форматирует число для вывода
func formatCount(count int) string {
	if count == 0 {
		return "0"
	}

	var buf [20]byte
	i := len(buf)

	for count > 0 {
		i--
		buf[i] = byte('0' + count%10)
		count /= 10
	}

	return string(buf[i:])
}

// ReadLines читает строки из io.Reader
func ReadLines(reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// WriteLines записывает строки в io.Writer
func WriteLines(writer io.Writer, lines []string) error {
	for _, line := range lines {
		if _, err := io.WriteString(writer, line+"\n"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
    opts, inputFile, outputFile, err := parseFlags()
    if err != nil {
        showUsage(err)
        os.Exit(1)
    }

    input := getInputReader(inputFile)
    output := getOutputWriter(outputFile)

    processUniq(input, output, opts)
}

// showUsage показывает сообщение об ошибке и использование
func showUsage(err error) {
    fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
    fmt.Fprintf(os.Stderr, "Использование: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]\n")
}

// getInputReader возвращает reader для ввода
func getInputReader(inputFile string) io.Reader {
    if inputFile == "" {
        return os.Stdin
    }

    file, err := os.Open(inputFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка открытия файла %s: %v\n", inputFile, err)
        os.Exit(1)
    }
    return file
}

// getOutputWriter возвращает writer для вывода
func getOutputWriter(outputFile string) io.Writer {
    if outputFile == "" {
        return os.Stdout
    }

    file, err := os.Create(outputFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка создания файла %s: %v\n", outputFile, err)
        os.Exit(1)
    }
    return file
}

// processUniq обрабатывает уникализацию
func processUniq(input io.Reader, output io.Writer, opts Options) {
    // Читаем строки
    lines, err := ReadLines(input)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
        os.Exit(1)
    }

    // Выполняем уникализацию
    result, err := Uniq(lines, opts)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка обработки: %v\n", err)
        os.Exit(1)
    }

    // Записываем результат
    if err := WriteLines(output, result); err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка записи: %v\n", err)
        os.Exit(1)
    }
}

// parseFlags разбирает аргументы командной строки
func parseFlags() (Options, string, string, error) {
    var opts Options

    flag.BoolVar(&opts.Count, "c", false, "подсчитать количество встречаний строки")
    flag.BoolVar(&opts.Duplicate, "d", false, "вывести только повторяющиеся строки")
    flag.BoolVar(&opts.Unique, "u", false, "вывести только уникальные строки")
    flag.BoolVar(&opts.IgnoreCase, "i", false, "игнорировать регистр")
    flag.IntVar(&opts.NumFields, "f", 0, "игнорировать первые num_fields полей")
    flag.IntVar(&opts.NumChars, "s", 0, "игнорировать первые num_chars символов")

    // Используем custom Usage для избежания вывода справки при ошибках
    flag.Usage = func() {
        // Пустая функция, чтобы не выводить справку автоматически
    }

    flag.Parse()

    // Валидируем опции ДО возврата из функции
    if err := opts.Validate(); err != nil {
        return Options{}, "", "", err
    }

    // Обрабатываем позиционные аргументы
    args := flag.Args()
    var inputFile, outputFile string

    switch len(args) {
    case 0:
        // нет файлов
    case 1:
        inputFile = args[0]
    case 2:
        inputFile = args[0]
        outputFile = args[1]
    default:
        return Options{}, "", "", fmt.Errorf("слишком много аргументов")
    }

    return opts, inputFile, outputFile, nil
}
