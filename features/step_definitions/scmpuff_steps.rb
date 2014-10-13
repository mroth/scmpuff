include FileUtils

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

#
# Make table/list versions of common Aruba functions:
#

Given(/^I successfully run the following commands:$/) do |list|
  # list is a Cucumber::Ast::Table
  list.raw.each do |item|
    step "I successfully run `#{item.first}`"
  end
end
