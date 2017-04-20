Aruba.configure do |config|
  # root directory for aruba tests
  # default is ".", but...
  # we want to move this OUTSIDE of the project directory though, since aruba
  # isolation is insufficient otherwise for git not to see the parent repo.
  if (RUBY_PLATFORM =~ /darwin/)
    # on macOS, /tmp and /var/folders are both actually symlinked into /private,
    # which causes some reliability issues because Golang's os.Getwd does not
    # consistently return a path, and it seems to like to return /private/tmp
    # etc, whereas most ways the user might select the TMPDIR omit it, which
    # makes our test matching more difficult.
    #
    # related, see also https://github.com/mroth/scmpuff/issues/11
    #
    # To get around this in tests, for now we just manually specify /private/tmp
    # on macOS, since current versions of macOS/Golang seem to agree on.
    config.root_directory = '/private/tmp'
  else
    config.root_directory = Dir.tmpdir
  end

  # the working directory inside the root directory note that our hacky trick to
  # speed up mocks walks up one level from here and makes a sibling directory
  # (see "mocked git repository with commited subdirectory and file" in
  # scmpuff_steps.rb for details)
  #
  # default is "tmp/aruba", but that is confusing when nested inside '/tmp'
  config.working_directory = 'aruba/workdir'
end
