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
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type "git add 1"
      And I close the shell `<shell>`
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
      | fish  |


  Scenario Outline: Wrapped `git add` can handle files with spaces properly
    Given I am in a git repository
      And an empty file named "file with spaces.txt"
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type "git add 1"
      And I close the shell `<shell>`
    Then the exit status should be 0
    And the output should match /new file:\s+\[1\] file with spaces.txt/
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |


  Scenario Outline: Wrapped `git reset` can handle files with spaces properly
    This is different and more complex because `git status --porcelain` puts it
    inside quotes for the case where it is already added (but doesnt in the ??
    case surprisingly), and also it expands using --relative.

    Given I am in a git repository
      And an empty file named "file with spaces.txt"
    And I successfully run `git add "file with spaces.txt"`
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type "git reset 1"
      And I close the shell `<shell>`
    Then the exit status should be 0
    When I run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      untracked:  [1] file with spaces.txt
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |


  @recent-git-only
  Scenario Outline: Wrapped `git restore` works
    Given I am in a git repository
      And a 2 byte file named "foo.bar"
      And I successfully run `git add foo.bar`
      And I successfully run `git commit -m "initial commit"`
      And a 4 byte file named "foo.bar"
      And I successfully run `git add foo.bar`
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type "git restore --staged 1"
      And I close the shell `<shell>`
    Then the exit status should be 0
    When I run `scmpuff status`
    Then the stdout from "scmpuff status" should contain:
      """
      ➤ Changes not staged for commit
      #
      #       modified:  [1] foo.bar
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: Wrapped `git add` can handle shell expansions
    Given I am in a git repository
      And an empty file named "file with spaces.txt"
      And an empty file named "file2.txt"
      And an empty file named "untracked file.txt"
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type `<setfile>`
      And I type `git add "$FILE" 2`
      And I close the shell `<shell>`
    And the output should contain:
      """
      new file:  [1] file with spaces.txt
      """
    And the output should contain:
      """
      new file:  [2] file2
      """
    And the output should contain:
      """
      untracked:  [3] untracked file.txt
      """
    Then the exit status should be 0
    Examples:
      | shell | setfile                         |
      | bash  | FILE="file with spaces.txt"     |
      | zsh   | FILE="file with spaces.txt"     |
      | fish  | set FILE "file with spaces.txt" |
