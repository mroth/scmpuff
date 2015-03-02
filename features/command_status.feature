Feature: status command

  Background:
    Given a mocked home directory

  @outside-repo
  Scenario: Appropriate error status when not in a git repo
    When I run `scmpuff status`
    Then the exit status should be 128
    #                              ^^^ same as `git status`
    And the output should contain:
      """
      Not a git repository (or any of the parent directories)
      """

  Scenario: Banner shows no changes when in an unchanged git repo
    Given I am in a git repository
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "No changes (working directory clean)"

  Scenario: Banner shows expansion reminder when in a changed git repo
    Given I am in a git repository
    And an empty file named "whatever"
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "|  [*] => $e*"

  Scenario: Banner shows current branch name
    Given I am in a git repository
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "On branch: master"

    When I switch to git branch "foobar"
    And  I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "On branch: foobar"

  Scenario: Banner shows position relative to remote status
    # Simulate a remote repository
    Given a git repository named "simulatedremote"
      And I cd to "simulatedremote"
      And a 4 byte file named "a.txt"
      And I successfully run the following commands:
        | git config receive.denyCurrentBranch ignore |
        | git add a.txt                               |
        | git commit -m "made a file"                 |
      And I cd to ".."
    Given I clone "simulatedremote" to "local"
      And I cd to "local"
      And a 4 byte file named "b.txt"
    # Check ahead of remote
    Given I successfully run the following commands:
      | git add b.txt                     |
      | git commit -m "made another file" |
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "|  +1  |"
    # Check behind from remote
    Given I successfully run the following commands:
      | git push         |
      | git reset HEAD~1 |
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "|  -1  |"

  Scenario: Status properly reports all file changes
    # See: http://git.io/IR8qcg for scm_breeze version of test
    Given I am in a complex working tree status matching scm_breeze tests
    When I run `scmpuff status`
    Then the exit status should be 0
      And the output should match / new file: *\[1\] *new_file/
      And the output should match /  deleted: *\[2\] *deleted_file/
      And the output should match / modified: *\[3\] *new_file/
      And the output should match /untracked: *\[4\] *untracked_file/


  @wip @focus
  Scenario: Status shows relative paths
    Some people might care about relative paths.
    TOOD: figure out how this should work w/r/t defaults and --relative.

    To do this, let's replicate the test_git_status_produces_relative_paths()
    test function from scm_breeze, located in status_shortcuts_test.sh:116,
    so that people who want to migrate over can have expected behavior.

    Given I am in a git repository
      And a directory named "dir1/sub1/subsub1"
      And a directory named "dir1/sub2"
      And a directory named "dir2"
      And an empty file named "dir1/sub1/subsub1/testfile"
      And an empty file named "dir1/sub2/testfile"
      And an empty file named "dir2/testfile"
      And I successfully run `git add .`

    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "dir1/sub1/subsub1/testfile"

    When I cd to "dir1"
    And I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain " sub1/subsub1/testfile"
     And the stdout from "scmpuff status" should contain " sub2/testfile"
     And the stdout from "scmpuff status" should contain "../dir2/testfile"

    When I cd to "sub1"
    And I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain " subsub1/testfile"
     And the stdout from "scmpuff status" should contain " ../sub2/testfile"
     And the stdout from "scmpuff status" should contain "../../dir2/testfile"

    When I cd to "subsub1"
    And I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain " testfile"
     And the stdout from "scmpuff status" should contain " ../../sub2/testfile"
     And the stdout from "scmpuff status" should contain "../../../dir2/testfile"


  Scenario: Status for a complex merge conflict
    Given PENDING: port from scm_breeze
  # TODO: port test_git_status_shortcuts_merge_conflicts()  from scm_breeze
