require 'aruba/cucumber'

Before do
  author_name  = "SCM Puff"
  author_email = "scmpuff@test.local"
  set_environment_variable 'GIT_AUTHOR_NAME',     author_name
  set_environment_variable 'GIT_COMMITTER_NAME',  author_name
  set_environment_variable 'GIT_AUTHOR_EMAIL',    author_email
  set_environment_variable 'GIT_COMMITTER_EMAIL', author_email
end
