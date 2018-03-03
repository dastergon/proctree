# proctree

proctree is a tool to display a tree of running processes. The tree is rooted at either pid or PID 1 if pid argument is omitted.

proctree is influenced by the [pstree](www.thp.uni-duisburg.de/pstree) command line utility.

## Installation

    go get -u github.com/dastergon/proctree

## Example

To display a process tree starting from PID 79408:

    proctree 79408

The output looks like the following:

```
79408 dastergon (/Applications/iTerm.app/Contents/MacOS/iTerm2)
└── 79412 root (/usr/bin/login)
    └── 79413 dastergon (-zsh)
        └── 81433 dastergon (vim)
```

## Usage

```
Usage: ./proctree [options...] [<pid>] (defaults to PID 1)
  -U	Do not show branches containing only root processes
  -g	Show process group ids
  -l int
    	Print tree to n level deep (default -1)
  -p int
    	Show only branches containing process <pid> (default -1)
  -u string
    	Show only branches containing process of <user>
  -version
    	Outputs the version of proctree.
```
