Feature: scmpuff_status function
  The scmpuff_status shell function wraps the underlying
  `scmpuff status` command, passing along the `--filelist` option and then
  parsing the results to set environment variables in the current shell.

  Background:
    Given a mocked home directory

  @outside-repo
  Scenario Outline: Handle error conditions from wrapped binary command
    It is possible for the underlying `scmpuff status` command wrapped by our
    shell function to produce errors, for example, when not in a git repository.

    In keeping with the design theory of this program (handle as much as
    possible in the binary), we want to make sure those error messages are
    propogated to the user and not swallowed by the shell function, and that
    non-zero exit codes from the underlying process are preserved.

    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I close the shell `<shell>`
    Then the exit status should be 128
    Then the stderr should contain:
      """
      Not a git repository (or any of the parent directories)
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: Basic functionality works with shell wrapper.
    Given I am in a git repository
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I close the shell `<shell>`
    Then the exit status should be 0
    And the output should contain "No changes (working directory clean)"
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: Sets proper environment variables in shell
    Given I am in a complex working tree status matching scm_breeze tests
      And the scmpuff environment variables have been cleared
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type `echo -e "e1:$e1\ne2:$e2\ne3:$e3\ne4:$e4\ne5:$e5\n"`
      And I close the shell `<shell>`
    Then the output should match /^e1:.*new_file$/
      And the output should match /^e2:.*deleted_file$/
      And the output should match /^e3:.*new_file$/
      And the output should match /^e4:.*untracked_file$/
      And the output should match /^e5:$/
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: Sets proper environment variables in shell with weird filenames
    Given I am in a git repository
      And an empty file named "aa bb"
      And an empty file named "bb|cc"
      And an empty file named "cc*dd"
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type `echo -e "e1:$e1\ne2:$e2\ne3:$e3\ne4:$e4\n"`
      And I close the shell `<shell>`
    Then the output should match /^e1:.*aa bb$/
      And the output should match /^e2:.*bb\|cc$/
      And the output should match /^e3:.*cc\*dd$/
      And the output should match /^e4:$/
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: Clears extra environment variables from before
    Given I am in a complex working tree status matching scm_breeze tests
      And the scmpuff environment variables have been cleared
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "scmpuff_status"
      And I type "git add new_file"
      And I type "git commit -m 'so be it'"
      And I type "scmpuff_status"
      And I type `echo -e "e1:$e1\ne2:$e2\ne3:$e3\ne4:$e4\ne5:$e5\n"`
      And I close the shell `<shell>`
    Then the output should match /^e1:.*deleted_file$/
      And the output should match /^e2:.*untracked_file$/
      And the output should match /^e3:$/
      And the output should match /^e4:$/
      And the output should match /^e5:$/
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: default SCMPUFF_GIT_CMD is set to absolute path of a git command
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "echo $SCMPUFF_GIT_CMD"
      And I close the shell `<shell>`
    Then the output should match %r<^/.+/git$>
    # ^^ is absolute path to git: begins with a /, and ends with /git
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: SCMPUFF_GIT_CMD is set to absolute path of a git command, eliminating aliases
    When I run `<shell>` interactively
      And I type "alias git=/foo/bar"
      And I initialize scmpuff in `<shell>`
      And I type "echo $SCMPUFF_GIT_CMD"
    And I close the shell `<shell>`
    Then the output should match %r<^/.+/git$>
    # ^^ is absolute path to git: begins with a /, and ends with /git
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

  Scenario Outline: SCMPUFF_GIT_CMD respects existing environment variables
    When I run `<shell>` interactively
      And I type "export SCMPUFF_GIT_CMD=/foo/hub"
      And I initialize scmpuff in `<shell>`
      And I type "echo $SCMPUFF_GIT_CMD"
    And I close the shell `<shell>`
    Then the output should contain exactly "/foo/hub"
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |
