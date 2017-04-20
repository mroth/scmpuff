Aruba.configure do |config|
  # root directory for aruba tests
  # default is ".", but...
  # we want to move this OUTSIDE of the project directory though, since aruba
  # isolation is insufficient otherwise for git not to see the parent repo.
  #
  # TODO: make this OS sensitive
  config.root_directory = '/tmp/'

  # the working directory inside the root directory note that our hacky trick to
  # speed up mocks walks up one level from here and makes a sibling directory
  # (see "mocked git repository with commited subdirectory and file" in
  # scmpuff_steps.rb for details)
  #
  # default is "tmp/aruba", but that is confusing when nested inside '/tmp'
  config.working_directory = 'aruba/workdir'
end
