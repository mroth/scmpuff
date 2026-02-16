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


  Scenario: Status properly reports all file changes w/ symlink
    Given I am in a complex working tree status matching scm_breeze tests
    When I cd to ".."
      And I run `sh -c 'ln -s mygitrepo symlink && cd symlink && scmpuff status'`
    Then the exit status should be 0
      And the output should match / new file: *\[1\] *new_file/
      And the output should match /  deleted: *\[2\] *deleted_file/
      And the output should match / modified: *\[3\] *new_file/
      And the output should match /untracked: *\[4\] *untracked_file/


  # CURRENT STATUS UNKNOWN! This behavior appears to have changed again in more
  # recent versions of git, as the test case does not trigger the expected
  # condition during `git status --short`. Disabling this test for now until it
  # can be reproduced reliably in modern git.

  # @recent-git-only
  # Scenario: Handle change detection properly
  #   Change detection is currently fairly rare in `git status`, mostly it only
  #   happens after in index via diff or show.  But it can occur, so make sure we
  #   support it when it happens, as it may be baked in better in the future.
  #   Change detection naturally involves two filepaths, like rename.

  #   In theory this is redundant with the "multiple filenames" scenario above,
  #   but since change detection seems somewhat in flux we want to test for it
  #   seperately in case its behavior changes in future versions of git.

  #   Thanks to @peff on git mailing list for conditions to reproduce this.

  #   Given I am in a git repository
  #   And a 1000 byte file named "file"
  #   And I successfully run the following commands:
  #     | git add file       |
  #     | git commit -m base |
  #     | mv file other      |
  #   Then I append to "file" with "foo"
  #   And I successfully run `git add .`
  #   # verify git behavior has not changed since this is hard to reproduce
  #   When I successfully run `git status --short`
  #   Then the stdout from "git status --short" should contain:
  #     """
  #     M  file
  #     C  file -> other
  #     """
  #   # actual behavior test
  #   When I successfully run `scmpuff status`
  #   Then the stdout from "scmpuff status" should contain:
  #     """
  #     modified:  [1] file
  #     """
  #   And the stdout from "scmpuff status" should contain:
  #     """
  #     copied:  [2] file -> other
  #     """

