package artisan_core

type PackageSet = map[string]bool

//noinspection GoUnusedExportedFunction
func PackageSetClone(m PackageSet) (n PackageSet) {
	if m == nil {
		return nil
	}
	n = make(PackageSet)
	for k, v := range m {
		n[k] = v
	}
	return n
}

//noinspection GoUnusedExportedFunction
func PackageSetMerge(pac PackageSet, oth PackageSet) PackageSet {
	newPac := make(PackageSet)
	for k, v := range pac {
		newPac[k] = v
	}
	for k, v := range oth {
		newPac[k] = v
	}
	return newPac
}

func PackageSetInPlaceMerge(pac PackageSet, oth PackageSet) PackageSet {
	if pac == nil {
		pac = make(PackageSet)
	}
	for k, v := range oth {
		pac[k] = v
	}
	return pac
}

func PackageSetAppend(packageSet PackageSet, pkg string) PackageSet {
	if packageSet == nil {
		packageSet = make(map[string]bool)
	}

	if len(pkg) != 0 {
		packageSet[pkg] = true
	}
	return packageSet
}
