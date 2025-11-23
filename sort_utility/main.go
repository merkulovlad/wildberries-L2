package main

import (
	"bufio"
	"fmt"
	"github.com/merkulovlad/wildberries-L2/sort_utility/internal"
	"github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	opts := internal.ParseFlags()

	files := pflag.Args()
	records := make([]*internal.Record, 0, 500)

	if len(files) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			records = append(records, internal.NewRecord(scanner.Text(), opts))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, file := range files {
			// #nosec G304 -- читаем файлы, указанные пользователем в CLI;
			f, err := os.Open(file)
			if err != nil {
				log.Fatalf("Failed to open file %s: %v", file, err)
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				records = append(records, internal.NewRecord(scanner.Text(), opts))
			}

			if err := scanner.Err(); err != nil {
				log.Fatalf("Error reading file %s: %v", file, err)
			}

			_ = f.Close()
		}
	}

	if opts.CheckOnly {
		if internal.CheckIsSorted(records, opts) {
			fmt.Println("The input is sorted.")
			os.Exit(0)
		} else {
			fmt.Println("The input is not sorted.")
			os.Exit(1)
		}
	} else {
		internal.SortRecords(records, opts)

		for _, record := range records {
			fmt.Println(record.OriginalLine)
		}
	}
}
