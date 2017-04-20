Feature: optional wrapping of normal git cmds in the shell
  In order to verify the shell git wrappers work correctly
  I want to make sure they intercept and wrap naked git commands properly

  Background:
    Given a mocked home directory

  Scenario Outline: Wrapped `git add` adds by number and echos status after
    Given I am in a git repository
      And a 4 byte file named "foo.bar"
      And a 4 byte file named "bar.foo"
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -ws)"`
      And I type "scmpuff_status"
      And I type "git add 1"
      And I type "exit"
    Then the output should contain:
      """
      # On branch: master  |  [*] => $e*
      #
      ➤ Changes to be committed
      #
      #       new file:  [1] bar.foo
      #
      ➤ Untracked files
      #
      #      untracked:  [2] foo.bar
      #
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |


  Scenario Outline: Wrapped `git add` can handle files with spaces properly
    Given I am in a git repository
      And an empty file named "file with spaces.txt"
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -ws)"`
      And I type "scmpuff_status"
      And I type "git add 1"
      And I type "exit"
    Then the exit status should be 0
    And the output should match /new file:\s+\[1\] file with spaces.txt/
    Examples:
      | shell |
      | bash  |
      | zsh   |


  Scenario Outline: Wrapped `git reset` can handle files with spaces properly
    This is different and more complex because `git status --porcelain` puts it
    inside quotes for the case where it is already added (but doesnt in the ??
    case surprisingly), and also it expands using --relative.

    Given I am in a git repository
    Given an empty file named "file with spaces.txt"
    ###AAA
    Given I successfully run `git add "file with spaces.txt"`

    ###BBB
    Given I run `<shell>` interactively
      And I type `eval "$(scmpuff init -ws)"`
      And I type "scmpuff_status"
      And I type "git reset 1"
      And I type "echo 'DEBUG: PHASE BBB'"
      And I type "exit"
    Then the exit status should be 0
    ### ^^^ this is checking exit status of AAA, not BBB!

    # Then the output should contain "BBB"
    ### ^^^ this command however, would force aruba to wait until BBB completes (super hacky option #1)

    # Then I stop the command started last
    ### ^^^ (hacky option #2) also seems to ensure BBB completes, but throws a deprecation error

    ###CCC
    ### currently, this phase is failing, because it is taking place *BEFORE* BBB completes
    When I run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      untracked:  [1] file with spaces.txt
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |
