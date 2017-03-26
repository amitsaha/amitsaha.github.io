:Title: Notes on using Golang to write gitbackup
:Date: 2017-03-26 10:00
:Category: golang

`gitbackup <https://github.com/amitsaha/gitbackup>`__ is a tool to backup your git repositories from GitHub and GitLab. I wrote the `initial version <https://github.com/amitsaha/gitbackup/releases/tag/lj-0.1>`__ as a project for a golang article which is in review for publication in a Linux magazine.

Since the initial version, the project's code ha seen number of changes which has been a learning experience for me since I am still fairly new to Golang. In the rest of this post, I describe these and some personal notes about the experience.

Using afero for filesystem operations
=====================================

``gitbackup`` needs to do some basic filesystem operations - create directories and check for existence of directories. In the initial version, I was using the ``os`` package directly which meant any test invoking the code which performed these operations were actually performing those on the underlying filesystem. I could of course
perform cleanup after these tests so that my filesystem would not remain polluted. However, then I decided to check what `afero <https://github.com/spf13/afero>`__ had to offer. It had exactly what I needed - a memory backed filesystem (`NewMemMapFs`).This `section <https://github.com/spf13/afero#using-afero-for-testing>`__ in the project homepage was all I needed to switch to using `afero` instead of `os` package drirectly. And hence I didn't need to worry about cleaning up my filesystem after a test run or worry about starting from a known clean state!

To show some code, this is the `git diff` of introducing `afero` and switching out direct use of `os`:

.. code:: diff

    diff --git a/src/gitbackup/main.go b/src/gitbackup/main.go
    index 500d9a2..6e71beb 100644
    --- a/src/gitbackup/main.go
    +++ b/src/gitbackup/main.go
    @@ -3,6 +3,7 @@ package main
     import (
            "flag"
            "github.com/mitchellh/go-homedir"
    +       "github.com/spf13/afero"
            "log"
            "os"
            "os/exec"
    @@ -14,6 +15,7 @@ import (
     var MAX_CONCURRENT_CLONES int = 20

     var execCommand = exec.Command
    +var appFS = afero.NewOsFs()
     var gitCommand = "git"

     // Check if we have a copy of the repo already, if
    @@ -22,7 +24,7 @@ func backUp(backupDir string, repo *Repository, wg *sync.WaitGroup) ([]byte, err
            defer wg.Done()

            repoDir := path.Join(backupDir, repo.Name)
    -       _, err := os.Stat(repoDir)
    +       _, err := appFS.Stat(repoDir)

            var stdoutStderr []byte
            if err == nil {
    @@ -83,7 +85,7 @@ func main() {
            } else {
                    *backupDir = path.Join(*backupDir, *service)
            }
    -       _, err := os.Stat(*backupDir)
    +       _, err := appFS.Stat(*backupDir)
            if err != nil {
                    log.Printf("%s doesn't exist, creating it\n", *backupDir)
                    err := os.MkdirAll(*backupDir, 0771)

When we declare `appFS` above outside all functions, it becomes a package level
variable and we set it to `NewOsFs()` and replace function calls such as `os.Stat` by `appFS.Stat()`. Since the variable name starts with a small letter, this variable is not visible outside the package.

Then, in the test, I will do:

.. code::

    appFS = afero.NewMemMapFs()

Hence, all operations will happen in the memory based filesystem rather than the "real" underlying filesystem.

Testing shell commands
======================

One of the first roadblocks to writing tests I faced was how to test functions which were invoking external programs (``git`` in this case). This post here titled `Testing os/exec.Command <https://npf.io/2015/06/testing-exec-command/>`__ had my answer. However, it took me a while to correctly apply it. And that post is still the reference if you want to understand what's going on.

Here's basically what I did:

.. code:: golang

    var execCommand = exec.Command
    ..

    func backUp(backupDir string, repo *Repository, wg *sync.WaitGroup) ([]byte, error) {
        ...
        if err == nil {
            ..
            cmd := execCommand(gitCommand, "-C", repoDir, "pull")
            ..
        } else {
            ..
            cmd := execCommand(gitCommand, "clone", repo.GitURL, repoDir)
            ..
        }
        ...
    }

We declare a package variable, ``execCommand`` which is intialized with ``exec.Command`` from the ``os/exec`` package. Then, in the tests, I do the following:

.. code:: golang

    func TestHelperCloneProcess(t *testing.T) {
        if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
            return
        }
        // Check that git command was executed
        if os.Args[3] != "git" || os.Args[4] != "clone" {
            fmt.Fprintf(os.Stdout, "Expected git clone to be executed. Got %v", os.Args[3:])
            os.Exit(1)
        }
        os.Exit(0)
    }


    func fakeCloneCommand(command string, args ...string) (cmd *exec.Cmd) {
        cs := []string{"-test.run=TestHelperCloneProcess", "--", command}
        cs = append(cs, args...)
        cmd = exec.Command(os.Args[0], cs...)
        cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
        return cmd
    }

    execCommand = fakeCloneCommand
    stdoutStderr, err := backUp(backupDir, &repo, &wg)

Switching from ``gb`` to standard go tooling
============================================





