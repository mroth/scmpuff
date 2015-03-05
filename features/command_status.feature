Feature: status command

  Background:
    Given a mocked home directory

  @outside-repo
  Scenario: Appropriate error status when not in a git repo
    We can make this pretty, but we also want to be sure to use the same exit
    code as the normal 'git' command line client for consistency.

    When I run `scmpuff status`
    Then the exit status should be 128
    And the output should contain:
      """
      Not a git repository (or any of the parent directories)
      """

  Scenario: Banner shows no changes when in an unchanged git repo
    Given I am in a git repository
    When I successfully run `scmpuff status`
    Then the output should contain "No changes (working directory clean)"


  Scenario: Banner shows expansion reminder when in a changed git repo
    Given I am in a git repository
    And an empty file named "whatever"
    When I successfully run `scmpuff status`
    Then the output should contain "|  [*] => $e*"


  Scenario: Banner shows current branch name
    Given I am in a git repository
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "On branch: master"

    When I switch to git branch "foobar"
    And  I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain "On branch: foobar"


  Scenario: Banner shows position relative to remote status
    # Simulate a remote git repository situation
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


  @wip
  Scenario Outline: Handles file path magic properly for untracked files
    Given I am in a git repository
      And a directory named "foo"
      And an empty file named "foo/placeholder.txt"
      And I successfully run `git add .`
      And I successfully run `git commit -am.`
    Given an empty file named "<gitpath>"
    And I cd to "<cwd>"
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      untracked:  [1] <displaypath>
      """
    When I successfully run `scmpuff status -f --display=false`
    Then the stdout from "scmpuff status -f --display=false" should contain:
      """
      <abspath_end>
      """
    Examples:
      | cwd | gitpath    | abspath_end  | displaypath |
      | .   | a.txt      | /a.txt       | a.txt       |
      | .   | foo/b.txt  | /foo/b.txt   | foo/b.txt   |
      | foo | foo/b.txt  | /foo/b.txt   | b.txt       |
      | foo | a.txt      | /a.txt       | ../a.txt    |
      | .   | hi mom.txt | /hi mom.txt  | hi mom.txt  |
      | .   | (x).txt    | /(x).txt     | (x).txt     |

  @wip
  Scenario Outline: Handles file path magic properly for new files
    You would think this would be the same as the previous test, but in fact
    the way `git status --porcelain` outputs these is different, so we need
    to test this scenario seperately. (For example, prior to a4f2282 this would
    fail!)

    Given I am in a git repository
      And a directory named "foo"
      And an empty file named "foo/placeholder.txt"
      And I successfully run `git add .`
      And I successfully run `git commit -am.`
    Given an empty file named "<gitpath>"
    # below is the one difference!
    And I successfully run `git add .`
    And I cd to "<cwd>"
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      new file:  [1] <displaypath>
      """
    When I successfully run `scmpuff status -f --display=false`
    Then the stdout from "scmpuff status -f --display=false" should contain:
      """
      <abspath_end>
      """
    Examples:
      | cwd | gitpath    | abspath_end  | displaypath |
      | .   | a.txt      | /a.txt       | a.txt       |
      | .   | foo/b.txt  | /foo/b.txt   | foo/b.txt   |
      | foo | foo/b.txt  | /foo/b.txt   | b.txt       |
      | foo | a.txt      | /a.txt       | ../a.txt    |
      | .   | hi mom.txt | /hi mom.txt  | hi mom.txt  |
      | .   | (x).txt    | /(x).txt     | (x).txt     |

  @wip
  Scenario: Handle changes involving multiple filenames properly (UI)
    Certain operations (rename) can involve multiple filenames.

    The ideal scenario is that the destination filename gets set as the path for
    environment (so it can be references in git cmds), and the display shows a
    pretty arrowized status, e.g. foo -> bar, which should also be path aware.

    Given I am in a git repository
    And an empty file named "a.txt"
    And a directory named "foo"
    And I successfully run the following commands:
      | git add a.txt      |
      | git commit -am.    |
      | git mv a.txt b.txt |
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      renamed:  [1] a.txt -> b.txt
      """
    When I cd to "foo"
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      renamed:  [1] ../a.txt -> ../b.txt
      """

  @wip
  Scenario: Handle changes involving multiple filenames properly (vars)
    Given I am in a git repository
    And an empty file named "a.txt"
    And I successfully run the following commands:
      | git add a.txt      |
      | git commit -am.    |
      | git mv a.txt b.txt |
    When I successfully run `scmpuff status -f --display=false`
    Then the stdout from "scmpuff status -f --display=false" should contain:
      """
      /tmp/aruba/b.txt\n
      """
    When I cd to "foo"
    When I successfully run `scmpuff status`
    Then the stdout from "scmpuff status -f --display=false" should contain:
      """
      /tmp/aruba/b.txt\n
      """

  Scenario: Status shows relative paths (scm_breeze reference)
    SCM Breeze handles relative paths in a particular way.

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
    Test by duplicating exactly the test_git_status_shortcuts_merge_conflicts()
    tests from scm_breeze.

    Given I am in a git repository
    And an empty file named "both_modified"
    And an empty file named "both_deleted"
    And an empty file named "deleted_by_them"
    And an empty file named "deleted_by_us"
    And a file named "renamed_file" with:
      """
      renamed file needs some content
      """
    And I successfully run `git add both_modified both_deleted renamed_file deleted_by_them deleted_by_us`
    And I successfully run `git commit -m "First commit"`

    And I successfully run `git checkout -b conflict_branch`
    And a file named "both_added" with:
      """
      added by branch
      """
    And I append to "both_modified" with "branch line"
    And I append to "deleted_by_us" with "deleted by us"
    And I successfully run `git rm deleted_by_them both_deleted`
    And I successfully run `git mv renamed_file renamed_file_on_branch`
    And I successfully run `git add both_added both_modified deleted_by_us`
    And I successfully run `git commit -m "Branch commit"`

    And I successfully run `git checkout master`
    And I append to "both_added" with "added by master"
    And I append to "both_modified" with "master line"
    And I append to "deleted_by_them" with "deleted by them"
    And I successfully run `git rm deleted_by_us both_deleted`
    And I successfully run `git mv renamed_file renamed_file_on_master`
    And I successfully run `git add both_added both_modified deleted_by_them`
    And I successfully run `git commit -m "Master commit"`
    And I run `git merge conflict_branch`

    When I successfully run `scmpuff status`
    Then the output should match /both added: *\[[0-9]*\] *both_added/
    Then the output should match /both modified: *\[[0-9]*\] *both_modified/
    Then the output should match /deleted by them: *\[[0-9]*\] *deleted_by_them/
    Then the output should match /deleted by us: *\[[0-9]*\] *deleted_by_us/
    Then the output should match /both deleted: *\[[0-9]*\] *renamed_file/
    Then the output should match /added by them: *\[[0-9]*\] *renamed_file_on_branch/
    Then the output should match /added by us: *\[[0-9]*\] *renamed_file_on_master/
