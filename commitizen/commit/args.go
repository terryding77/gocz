package commit

type Args struct {
	//    usage: git commit [<options>] [--] <pathspec>...
	//
	//    -q, --quiet           suppress summary after successful commit
	Quiet bool
	//    -v, --verbose         show diff in commit message template
	Verbose bool
	//
	//Commit message options
	//    -F, --file <file>     read message from file
	//    --author <author>     override author for commit
	Author string
	//    --date <date>         override date for commit
	Date string
	//    -m, --message <message>
	//                          commit message
	//    -c, --reedit-message <commit>
	//                          reuse and edit message from specified commit
	//    -C, --reuse-message <commit>
	//                          reuse message from specified commit
	//    --fixup <commit>      use autosquash formatted message to fixup specified commit
	//    --squash <commit>     use autosquash formatted message to squash specified commit
	//    --reset-author        the commit is authored by me now (used with -C/-c/--amend)
	//    -s, --signoff         add Signed-off-by:
	//    -t, --template <file>
	//                          use specified template file
	//    -e, --edit            force edit of commit
	//    --cleanup <default>   how to strip spaces and #comments from message
	//    --status              include status in commit message template
	//    -S, --gpg-sign[=<key-id>]
	//                          GPG sign commit
	//
	//Commit contents options
	//    -a, --all             commit all changed files
	All bool
	//    -i, --include         add specified files to index for commit
	//    --interactive         interactively add files
	//    -p, --patch           interactively add changes
	//    -o, --only            commit only specified files
	//    -n, --no-verify       bypass pre-commit and commit-msg hooks
	//    --dry-run             show what would be committed
	//    --short               show status concisely
	//    --branch              show branch information
	//    --ahead-behind        compute full ahead/behind values
	//    --porcelain           machine-readable output
	//    --long                show status in long format (default)
	//    -z, --null            terminate entries with NUL
	//    --amend               amend previous commit
	Amend bool
	//    --no-post-rewrite     bypass post-rewrite hook
	//    -u, --untracked-files[=<mode>]
	//                          show untracked files, optional modes: all, normal, no. (Default: all)
}

func NewDefaultArgs() *Args {
	return &Args{
		Quiet: false,
	}
}

func (args Args) Combination(commitMsg string) (rawArgs []string) {
	rawArgs = append(rawArgs, "commit")
	if args.Quiet {
		rawArgs = append(rawArgs, "-q")
	}
	if args.Verbose {
		rawArgs = append(rawArgs, "-v")
	}
	if args.Author != "" {
		rawArgs = append(rawArgs, "--author", args.Author)
	}
	if args.Date != "" {
		rawArgs = append(rawArgs, "--date", args.Date)
	}
	rawArgs = append(rawArgs, "--message", commitMsg)
	if args.All {
		rawArgs = append(rawArgs, "-a")
	}
	if args.Amend {
		rawArgs = append(rawArgs, "--amend")
	}

	return rawArgs
}
