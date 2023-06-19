package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/graph-guard/gqlscan"
)

func main() {
	var fInDirPath string
	var fSkipErr bool
	var fParallel int
	var fDur time.Duration
	var fOrder string
	flag.StringVar(
		&fInDirPath,
		"i",
		"./inputs",
		"path to inputs directory containing `.graphql` files.",
	)
	flag.BoolVar(
		&fSkipErr,
		"skiperr",
		false,
		"don't stop if an error is encountered.",
	)
	flag.IntVar(
		&fParallel,
		"p",
		runtime.NumCPU(),
		"number of workers to run in parallel.",
	)
	flag.DurationVar(
		&fDur,
		"d",
		5*time.Second,
		"test duration.",
	)
	flag.StringVar(
		&fOrder,
		"o",
		"rnd",
		"defines the execution order "+
			"(`rnd` for random, `seq` for sequential)",
	)
	flag.Parse()

	var inputs []InputFile
	if err := find(fInDirPath, ".graphql", func(path string) {
		c, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("reading file %q: %v", path, err)
		}
		inputs = append(inputs, InputFile{
			Path:     path,
			Name:     filepath.Base(path),
			Contents: c,
		})
	}); err != nil {
		log.Fatalf("reading input directory: %s", err)
	}

	if len(inputs) < 1 {
		log.Fatal("no input files")
	}

	fmt.Printf("gathered %d input file(s):\n", len(inputs))
	errAt := -1
	for i, f := range inputs {
		fmt.Printf(
			" %s: %s\n",
			f.Path, humanize.Bytes(uint64(len(f.Contents))),
		)
		err := gqlscan.ScanAll(f.Contents, func(i *gqlscan.Iterator) {})
		if err.IsErr() {
			fmt.Printf("  %s\n", err.Error())
			errAt = i
		}
	}
	if !fSkipErr && errAt != -1 {
		log.Fatalf("use -skiperr to ignore error in input files\n")
	}

	var wg sync.WaitGroup
	wg.Add(fParallel)
	var totalBytes uint64
	var totalErrs uint64
	var timeStart time.Time

	fn, ok := modes[fOrder]
	if !ok {
		log.Fatalf("unknown order: %q", fOrder)
	}
	fmt.Printf("running %d parallel goroutine(s) for %s\n", fParallel, fDur)

	timeStart = time.Now()
	for i := 0; i < fParallel; i++ {
		go fn(fDur, inputs, timeStart, &wg, &totalBytes, &totalErrs)
	}

	wg.Wait()
	timeEnd := time.Now()
	totalDur := timeEnd.Sub(timeStart)
	tb, te := atomic.LoadUint64(&totalBytes), atomic.LoadUint64(&totalErrs)
	bytesPerSec := math.Floor(float64(tb) / totalDur.Seconds())
	bps := fmt.Sprintf("%.2f gbit/s", gbits(float64(tb), totalDur.Seconds()))

	fmt.Printf(
		"finished (%s)\ntotal processed: %s (%s/s; %s)\ntotal errors: %s\n",
		totalDur.String(),
		humanize.Bytes(tb),
		humanize.Bytes(uint64(bytesPerSec)),
		bps,
		humanize.Comma(int64(te)),
	)
}

func gbits(totalBytes, seconds float64) float64 {
	return ((totalBytes * 8) / 1e9) / seconds
}

func find(root, ext string, fn func(path string)) error {
	return filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			fn(s)
		}
		return nil
	})
}

var modes = map[string]func(
	testDur time.Duration,
	inputs []InputFile,
	timeStart time.Time,
	wg *sync.WaitGroup,
	totalBytes, totalErrs *uint64,
){
	"rnd": func(
		testDur time.Duration,
		inputs []InputFile,
		timeStart time.Time,
		wg *sync.WaitGroup,
		totalBytes, totalErrs *uint64,
	) {
		rnd := NewRandSrc(time.Now().Unix())
		var tb, te uint64
		defer wg.Done()
		for time.Since(timeStart) < testDur {
			in := inputs[rnd.IntRange(0, len(inputs))]
			err := gqlscan.ScanAll(in.Contents, func(i *gqlscan.Iterator) {})
			if err.IsErr() {
				te++
			}
			tb += uint64(len(in.Contents))
		}
		atomic.AddUint64(totalBytes, tb)
		atomic.AddUint64(totalErrs, te)
	},
	"seq": func(
		testDur time.Duration,
		inputs []InputFile,
		timeStart time.Time,
		wg *sync.WaitGroup,
		totalBytes, totalErrs *uint64,
	) {
		var tb, te uint64
		defer wg.Done()
		var i int
		for time.Since(timeStart) < testDur {
			in := inputs[i]
			err := gqlscan.ScanAll(in.Contents, func(i *gqlscan.Iterator) {})
			if err.IsErr() {
				te++
			}
			tb += uint64(len(in.Contents))
			i++
			if i >= len(inputs) {
				i = 0
			}
		}
		atomic.AddUint64(totalBytes, tb)
		atomic.AddUint64(totalErrs, te)
	},
}

type InputFile struct {
	Name     string
	Path     string
	Contents []byte
}

type RandSrc struct{ s *rand.Rand }

func NewRandSrc(seed int64) *RandSrc {
	return &RandSrc{s: rand.New(rand.NewSource(seed))}
}

func (r *RandSrc) IntRange(greaterThanOrEq, lessThan int) int {
	return r.s.Intn(lessThan-greaterThanOrEq) + greaterThanOrEq
}
