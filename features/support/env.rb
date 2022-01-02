require 'aruba/cucumber'

Before do
  author_name  = "SCM Puff"
  author_email = "scmpuff@test.local"
  set_environment_variable 'GIT_AUTHOR_NAME',     author_name
  set_environment_variable 'GIT_COMMITTER_NAME',  author_name
  set_environment_variable 'GIT_AUTHOR_EMAIL',    author_email
  set_environment_variable 'GIT_COMMITTER_EMAIL', author_email
end

# since tmp/aruba is nested within the git repo of this program, we need to
# be somewhere else to test behavior of the binary when outside the git repo.
Before('@outside-repo') do
  tmpdir = Dir.mktmpdir("aruba")
  cd tmpdir
end

When(/I initialize scmpuff in `(.*)`/) do |shell|
  if shell == "fish"
    type %{scmpuff init -s --shell=fish | source}
  else
    type %{eval "$(scmpuff init -ws)"}
  end
end

When(/I close the shell `(.*)`/) do |shell|
  status_var = shell == "fish" ? "$status" : "$?"
  type "exit #{status_var}"
  # fish doesn't run the inputted commands until stdin is closed
  close_input
  step("I stop the command \"#{shell}\"")
end
