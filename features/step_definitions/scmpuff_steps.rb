require 'fileutils'
require 'aruba/api'

Given(/^a git repository named "([^"]*)"$/) do |repo_name|
  steps %(
    Given I successfully run `git init --quiet #{repo_name}`
  )
end

Given(/^I am in a git repository named "([^"]*)"$/) do |repo_name|
  steps %(
    Given a git repository named "#{repo_name}"
    And I cd to "#{repo_name}"
  )
end

Given(/^I am in a git repository$/) do
  repo_name = 'mygitrepo'
  steps %(
    Given a git repository named "#{repo_name}"
    And I cd to "#{repo_name}"
  )
end

Given(/^I switch to git branch "([^"]*)"$/) do |branch_name|
  steps %(
    Given I successfully run `git checkout -b #{branch_name}`
  )
end

Given(/^I switch to existing git branch "([^"]*)"$/) do |branch_name|
  steps %(
    Given I successfully run `git checkout #{branch_name}`
  )
end

Given(/^I clone "(.*?)" to "(.*?)"$/) do |r1, r2|
  steps %(
    Given I successfully run `git clone #{r1}/.git #{r2}`
  )
end

Given(/^I am in a complex working tree status matching scm_breeze tests$/) do
  steps %(
    Given I am in a git repository
    And an empty file named "deleted_file"
    And I successfully run `git add deleted_file`
    And I successfully run `git commit -m "Test commit"`
    And an empty file named "new_file"
    And an empty file named "untracked_file"
    And I successfully run `git add new_file`
    And I overwrite "new_file" with:
      """
      changed contents lolol
      """
    And I remove the file "deleted_file"
  )
end

# Create a filesystem mock of the git repo, and copy it in.
# This speeds up tests because shell commands are the slowest thing, so we dont
# want to rerun the same git init every iteration, rather we just copy a fresh
# copy of the mocked directory!
Given(/^I am in the mocked git repository with commited subdirectory and file$/) do
  MOCK ||= File.join(expand_path(".."), "mock", "gitsubdir") #needs to be outside of aruba clobber dir
  unless File.directory? MOCK
    FileUtils.mkdir_p MOCK
    Dir.chdir MOCK do
      FileUtils.mkdir "foo"
      FileUtils.touch "foo/placeholder.txt"
      system "git init --quiet"
      system "git config --local user.name 'scmpuff mocker'"
      system "git config --local user.email mocked@scmpuff.github.io"
      system "git add ."
      system "git commit -m."
    end
  end
  FileUtils.cp_r MOCK, expand_path(".")
  cd "gitsubdir"
end

Given(/^the scmpuff environment variables have been cleared$/) do
  (1..50).each do |n|
    delete_environment_variable "e#{n}"
  end
end

# identical to what is in the built-in aruba steps, but without stupid
# forced variable upcasing, since I use lowercase in the app!
Given(/^I override the environment variables to:/) do |table|
  table.hashes.each do |row|
    variable = row['variable'].to_s
    value = row['value'].to_s
    set_environment_variable(variable, value)
  end
end

Given(/^I override environment variable "(.*?)" to the absolute path of "(.*?)"$/) do |e, f|
  filepath = expand_path(f)
  set_environment_variable(e, filepath)
end

#
# Handle unknown absolute paths in output
#
Then(/^the stdout from "([^"]*)" should contain the absolute path of "([^"]*)"$/) do |cmd, f|
  filepath = expand_path(f)
  step %Q(the stdout from "#{cmd}" should contain "#{filepath}")
end


#
# Backtick version for "when I type" step to enable passing quotation marks
#
When(/^I type `(.*?)`$/) do |cmd|
  type(cmd)
end

#
# Make table/list versions of common Aruba functions:
#
Given(/^I successfully run the following commands:$/) do |list|
  # list is a Cucumber::Ast::Table
  list.raw.each do |item|
    step "I successfully run `#{item.first}`"
  end
end

#
# Helpful to define actual pending tests
#
Given(/^PENDING/) do
  pending
end
