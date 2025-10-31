package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func findGitRepo() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Достигли корня файловой системы
			return "", fmt.Errorf("git repository not found")
		}
		dir = parent
	}
}

func main() {
	fmt.Println("Generating version...")

	repoPath, err := findGitRepo()
	if err != nil {
		log.Fatal("Git repository not found:", err)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	refIter, err := repo.Tags()
	if err != nil {
		log.Fatal(err)
	}

	var latestTag *plumbing.Reference
	var latestTime time.Time

	err = refIter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}

		obj, err := repo.Object(plumbing.AnyObject, ref.Hash())
		if err != nil {
			log.Printf("Не удалось получить объект для %s: %v", ref.Name(), err)
			return nil
		}

		var commitHash plumbing.Hash

		switch obj := obj.(type) {
		case *object.Tag:
			commitHash = obj.Target
		case *object.Commit:
			commitHash = ref.Hash()
		default:
			log.Printf("Неподдерживаемый тип объекта для тега %s: %T", ref.Name(), obj)
			return nil
		}

		commit, err := repo.CommitObject(commitHash)
		if err != nil {
			log.Printf("Не удалось получить коммит для %s: %v", ref.Name(), err)
			return nil
		}

		if commit.Author.When.After(latestTime) {
			latestTime = commit.Author.When
			latestTag = ref
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	versionPath := filepath.Join(repoPath, "cmd/VERSION")
	if latestTag != nil {
		err = os.WriteFile(versionPath, []byte(latestTag.Name().Short()), 0o644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Version written: %s\n", latestTag.Name().Short())
	} else {
		err = os.WriteFile(versionPath, []byte("v0.0.0"), 0o644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Version written: v0.0.0 (no tags found)")
	}
}
