require 'aruba/api'

Given /^a git repository named "([^"]*)"$/ do |repo_name|
  # system `git init --quiet #{repo_name}`
  steps %Q{
    Given I successfully run `git init --quiet #{repo_name}`
  }
end

Given /^I am in a git repository named "([^"]*)"$/ do |repo_name|
  steps %Q{
    Given a git repository named "#{repo_name}"
    And I cd to "#{repo_name}"
  }
end

Given /^I am in a git repository$/ do
  repo_name = 'mygitrepo'
  steps %Q{
    Given a git repository named "#{repo_name}"
    And I cd to "#{repo_name}"
  }
end

Given /^I switch to git branch "([^"]*)"$/ do |branch_name|
  steps %Q{
    Given I successfully run `git checkout -b #{branch_name}`
  }
end

Given(/^I clone "(.*?)" to "(.*?)"$/) do |r1, r2|
  steps %Q{
    Given I successfully run `git clone #{r1}/.git #{r2}`
  }
end

Given /^I am in a complex working tree status matching scm_breeze tests$/ do
  steps %Q{
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
  }
end

Given(/^the scmpuff environment variables have been cleared$/) do
  (1..50).each do |n|
    set_env("e#{n}", nil)
  end
end


# TODO: no longer needed?
#Then(/^the environment variable "(.*?)" should equal the absolute path for "(.*?)"$/) do |var, filename|
#  expect(ENV[var]).to eq(File.expand_path("~/mygitrepo/#{filename}"))
#end


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
