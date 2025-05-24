Feature: branch command

  Background:
    Given a mocked home directory

  @outside-repo
  Scenario: Appropriate error status when not in a git repo
    When I run `scmpuff branch`
    Then the exit status should be 128
    And the output should contain:
      """
      Not a git repository (or any of the parent directories)
      """

  Scenario: Numbered output of local branches
    Given I am in a git repository
    And a file named "base" with:
      """
      foo
      """
    And I successfully run `git add base`
    And I successfully run `git commit -m "init"`
    And I successfully run the following commands:
      | git branch branch_a |
      | git branch branch_b |
    When I successfully run `scmpuff branch`
    Then the stdout from "scmpuff branch" should contain "* [1] master"
    And the stdout from "scmpuff branch" should contain "  [2] branch_a"
    And the stdout from "scmpuff branch" should contain "  [3] branch_b"

  Scenario: Detached HEAD output
    Given I am in a git repository
    And a file named "foo" with:
      """
      bar
      """
    And I successfully run `git add foo`
    And I successfully run `git commit -m "first"`
    And I successfully run `git branch feature`
    And I run `git checkout HEAD~0`
    When I successfully run `scmpuff branch`
    Then the stdout from "scmpuff branch" should match /\* \(HEAD detached/
    And the stdout from "scmpuff branch" should contain "[1] master"
    And the stdout from "scmpuff branch" should contain "[2] feature"

