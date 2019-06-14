package internal

import (
	"fmt"
	"strings"
)

func runWhichLanguage() {
	backend := getBackend("")
	fmt.Println(backend.name)
}

func runListLanguages() {
	for _, backendName := range getBackendNames() {
		fmt.Println(backendName)
	}
}

func runSearch(language string, queries []string, outputFormat outputFormat) {
	results := getBackend(language).search(queries)
	fmt.Printf("output %#v in format %#v\n", results, outputFormat)
	notImplemented()
}

func runInfo(language string, pkg string, outputFormat outputFormat) {
	info := getBackend(language).info(pkgName(pkg))
	fmt.Printf("output %#v in format %#v\n", info, outputFormat)
	notImplemented()
}

func runAdd(language string, args []string, guess bool) {
	pkgs := map[pkgName]pkgSpec{}
	for _, arg := range args {
		fields := strings.SplitN(arg, " ", 2)
		name := pkgName(fields[0])
		var spec pkgSpec
		if len(fields) >= 2 {
			spec = pkgSpec(fields[1])
		}

		pkgs[name] = spec
	}

	backend := getBackend(language)

	if guess {
		for name, _ := range backend.guess() {
			if _, ok := pkgs[name]; !ok {
				pkgs[name] = ""
			}
		}
	}

	for name, _ := range backend.listSpecfile() {
		delete(pkgs, name)
	}

	if len(pkgs) >= 1 {
		backend.add(pkgs)
	}

	store := readStore()

	if !doesStoreSpecfileHashMatch(store, backend.specfile) {
		backend.lock()
	}

	if !doesStoreLockfileHashMatch(store, backend.lockfile) {
		backend.install()
	}

	updateStoreHashes(store, backend.specfile, backend.lockfile)
}

func runRemove(language string, args []string) {
	backend := getBackend(language)
	specfilePkgs := backend.listSpecfile()

	pkgs := map[pkgName]bool{}
	for _, arg := range args {
		name := pkgName(arg)
		if _, ok := specfilePkgs[name]; ok {
			pkgs[name] = true
		}
	}

	if len(pkgs) >= 1 {
		backend.remove(pkgs)
	}

	store := readStore()

	if !doesStoreSpecfileHashMatch(store, backend.specfile) {
		backend.lock()
	}

	if !doesStoreLockfileHashMatch(store, backend.lockfile) {
		backend.install()
	}

	updateStoreHashes(store, backend.specfile, backend.lockfile)
}

func runLock(language string, force bool) {
	backend := getBackend(language)
	store := readStore()
	if doesStoreSpecfileHashMatch(store, backend.specfile) && !force {
		return
	}
	backend.lock()

	if !doesStoreLockfileHashMatch(store, backend.lockfile) {
		backend.install()
	}

	updateStoreHashes(store, backend.specfile, backend.lockfile)
}

func runInstall(language string, force bool) {
	backend := getBackend(language)
	store := readStore()
	if doesStoreLockfileHashMatch(store, backend.lockfile) && !force {
		return
	}
	backend.install()
	updateStoreHashes(store, backend.specfile, backend.lockfile)
}

func runList(language string, all bool, outputFormat outputFormat) {
	backend := getBackend(language)
	if all {
		results := backend.listLockfile()
		fmt.Printf("output %#v in format %#v\n", results, outputFormat)
		notImplemented()
	} else {
		results := backend.listSpecfile()
		fmt.Printf("output %#v in format %#v\n", results, outputFormat)
		notImplemented()
	}
}

func runGuess(language string, all bool) {
	backend := getBackend(language)
	pkgs := backend.guess()

	if (!all) {
		for name, _ := range backend.listSpecfile() {
			delete(pkgs, name)
		}
	}

	for name, _ := range pkgs {
		fmt.Println(name)
	}
}
