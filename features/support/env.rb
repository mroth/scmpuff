require 'aruba/cucumber'

Before do
  # sigh, apparently overriding the tmp dir is buggy, still gets created empty,
  # and it doesnt properly clobber before scenarios are run.
  #
  # I'd like to keep it out of here so Dropbox doesnt mess with it, but
  # such is life I suppose.
  # @dirs = ["/tmp/aruba/scmpuff"]

  # override PATH (find binary)

  # override HOME (dont want to read ~/.gitconfig)
  # set_env 'HOME', File.expand_path(File.join(current_dir, 'home'))
  # FileUtils.mkdir_p ENV['HOME']
  # can maybe just use @mocked_home_directory

  # override GIT env vars
  author_name  = "SCM Puff"
  author_email = "scmpuff@test.local"
  set_env 'GIT_AUTHOR_NAME',     author_name
  set_env 'GIT_COMMITTER_NAME',  author_name
  set_env 'GIT_AUTHOR_EMAIL',    author_email
  set_env 'GIT_COMMITTER_EMAIL', author_email

end

After do
  # dont need to do these with set_env??
  #restore PATH
  #restore HOME
  #restore GIT env var
end

# since tmp/aruba is nested within the git repo of this program, we need to
# be somewhere else to test behavior of the binary when outside the git repo.
Before('@outside-repo') do
  @dirs = ["/tmp/aruba/scmpuff"]
end
