package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/go-semver/semver"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
)

func luaScriptLatestVersionFromTags(tags storer.ReferenceIter) *semver.Version {
	lVersion := semver.New("0.0.1")
	tags.ForEach(func(ref *plumbing.Reference) error {
		thisVerStr := ref.Name().Short()
		if thisVerStr[0] != 'v' {
			return nil
		}
		thisVerStr = strings.TrimPrefix(thisVerStr, "v")
		thisVer, err := semver.NewVersion(thisVerStr)
		if err != nil {
			log.Printf("[WARN] %s", err.Error())
			return err
		}
		if lVersion.LessThan(*thisVer) {
			lVersion = thisVer
		}
		return nil
	})
	return lVersion
}

func luaScriptLastestVersion(url string) (*semver.Version, error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:   url,
		Depth: 1,
		Tags:  git.AllTags,
	})
	if err != nil {
		return nil, err
	}
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}
	return luaScriptLatestVersionFromTags(tags), nil
}

func luaInstallScript(url string, version semver.Version) error {
	nameSplit := strings.Split(filepath.Base(url), ".")
	name := nameSplit[0]
	logInfo("Install/update script from '%s' (v%s).", url, version.String())
	pathTo := filepath.Join(getScriptPath(), name)
	os.RemoveAll(pathTo)
	_, err := git.PlainClone(
		pathTo,
		false,
		&git.CloneOptions{
			URL:           url,
			ReferenceName: plumbing.NewTagReferenceName("v" + version.String()),
			Depth:         1,
		},
	)
	return err
}

func (ls *luaScript) getGitRepository() (*git.Repository, error) {
	repo, err := git.PlainOpen(filepath.Dir(ls.Path))
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return nil, ErrNoGit
	}
	return repo, err
}

func (ls *luaScript) getLatestVersion() (*semver.Version, error) {
	repo, err := ls.getGitRepository()
	if err != nil {
		return nil, err
	}
	repo.Fetch(&git.FetchOptions{
		Tags: git.AllTags,
	})
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	return luaScriptLatestVersionFromTags(tags), nil
}

func (ls *luaScript) getCurrentVersion() (*semver.Version, error) {
	repo, err := ls.getGitRepository()
	if err != nil {
		return nil, err
	}
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}
	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		logWarn(err.Error())
		return nil, err
	}
	var ver *semver.Version
	tags.ForEach(func(ref *plumbing.Reference) error {
		commitHash, _ := repo.ResolveRevision(plumbing.Revision(ref.Name().String()))
		if *commitHash == headCommit.Hash {
			thisVerStr := ref.Name().Short()
			if thisVerStr[0] != 'v' {
				return nil
			}
			thisVerStr = strings.TrimPrefix(thisVerStr, "v")
			ver, err = semver.NewVersion(thisVerStr)
			return err
		}
		return nil
	})
	return ver, nil
}

func (ls *luaScript) update() error {
	repo, err := ls.getGitRepository()
	if err != nil {
		return err
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			return ErrGitNoRemote
		}
		return err
	}
	url := remote.Config().URLs[0]
	if url == "" {
		return ErrGitNoRemote
	}
	latestVersion, err := ls.getLatestVersion()
	if err != nil {
		return err
	}
	return luaInstallScript(url, *latestVersion)
}
